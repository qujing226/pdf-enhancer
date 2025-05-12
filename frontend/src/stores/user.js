import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || '{}')
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    userInfo: (state) => state.user
  },
  actions: {
    async login(email, password) {
      try {
        const response = await axios.post('/api/v1/login', { email, password })
        const { token, user } = response.data.data
        
        // 保存令牌和用户信息
        this.token = token
        this.user = user
        
        // 存储到本地存储
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        
        // 设置 axios 默认请求头
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        return Promise.resolve(response)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async register(name, email, password) {
      try {
        const response = await axios.post('/api/v1/register', { name, email, password })
        return Promise.resolve(response)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    logout() {
      // 清除状态
      this.token = ''
      this.user = {}
      
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      // 清除请求头
      delete axios.defaults.headers.common['Authorization']
    }
  }
})