package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePDFText_ValidPDFFile(t *testing.T) {
	// 获取当前文件所在目录（utils），然后定位到 backend/download.pdf
	projectRoot := filepath.Join("..") // 从 utils 回到 backend 目录
	pdfPath := filepath.Join(projectRoot, "download.pdf")

	file, err := os.Open(pdfPath)
	if err != nil {
		t.Fatalf("打开PDF文件失败: %v", err)
	}
	defer file.Close()

	// 调用被测函数
	text, err := ParsePDFText(file)
	if err != nil {
		t.Errorf("解析PDF时出错: %v", err)
	}

	// 检查是否成功提取内容
	assert.NotEmpty(t, text, "期望解析出文本内容，但结果为空")
	
	t.Logf("提取的文本内容: %s", text)
}
