import request from '@/utils/request'

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/system/user/list',
    method: 'get',
    params
  })
}

// 根据ID获取用户
export function getUserById(id) {
  return request({
    url: `/system/user/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/system/user',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(id, data) {
  return request({
    url: `/system/user/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(id) {
  return request({
    url: `/system/user/${id}`,
    method: 'delete'
  })
}

// 获取当前用户信息
export function getUserInfo() {
  return request({
    url: '/system/user/info',
    method: 'get'
  })
}

// 更新当前用户信息
export function updateUserInfo(data) {
  return request({
    url: '/system/user/info',
    method: 'put',
    data
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/system/user/password',
    method: 'put',
    data
  })
}

// 重置密码
export function resetPassword(id) {
  return request({
    url: `/system/user/${id}`,
    method: 'put',
    data: {
      password: 'RESET_PASSWORD_123456'
    }
  })
}

// 获取用户菜单（根据权限动态返回）
export function getUserMenus() {
  return request({
    url: '/system/user/menus',
    method: 'get'
  })
} 