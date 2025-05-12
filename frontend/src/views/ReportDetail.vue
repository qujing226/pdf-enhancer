<template>
  <div class="report-detail">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h2>{{ report.title }}</h2>
          <div class="header-actions">
            <el-button type="primary" @click="downloadPDF">下载PDF</el-button>
            <el-button type="success" @click="generateSummary" :loading="summaryLoading" :disabled="summaryLoading">
              {{ report.summary ? '重新生成摘要' : '生成摘要' }}
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="report-info">
        <p><strong>创建时间：</strong>{{ formatDateTime(report.created_at) }}</p>
        <p><strong>报告ID：</strong>{{ report.report_id }}</p>
      </div>

      <el-divider>摘要</el-divider>
      <div class="report-summary">
        <p v-if="report.summary">{{ report.summary }}</p>
        <el-empty v-else description="暂无摘要，点击'生成摘要'按钮生成" />
      </div>

      <el-divider>报告内容</el-divider>
      <div class="report-content">
        <div v-if="pdfPages.length > 0" class="pdf-viewer">
          <div class="pdf-pagination">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="1"
              layout="prev, pager, next"
              :total="pdfPages.length"
              @current-change="handlePageChange"
            />
            <span class="page-info">{{ currentPage }} / {{ pdfPages.length }}</span>
          </div>
          
          <div class="pdf-page">
            <div v-html="pdfPages[currentPage - 1]" class="pdf-content"></div>
          </div>
        </div>
        <el-empty v-else description="加载PDF内容中..." />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { reportApi } from '../api'
import dayjs from 'dayjs'

const route = useRoute()
const report = ref({})
const pdfPages = ref([])
const currentPage = ref(1)
const summaryLoading = ref(false)

// 获取报告详情
const fetchReportDetail = async () => {
  try {
    const reportId = route.params.id
    const response = await reportApi.getReportById(reportId)
    report.value = response.data.data.report
    
    // 获取PDF内容
    fetchPdfContent(reportId)
  } catch (error) {
    ElMessage.error('获取报告详情失败')
    console.error('获取报告详情失败:', error)
  }
}

// 获取PDF内容
const fetchPdfContent = async (reportId) => {
  try {
    // 这里应该调用后端API获取PDF内容
    // 由于后端API可能尚未实现，这里仍使用模拟数据
    simulatePdfPages()
  } catch (error) {
    ElMessage.error('获取PDF内容失败')
    console.error('获取PDF内容失败:', error)
  }
}

// 模拟PDF内容分页（实际项目中应从后端获取PDF内容）
const simulatePdfPages = () => {
  // 这里模拟PDF内容，实际项目中应从后端获取
  const dummyContent = report.value.content || '这是报告内容示例。'
  
  // 创建5页模拟内容
  pdfPages.value = Array(5).fill().map((_, index) => {
    return `<div class="pdf-page-content">
      <h3>第 ${index + 1} 页</h3>
      <p>${dummyContent}</p>
      <p>这是PDF内容的第 ${index + 1} 页，实际项目中应从后端获取真实PDF内容。</p>
    </div>`
  })
}

// 处理页面切换
const handlePageChange = (page) => {
  currentPage.value = page
}

// 生成摘要
const generateSummary = async () => {
  try {
    summaryLoading.value = true
    const reportId = route.params.id
    const response = await reportApi.generateSummary(reportId)
    report.value.summary = response.data.data.summary
    report.value.has_summary = true
    ElMessage.success('摘要生成成功')
  } catch (error) {
    ElMessage.error('摘要生成失败')
    console.error('摘要生成失败:', error)
  } finally {
    summaryLoading.value = false
  }
}

// 下载PDF
const downloadPDF = () => {
  const reportId = route.params.id
  const token = localStorage.getItem('token')
  const pdfUrl = reportApi.getReportPdfUrl(reportId)
  // 创建一个带有认证信息的链接
  const downloadLink = document.createElement('a')
  downloadLink.href = `${pdfUrl}?token=${token}`
  downloadLink.target = '_blank'
  document.body.appendChild(downloadLink)
  downloadLink.click()
  document.body.removeChild(downloadLink)
}

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return ''
  return dayjs(dateTimeStr).format('YYYY-MM-DD HH:mm')
}

// 页面加载时获取报告详情
onMounted(() => {
  fetchReportDetail()
})
</script>

<style scoped>
.report-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.box-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 10px;
}

.card-header h2 {
  margin: 0;
  font-size: 1.5rem;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.report-info {
  margin-bottom: 20px;
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.report-info p {
  margin: 0;
}

.report-summary {
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
  line-height: 1.6;
}

.pdf-viewer {
  margin-top: 20px;
}

.pdf-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-info {
  margin-left: 15px;
  color: #606266;
  font-weight: bold;
}

.pdf-page {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 20px;
  min-height: 500px;
  background-color: #fff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow-x: auto;
}

.pdf-content {
  line-height: 1.6;
  font-size: 16px;
}

.pdf-page-content h3 {
  color: #409EFF;
  margin-bottom: 15px;
  text-align: center;
}

.pdf-page-content p {
  margin-bottom: 15px;
  text-align: justify;
}

@media screen and (max-width: 768px) {
  .report-detail {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    margin-top: 10px;
    width: 100%;
    justify-content: space-between;
  }
  
  .pdf-page {
    padding: 10px;
    min-height: 400px;
  }
}
</style>