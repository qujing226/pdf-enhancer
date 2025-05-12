import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// 如果本地存储中有token，则设置默认请求头
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

// 报告相关API
export const reportApi = {
  // 获取报告列表
  getReports() {
    return axios.get('/api/v1/reports')
  },
  
  // 获取报告详情
  getReportById(reportId) {
    return axios.get(`/api/v1/report/${reportId}`)
  },
  
  // 生成报告摘要
  generateSummary(reportId) {
    return axios.post(`/api/v1/report/${reportId}/summary`)
  },
  
  // 获取报告PDF的URL
  getReportPdfUrl(reportId) {
    return `${axios.defaults.baseURL}/api/v1/report/${reportId}/pdf`
  },
  
  // 上传报告
  uploadReport(formData) {
    return axios.post('/api/v1/reports/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

// 用户相关API
export const userApi = {
  // 用户登录
  login(email, password) {
    return axios.post('/api/v1/login', { email, password })
  },
  
  // 用户注册
  register(name, email, password) {
    return axios.post('/api/v1/register', { name, email, password })
  }
}