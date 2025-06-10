import request from '@/utils/request'

// 获取角色的菜单权限
export function getRoleMenus(roleId) {
  return request({
    url: `/system/role/${roleId}/menus`,
    method: 'get'
  })
}

// 给角色分配菜单权限
export function assignRoleMenus(roleId, menuIds) {
  return request({
    url: `/system/role/${roleId}/menus`,
    method: 'post',
    data: {
      menuIds
    }
  })
} 