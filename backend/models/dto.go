package models

import (
	"time"
)

// ReportDTO 报告数据传输对象，用于API交互
type ReportDTO struct {
	ID        string    `json:"report_id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` // API层的内容字段
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PDFPath   string    `json:"pdf_path"`
}

// ToReport 将DTO转换为业务模型
func (dto *ReportDTO) ToReport() *Report {
	return &Report{
		ID:        dto.ID,
		UserID:    dto.UserID,
		Title:     dto.Title,
		Content:   dto.Content,
		Summary:   dto.Summary,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		PDFPath:   dto.PDFPath,
	}
}

// ToReportDTO 将业务模型转换为DTO
func ToReportDTO(report *Report) *ReportDTO {
	return &ReportDTO{
		ID:        report.ID,
		UserID:    report.UserID,
		Title:     report.Title,
		Content:   report.Content,
		Summary:   report.Summary,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
		PDFPath:   report.PDFPath,
	}
}
