import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getToken } from '@/utils/auth'

// 创建 axios 实例
const service = axios.create({
  baseURL: '/api', // api 的 base_url
  timeout: 5000 // 请求超时时间
})

// request 拦截器
service.interceptors.request.use(
  config => {
    // 在请求发送之前做一些处理
    const token = getToken()
    if (token) {
      // 让每个请求携带自定义 token 请根据实际情况自行修改
      config.headers['Authorization'] = 'Bearer ' + token
    }
    return config
  },
  error => {
    // 处理请求错误
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response 拦截器
service.interceptors.response.use(
  response => {
    const res = response.data

    // 如果自定义代码不是 0，则判断为错误
    if (res.code !== 0) {
      // 401: 未授权
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        return Promise.reject(new Error('登录已过期，请重新登录'))
      }

      ElMessage({
        message: res.msg || 'Error',
        type: 'error',
        duration: 5 * 1000
      })

      return Promise.reject(new Error(res.msg || 'Error'))
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error) // for debug
    
    // 处理HTTP状态码401
    if (error.response && error.response.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      return Promise.reject(new Error('登录已过期，请重新登录'))
    }
    
    ElMessage({
      message: error.message || '网络错误',
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service 