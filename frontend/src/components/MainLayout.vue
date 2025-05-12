<template>
  <div class="main-layout">
    <!-- 顶部导航栏 -->
    <el-header class="header">
      <div class="logo">
        <h2>AI增强报告系统</h2>
      </div>
      <div class="user-info" v-if="userStore.isLoggedIn">
        <span>{{ userStore.userInfo.name }}</span>
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link">
            <el-avatar :size="32" :src="userAvatar">{{ userInitials }}</el-avatar>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <!-- 主要内容区域 -->
    <el-main class="main-content">
      <slot></slot>
    </el-main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

// 计算用户头像或初始字母
const userInitials = computed(() => {
  if (userStore.userInfo && userStore.userInfo.name) {
    return userStore.userInfo.name.substring(0, 1).toUpperCase()
  }
  return 'U'
})

const userAvatar = computed(() => {
  return userStore.userInfo.avatar || ''
})

// 处理下拉菜单命令
const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    ElMessage.success('退出登录成功')
    router.push('/login')
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #409EFF;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo h2 {
  margin: 0;
  font-size: 1.5rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
  overflow-y: auto;
}
</style>