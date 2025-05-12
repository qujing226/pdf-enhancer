package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/repository"
	"github.com/qujing226/pdf-enhancer/backend/utils"
)

// ReportService 报告服务
type ReportService struct {
	reportRepo     repository.IReportRepository
	MinioClient    *minio.Client
	DeepSeekClient *DeepSeekClient
	BucketName     string
}

// NewReportService 创建新的报告服务
func NewReportService(reportRepo repository.IReportRepository, minioClient *minio.Client, deepseekClient *DeepSeekClient, bucketName string) *ReportService {
	return &ReportService{
		reportRepo:     reportRepo,
		MinioClient:    minioClient,
		DeepSeekClient: deepseekClient,
		BucketName:     bucketName,
	}
}

// CreateReportFromUpload 从上传的文件创建报告
func (s *ReportService) CreateReportFromUpload(ctx context.Context, userID string, title string, fileHeader *multipart.FileHeader) (*models.Report, error) {
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer file.Close()

	// 提取PDF文本内容
	// 由于 ParsePDFText 需要 io.Reader，而 file 已经是 io.Reader，可以直接传递
	// 但为了确保文件可以被多次读取（例如，一次用于解析，一次用于上传），我们先读取到内存
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 重置文件读取器以便 ParsePDFText 可以从头读取
	fileReaderForParse := bytes.NewReader(body)
	textContent, err := utils.ParsePDFText(fileReaderForParse)
	if err != nil {
		log.Printf("解析PDF文本内容失败: %v", err) // 记录错误，但可能仍希望保存文件
		// 根据需求决定是否在此处返回错误，或者允许没有文本内容的报告
	}

	// 生成报告ID和PDF文件名
	reportID := uuid.New().String()
	pdfObjectName := fmt.Sprintf("%s.pdf", reportID)

	// 上传PDF到MinIO
	// 重置文件读取器以便上传
	fileReaderForUpload := bytes.NewReader(body)
	_, err = s.MinioClient.PutObject(ctx, s.BucketName, pdfObjectName, fileReaderForUpload, int64(len(body)), minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"), // "application/pdf"
	})
	if err != nil {
		return nil, fmt.Errorf("上传PDF到MinIO失败: %w", err)
	}

	// 创建报告模型
	report := &models.Report{
		ID:        reportID,
		UserID:    userID,
		Title:     title,
		Content:   textContent, // 保存提取的文本内容
		Summary:   "",          // 初始摘要为空
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		PDFPath:   filepath.Join(s.BucketName, pdfObjectName), // MinIO中的对象路径
	}

	// 将报告保存到数据库
	err = s.reportRepo.Create(ctx, report)
	if err != nil {
		// 如果数据库插入失败，可能需要考虑删除已上传到MinIO的文件以保持一致性
		// s.MinioClient.RemoveObject(ctx, s.BucketName, pdfObjectName, minio.RemoveObjectOptions{})
		return nil, err
	}

	return report, nil
}

// GetReportsByUserID 根据用户ID获取报告列表
func (s *ReportService) GetReportsByUserID(userID string) ([]models.ReportListItem, error) {
	return s.reportRepo.GetByUserID(context.Background(), userID)
}

// GetReportByID 根据报告ID和用户ID获取报告详情
func (s *ReportService) GetReportByID(reportID string, userID string) (*models.Report, error) {
	return s.reportRepo.GetByID(context.Background(), reportID, userID)
}

// GeneratePDFURL 生成PDF的预签名URL
func (s *ReportService) GeneratePDFURL(reportID string) (string, error) {
	// 假设PDFPath存储的是 bucketName/objectName
	report, err := s.GetReportByID(reportID, "") // 注意：这里可能需要正确的userID，或者修改GetReportByID的逻辑
	if err != nil {
		return "", fmt.Errorf("获取报告信息失败: %w", err)
	}

	reqParams := make(url.Values)
	reqParams["response-content-disposition"] = []string{fmt.Sprintf("attachment; filename=\"%s.pdf\"", report.Title)}

	presignedURL, err := s.MinioClient.PresignedGetObject(context.Background(), s.BucketName, filepath.Base(report.PDFPath), time.Second*24*60*60, reqParams)
	if err != nil {
		return "", fmt.Errorf("生成MinIO预签名URL失败: %w", err)
	}
	return presignedURL.String(), nil
}

// GetReportPDF 获取报告的PDF文件流和信息
func (s *ReportService) GetReportPDF(reportID, userID string) (io.ReadCloser, minio.ObjectInfo, error) {
	report, err := s.GetReportByID(reportID, userID) // 注意：这里可能需要正确的userID
	if err != nil {
		return nil, minio.ObjectInfo{}, fmt.Errorf("获取报告信息失败: %w", err)
	}

	object, err := s.MinioClient.GetObject(context.Background(), s.BucketName, filepath.Base(report.PDFPath), minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, fmt.Errorf("从MinIO获取对象失败: %w", err)
	}

	objInfo, err := object.Stat()
	if err != nil {
		object.Close() // Ensure the object is closed on error
		return nil, minio.ObjectInfo{}, fmt.Errorf("获取MinIO对象信息失败: %w", err)
	}

	return object, objInfo, nil
}

// GenerateSummary 生成报告摘要
func (s *ReportService) GenerateSummary(report *models.Report) (string, error) {
	if report.Content == "" {
		return "", fmt.Errorf("报告内容为空，无法生成摘要")
	}

	summary, err := s.DeepSeekClient.GenerateSummary(context.Background(), fmt.Sprintf("Title: %s\nContent:%s", report.Title, report.Content))
	if err != nil {
		return "", fmt.Errorf("调用DeepSeek API生成摘要失败: %w", err)
	}

	// 更新数据库中的摘要
	err = s.reportRepo.UpdateSummary(context.Background(), report.ID, summary)
	if err != nil {
		return "", err
	}

	return summary, nil
}

// Helper function to read file into a byte slice and return a new reader
// This is useful if you need to read the file multiple times (e.g., for parsing and uploading)
func getReusableReader(file multipart.File) (io.Reader, []byte, error) {
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewReader(body), body, nil
}

// bytes.NewReader is needed for ParsePDFText and PutObject if we read the file into memory first
// Ensure to import "bytes"
