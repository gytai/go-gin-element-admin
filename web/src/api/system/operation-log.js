import request from '@/utils/request'

// 获取操作日志列表
export function getOperationLogList(params) {
  return request({
    url: '/system/operation-log/list',
    method: 'get',
    params
  })
}

// 获取操作日志详情
export function getOperationLogById(id) {
  return request({
    url: `/system/operation-log/${id}`,
    method: 'get'
  })
}

// 删除操作日志
export function deleteOperationLog(id) {
  return request({
    url: `/system/operation-log/${id}`,
    method: 'delete'
  })
}

// 批量删除操作日志
export function deleteOperationLogsByIds(ids) {
  return request({
    url: '/system/operation-log/batch',
    method: 'delete',
    data: ids
  })
}

// 清空操作日志
export function clearOperationLogs() {
  return request({
    url: '/system/operation-log/clear',
    method: 'delete'
  })
}

// 清理指定天数前的操作日志
export function clearOperationLogsByDays(days) {
  return request({
    url: '/system/operation-log/clear-by-days',
    method: 'delete',
    data: { days }
  })
}

// 获取操作统计信息
export function getOperationStats() {
  return request({
    url: '/system/operation-log/stats',
    method: 'get'
  })
}

// 导出操作日志
export function exportOperationLogs(params) {
  return request({
    url: '/system/operation-log/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}
