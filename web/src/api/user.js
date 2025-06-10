import request from '@/utils/request'

// 登录
export function login(data) {
  return request({
    url: '/base/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/getUserInfo',
    method: 'get'
  })
}

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/user/getUserList',
    method: 'get',
    params
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(data) {
  return request({
    url: '/user/setUserInfo',
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(data) {
  return request({
    url: '/user/deleteUser',
    method: 'delete',
    data
  })
}

// 重置密码
export function resetPassword(data) {
  return request({
    url: '/user/resetPassword',
    method: 'post',
    data
  })
} 