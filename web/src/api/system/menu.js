import request from '@/utils/request'

// 获取菜单列表
export function getMenuList(params) {
  return request({
    url: '/system/menu/list',
    method: 'get',
    params
  })
}

// 获取菜单树结构
export function getMenuTree() {
  return request({
    url: '/system/menu/tree',
    method: 'get'
  })
}

// 根据ID获取菜单
export function getMenuById(id) {
  return request({
    url: `/system/menu/${id}`,
    method: 'get'
  })
}

// 创建菜单
export function createMenu(data) {
  return request({
    url: '/system/menu',
    method: 'post',
    data
  })
}

// 更新菜单
export function updateMenu(id, data) {
  return request({
    url: `/system/menu/${id}`,
    method: 'put',
    data
  })
}

// 删除菜单
export function deleteMenu(id) {
  return request({
    url: `/system/menu/${id}`,
    method: 'delete'
  })
} 