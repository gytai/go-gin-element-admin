import request from '@/utils/request'

// 获取仪表板统计数据
export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  })
}

// 获取系统信息
export function getSystemInfo() {
  return request({
    url: '/dashboard/systemInfo',
    method: 'get'
  })
} 