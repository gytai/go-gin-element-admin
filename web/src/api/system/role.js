import request from '@/utils/request'

// 获取角色列表
export function getRoleList(params) {
  return request({
    url: '/system/role/list',
    method: 'get',
    params
  })
}

// 获取所有角色（下拉选项用）
export function getAllRoles() {
  return request({
    url: '/system/role/all',
    method: 'get'
  })
}

// 根据ID获取角色
export function getRoleById(id) {
  return request({
    url: `/system/role/${id}`,
    method: 'get'
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/system/role',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id, data) {
  return request({
    url: `/system/role/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id) {
  return request({
    url: `/system/role/${id}`,
    method: 'delete'
  })
} 