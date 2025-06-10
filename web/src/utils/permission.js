import { useUserStore } from '@/stores/user'

/**
 * 检查用户是否有指定权限
 * @param {string|Array} permission - 权限编码或权限编码数组
 * @param {string} mode - 检查模式：'some'(任一权限) 或 'every'(全部权限)，默认'some'
 * @returns {boolean} 是否有权限
 */
export function hasPermission(permission, mode = 'some') {
  const userStore = useUserStore()
  const userPermissions = userStore.userPermissions || []
  const userInfo = userStore.userInfo || {}
  
  if (!permission) return true
  
  // 详细的调试信息
  console.log('=== 权限检查开始 ===')
  console.log('检查权限:', permission)
  console.log('用户信息:', {
    username: userInfo.username,
    authorityId: userInfo.authority?.authorityId,
    id: userInfo.id
  })
  console.log('用户权限列表:', userPermissions)
  console.log('权限列表长度:', userPermissions.length)
  
  // 管理员检查
  const isAdmin = userInfo?.username === 'admin' || userInfo?.authority?.authorityId === 888
  console.log('是否管理员:', isAdmin)
  
  if (isAdmin) {
    console.log('管理员用户，直接返回true')
    return true
  }
  
  // 如果是字符串，转为数组
  const permissions = Array.isArray(permission) ? permission : [permission]
  console.log('需要检查的权限:', permissions)
  
  let hasAuth = false
  if (mode === 'every') {
    // 检查是否拥有所有权限
    hasAuth = permissions.every(perm => {
      const has = userPermissions.includes(perm)
      console.log(`权限 ${perm}:`, has ? '✓' : '✗')
      return has
    })
  } else {
    // 检查是否拥有任一权限
    hasAuth = permissions.some(perm => {
      const has = userPermissions.includes(perm)
      console.log(`权限 ${perm}:`, has ? '✓' : '✗')
      return has
    })
  }
  
  console.log('最终结果:', hasAuth ? '有权限' : '无权限')
  console.log('=== 权限检查结束 ===')
  
  return hasAuth
}

/**
 * 权限指令
 * 用法：v-permission="'user:create'" 或 v-permission="['user:create', 'user:update']"
 */
export const permissionDirective = {
  mounted(el, binding) {
    const { value } = binding
    
    if (!hasPermission(value)) {
      // 如果没有权限，隐藏元素
      el.style.display = 'none'
      // 或者完全移除元素
      // el.parentNode && el.parentNode.removeChild(el)
    }
  },
  
  updated(el, binding) {
    const { value } = binding
    
    if (!hasPermission(value)) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}

/**
 * 角色权限指令
 * 用法：v-role="'admin'" 或 v-role="['admin', 'user']"
 */
export const roleDirective = {
  mounted(el, binding) {
    const { value } = binding
    const userStore = useUserStore()
    const userRoles = userStore.userInfo?.roles || []
    
    if (!value) return
    
    const roles = Array.isArray(value) ? value : [value]
    const hasRole = roles.some(role => userRoles.includes(role))
    
    if (!hasRole) {
      el.style.display = 'none'
    }
  },
  
  updated(el, binding) {
    const { value } = binding
    const userStore = useUserStore()
    const userRoles = userStore.userInfo?.roles || []
    
    if (!value) return
    
    const roles = Array.isArray(value) ? value : [value]
    const hasRole = roles.some(role => userRoles.includes(role))
    
    if (!hasRole) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}

/**
 * 超级管理员检查
 * @returns {boolean} 是否为超级管理员
 */
export function isSuperAdmin() {
  const userStore = useUserStore()
  return userStore.userInfo?.username === 'admin' || userStore.userInfo?.authorityId === 888
}

/**
 * 权限检查装饰器（用于组合式API）
 * @param {string|Array} permission - 权限编码
 * @param {Function} callback - 有权限时执行的回调
 * @param {Function} fallback - 无权限时执行的回调
 */
export function withPermission(permission, callback, fallback) {
  if (hasPermission(permission)) {
    return callback && callback()
  } else {
    return fallback && fallback()
  }
} 