import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi } from '@/api/user'
import { getUserInfo as getUserInfoApi, getUserMenus as getUserMenusApi } from '@/api/system/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(getToken())
  const userInfo = ref({})
  const userMenus = ref([])

  // 登录
  const login = async (loginForm) => {
    try {
      const response = await loginApi(loginForm)
      if (response.code === 0) {
        token.value = response.data.token
        setToken(response.data.token)
        return response
      } else {
        throw new Error(response.msg)
      }
    } catch (error) {
      throw error
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const response = await getUserInfoApi()
      if (response.code === 0) {
        userInfo.value = response.data
        return response.data
      } else {
        throw new Error(response.msg)
      }
    } catch (error) {
      throw error
    }
  }

  // 获取用户菜单
  const getUserMenus = async () => {
    try {
      const response = await getUserMenusApi()
      if (response.code === 0) {
        userMenus.value = response.data || []
        return response.data
      } else {
        throw new Error(response.msg)
      }
    } catch (error) {
      throw error
    }
  }

  // 检查登录状态
  const checkLogin = () => {
    const savedToken = getToken()
    if (savedToken && savedToken.trim()) {
      token.value = savedToken
      // 如果有token但没有用户信息，获取用户信息和菜单
      if (!userInfo.value.id) {
        getUserInfo().then(() => {
          getUserMenus()
        }).catch(() => {
          // 如果获取失败，说明token无效，清空状态
          logout()
        })
      }
    }
  }

  // 退出登录
  const logout = () => {
    token.value = ''
    userInfo.value = {}
    userMenus.value = []
    removeToken()
    router.push('/login')
  }

  return {
    token,
    userInfo,
    userMenus,
    login,
    getUserInfo,
    getUserMenus,
    checkLogin,
    logout
  }
}) 