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
  const userPermissions = ref([])

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
        const menuData = response.data || []
        userMenus.value = menuData
        
        // 提取权限编码
        const permissions = extractPermissions(menuData)
        userPermissions.value = permissions
        
        return menuData
      } else {
        throw new Error(response.msg)
      }
    } catch (error) {
      throw error
    }
  }

  // 提取权限编码的递归函数
  const extractPermissions = (menus) => {
    const permissions = []
    
    const traverse = (items) => {
      for (const item of items) {
        // 如果有权限编码，添加到权限列表
        if (item.permissionCode && item.permissionCode.trim()) {
          permissions.push(item.permissionCode.trim())
        }
        // 递归处理子菜单
        if (item.children && item.children.length > 0) {
          traverse(item.children)
        }
      }
    }
    
    traverse(menus)
    
    // 去重
    return [...new Set(permissions)]
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
    userPermissions.value = []
    removeToken()
    router.push('/login')
  }

  return {
    token,
    userInfo,
    userMenus,
    userPermissions,
    login,
    getUserInfo,
    getUserMenus,
    checkLogin,
    logout
  }
}) 