package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ledongthuc/pdf"
)

// ParsePDFText 从 io.Reader 读取 PDF 内容并提取文本
func ParsePDFText(reader io.Reader) (string, error) {
	// 读取全部内容到字节切片
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("读取PDF内容失败: %v", err)
	}

	// 将字节切片转换为 io.ReadSeeker（PDF解析需要随机访问）
	readSeeker := bytes.NewReader(content)

	// 创建PDF解析器
	pdfReader, err := pdf.NewReader(readSeeker, int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("解析PDF失败: %v", err)
	}

	// 提取所有页面文本
	var fullText string
	totalPages := pdfReader.NumPage()
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("提取第%d页文本失败: %v", pageNum, err)
		}
		fullText += text + "\n"
	}

	return fullText, nil
}
