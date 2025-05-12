<template>
  <div class="report-list">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h2>报告列表</h2>
          <el-button type="primary" @click="showUploadDialog">上传新报告</el-button>
        </div>
      </template>
      
      <el-table :data="reports" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="报告标题" width="300" />
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="has_summary" label="摘要状态" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.has_summary ? 'success' : 'info'">
              {{ scope.row.has_summary ? '已生成' : '未生成' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="摘要" min-width="200">
          <template #default="scope">
            <div v-if="scope.row.has_summary && scope.row.summary" class="summary-preview">
              {{ scope.row.summary.length > 100 ? scope.row.summary.substring(0, 100) + '...' : scope.row.summary }}
            </div>
            <span v-else class="no-summary">暂无摘要</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewReport(scope.row.report_id)">
              查看
            </el-button>
            <el-button type="success" size="small" @click="downloadPDF(scope.row.report_id)">
              下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 上传报告对话框 -->
      <el-dialog v-model="uploadDialogVisible" title="上传新报告" width="500px">
        <el-form :model="uploadForm" label-width="80px" ref="uploadFormRef">
          <el-form-item label="报告标题" prop="title" :rules="[{ required: true, message: '请输入报告标题', trigger: 'blur' }]">
            <el-input v-model="uploadForm.title" placeholder="请输入报告标题"></el-input>
          </el-form-item>
          <el-form-item label="PDF文件" prop="file" :rules="[{ required: true, message: '请选择PDF文件', trigger: 'change' }]">
            <el-upload
              class="upload-demo"
              action="#"
              :auto-upload="false"
              :on-change="handleFileChange"
              :limit="1"
              accept=".pdf"
            >
              <template #trigger>
                <el-button type="primary">选择文件</el-button>
              </template>
              <template #tip>
                <div class="el-upload__tip">只能上传PDF文件</div>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="uploadDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="uploadReport" :loading="uploading">上传</el-button>
          </span>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { reportApi } from '../api'
import dayjs from 'dayjs'

const router = useRouter()
const reports = ref([])
const loading = ref(false)
const uploadDialogVisible = ref(false)
const uploading = ref(false)
const uploadFormRef = ref(null)
const uploadForm = ref({
  title: '',
  file: null
})

// 获取报告列表
const fetchReports = async () => {
  loading.value = true
  try {
    const response = await reportApi.getReports()
    reports.value = response.data.data
  } catch (error) {
    ElMessage.error('获取报告列表失败')
    console.error('获取报告列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 查看报告详情
const viewReport = (reportId) => {
  router.push(`/report/${reportId}`)
}

// 下载PDF
const downloadPDF = (reportId) => {
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

// 显示上传对话框
const showUploadDialog = () => {
  uploadDialogVisible.value = true
  uploadForm.value = {
    title: '',
    file: null
  }
}

// 处理文件选择
const handleFileChange = (file) => {
  uploadForm.value.file = file.raw
}

// 上传报告
const uploadReport = async () => {
  if (!uploadFormRef.value) return
  
  await uploadFormRef.value.validate(async (valid) => {
    if (valid && uploadForm.value.file) {
      uploading.value = true
      try {
        const formData = new FormData()
        formData.append('title', uploadForm.value.title)
        formData.append('file', uploadForm.value.file)
        
        await reportApi.uploadReport(formData)
        ElMessage.success('报告上传成功')
        uploadDialogVisible.value = false
        fetchReports() // 刷新报告列表
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '上传报告失败')
        console.error('上传报告失败:', error)
      } finally {
        uploading.value = false
      }
    }
  })
}

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return ''
  return dayjs(dateTimeStr).format('YYYY-MM-DD HH:mm')
}

// 页面加载时获取报告列表
onMounted(() => {
  fetchReports()
})
</script>

<style scoped>
.report-list {
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

.summary-preview {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.no-summary {
  color: #909399;
  font-style: italic;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media screen and (max-width: 768px) {
  .report-list {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .el-table {
    width: 100%;
    overflow-x: auto;
  }
}
</style>