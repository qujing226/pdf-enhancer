<script setup>
import { onMounted, computed } from 'vue'
import { useUserStore } from './stores/user'
import { useRoute } from 'vue-router'
import axios from 'axios'
import MainLayout from './components/MainLayout.vue'

const userStore = useUserStore()
const route = useRoute()

// 计算当前路由是否需要主布局
const needsLayout = computed(() => {
  return route.meta.requiresAuth === true
})

// 设置请求拦截器，添加认证令牌
onMounted(() => {
  // 确保使用环境变量中的API基础URL
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  
  // 如果用户已登录，设置请求头
  if (userStore.isLoggedIn) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${userStore.token}`
  }
  
  // 添加响应拦截器，处理认证错误
  axios.interceptors.response.use(
    response => response,
    error => {
      // 如果是认证错误，清除用户状态并重定向到登录页面
      if (error.response && error.response.status === 401) {
        userStore.logout()
        window.location.href = '/login'
      }
      return Promise.reject(error)
    }
  )
})
</script>

<template>
  <div class="app-container">
    <!-- 使用主布局组件包装需要认证的页面 -->
    <MainLayout v-if="needsLayout">
      <router-view />
    </MainLayout>
    <router-view v-else />
  </div>
</template>

<style>
/* 全局样式 */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

.app-container {
  height: 100vh;
}

/* 报告列表和详情页样式 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.report-summary {
  margin: 20px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  border-left: 4px solid #409EFF;
}

.report-content {
  margin-top: 20px;
}
</style>
