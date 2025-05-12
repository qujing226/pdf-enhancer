package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/services"
	"github.com/qujing226/pdf-enhancer/backend/utils"
)

// ReportHandler 处理报告相关的请求
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler 创建新的报告处理器
func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// UploadReport 处理报告上传请求
func (h *ReportHandler) UploadReport(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.NewAPIResponse(http.StatusUnauthorized, "未授权的访问", nil))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewAPIResponse(http.StatusBadRequest, "未找到上传的文件", err.Error()))
		return
	}
	defer file.Close()

	// 获取文件标题
	title := header.Filename

	// 创建报告
	report, err := h.reportService.CreateReportFromUpload(c, userID, title, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "上传报告失败", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.NewAPIResponse(http.StatusCreated, "上传成功", report))
}

// GetReports 获取用户的报告列表
func (h *ReportHandler) GetReports(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.NewAPIResponse(http.StatusUnauthorized, "未授权的访问", nil))
		return
	}

	// 获取报告列表
	reports, err := h.reportService.GetReportsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "获取报告列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewAPIResponse(http.StatusOK, "获取成功", reports))
}

// GetReport 获取单个报告详情
func (h *ReportHandler) GetReport(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.NewAPIResponse(http.StatusUnauthorized, "未授权的访问", nil))
		return
	}

	// 获取报告ID
	reportID := c.Param("report_id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, models.NewAPIResponse(http.StatusBadRequest, "无效的报告ID", nil))
		return
	}

	// 获取报告详情
	report, err := h.reportService.GetReportByID(reportID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "获取报告详情失败", err.Error()))
		return
	}

	pdfURL := "http://localhost:8080/api/v1/report/" + reportID + ".pdf"

	c.JSON(http.StatusOK, models.NewAPIResponse(http.StatusOK, "获取成功", models.ReportDetailResponse{
		Report: *report,
		PDFURL: pdfURL,
	}))
}

// GetReportPDF 获取报告的PDF文件
func (h *ReportHandler) GetReportPDF(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.NewAPIResponse(http.StatusUnauthorized, "未授权的访问", nil))
		return
	}

	// 获取报告ID
	reportID := c.Param("report_id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, models.NewAPIResponse(http.StatusBadRequest, "无效的报告ID", nil))
		return
	}

	// 获取PDF文件流
	pdfStream, objInfo, err := h.reportService.GetReportPDF(reportID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "获取PDF文件失败", err.Error()))
		return
	}
	defer pdfStream.Close()

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objInfo.Key))
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// 发送文件
	c.DataFromReader(http.StatusOK, objInfo.Size, "application/pdf", pdfStream, nil)
}

// GenerateSummary 生成报告摘要
func (h *ReportHandler) GenerateSummary(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.NewAPIResponse(http.StatusUnauthorized, "未授权的访问", nil))
		return
	}

	// 获取报告ID
	reportID := c.Param("report_id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, models.NewAPIResponse(http.StatusBadRequest, "无效的报告ID", nil))
		return
	}

	// 获取报告详情
	report, err := h.reportService.GetReportByID(reportID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "获取报告详情失败", err.Error()))
		return
	}

	// 生成摘要
	summary, err := h.reportService.GenerateSummary(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewAPIResponse(http.StatusInternalServerError, "生成摘要失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewAPIResponse(http.StatusOK, "生成成功", gin.H{"summary": summary}))
}
