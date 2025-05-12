package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/repository/dao_models"
)

// IReportRepository 报告仓储接口
type IReportRepository interface {
	Create(ctx context.Context, report *models.Report) error
	GetByID(ctx context.Context, reportID string, userID string) (*models.Report, error)
	GetByUserID(ctx context.Context, userID string) ([]models.ReportListItem, error)
	UpdateSummary(ctx context.Context, reportID string, summary string) error
}

// ReportRepository 报告仓储实现
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository 创建报告仓储实例
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Create 创建新报告
func (r *ReportRepository) Create(ctx context.Context, report *models.Report) error {
	// 转换为DAO模型
	reportDAO := &dao_models.ReportDAO{
		ID:        report.ID,
		UserID:    report.UserID,
		Title:     report.Title,
		Content:   report.Content, // 内容将存储在longtext字段中
		Summary:   report.Summary,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
		PDFPath:   report.PDFPath,
	}

	query := `INSERT INTO reports (id, user_id, title, content, summary, created_at, updated_at, pdf_path)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		reportDAO.ID, reportDAO.UserID, reportDAO.Title, reportDAO.Content,
		reportDAO.Summary, reportDAO.CreatedAt, reportDAO.UpdatedAt, reportDAO.PDFPath)
	if err != nil {
		return fmt.Errorf("保存报告到数据库失败: %w", err)
	}
	return nil
}

// GetByID 根据报告ID和用户ID获取报告
func (r *ReportRepository) GetByID(ctx context.Context, reportID string, userID string) (*models.Report, error) {
	query := `SELECT id, user_id, title, content, summary, created_at, updated_at, pdf_path 
	          FROM reports WHERE id = ? AND user_id = ?`

	// 使用DAO模型接收数据库数据
	reportDAO := &dao_models.ReportDAO{}
	err := r.db.QueryRowContext(ctx, query, reportID, userID).Scan(
		&reportDAO.ID, &reportDAO.UserID, &reportDAO.Title, &reportDAO.Content,
		&reportDAO.Summary, &reportDAO.CreatedAt, &reportDAO.UpdatedAt, &reportDAO.PDFPath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			return nil, fmt.Errorf("报告不存在或无权访问")
		}
		return nil, fmt.Errorf("查询报告详情失败: %w", err)
	}

	// 转换为业务模型
	report := &models.Report{
		ID:        reportDAO.ID,
		UserID:    reportDAO.UserID,
		Title:     reportDAO.Title,
		Content:   reportDAO.Content,
		Summary:   reportDAO.Summary,
		CreatedAt: reportDAO.CreatedAt,
		UpdatedAt: reportDAO.UpdatedAt,
		PDFPath:   reportDAO.PDFPath,
	}

	return report, nil
}

// GetByUserID 获取用户的所有报告
func (r *ReportRepository) GetByUserID(ctx context.Context, userID string) ([]models.ReportListItem, error) {
	query := `SELECT id, title, created_at, summary FROM reports WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询报告列表失败: %w", err)
	}
	defer rows.Close()

	var reports []models.ReportListItem
	for rows.Next() {
		var reportItem models.ReportListItem
		var summary sql.NullString
		if err := rows.Scan(&reportItem.ReportID, &reportItem.Title, &reportItem.CreatedAt, &summary); err != nil {
			return nil, fmt.Errorf("扫描报告数据失败: %w", err)
		}
		reportItem.HasSummary = summary.Valid && summary.String != ""
		reports = append(reports, reportItem)
	}
	return reports, nil
}

// UpdateSummary 更新报告摘要
func (r *ReportRepository) UpdateSummary(ctx context.Context, reportID string, summary string) error {
	query := `UPDATE reports SET summary = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, summary, time.Now(), reportID)
	if err != nil {
		return fmt.Errorf("更新报告摘要到数据库失败: %w", err)
	}
	return nil
}
