// 动态路由工具
import router from '@/router'

// 组件映射表，用于动态导入组件
const componentMap = {
  // 布局组件
  'layout/index.vue': () => import('@/layout/index.vue'),
  
  // 页面组件
  'views/dashboard/index.vue': () => import('@/views/dashboard/index.vue'),
  'views/profile/index.vue': () => import('@/views/profile/index.vue'),
  'views/system/user/index.vue': () => import('@/views/system/user/index.vue'),
  'views/system/role/index.vue': () => import('@/views/system/role/index.vue'),
  'views/system/role/permission.vue': () => import('@/views/system/role/permission.vue'),
  'views/system/menu/index.vue': () => import('@/views/system/menu/index.vue'),
  'views/system/operation-log/index.vue': () => import('@/views/system/operation-log/index.vue'),
  
  // 错误页面
  'views/error/404.vue': () => import('@/views/error/404.vue'),
}

/**
 * 将后端菜单数据转换为前端路由格式
 */
export function transformMenusToRoutes(menus) {
  const routes = []
  
  for (const menu of menus) {
    // 只处理菜单类型，跳过按钮类型
    if (menu.menuType === 'button' || menu.hidden) {
      continue
    }
    
    const route = {
      path: menu.path,
      name: menu.name,
      component: getComponent(menu.component),
      meta: {
        title: menu.title,
        icon: menu.icon,
        keepAlive: menu.keepAlive,
        requiresAuth: true
      }
    }
    
    // 处理子路由
    if (menu.children && menu.children.length > 0) {
      const childRoutes = transformMenusToRoutes(menu.children)
      if (childRoutes.length > 0) {
        route.children = childRoutes
      }
    }
    
    routes.push(route)
  }
  
  return routes
}

/**
 * 根据组件路径获取组件
 */
function getComponent(componentPath) {
  if (!componentPath) {
    return () => import('@/views/error/404.vue')
  }
  
  // 从组件映射表中获取组件
  const component = componentMap[componentPath]
  if (component) {
    return component
  }
  
  // 如果映射表中没有，尝试动态导入
  console.warn(`组件 ${componentPath} 未在映射表中找到，尝试动态导入`)
  return () => import(`@/${componentPath}`).catch(() => {
    console.error(`无法加载组件: ${componentPath}`)
    return import('@/views/error/404.vue')
  })
}

/**
 * 动态添加路由
 */
export function addDynamicRoutes(menus) {
  const routes = transformMenusToRoutes(menus)
  
  // 清除之前添加的动态路由（如果需要的话）
  // 注意：Vue Router 4 没有直接的清除路由方法，需要记录添加的路由名称
  
  // 添加动态路由
  routes.forEach(route => {
    try {
      router.addRoute(route)
      console.log(`动态添加路由: ${route.path}`)
    } catch (error) {
      console.error(`添加路由失败: ${route.path}`, error)
    }
  })
  
  return routes
}

/**
 * 重置动态路由
 */
export function resetDynamicRoutes() {
  // 由于 Vue Router 4 没有直接的移除路由方法
  // 我们需要重新创建 router 实例或者记录动态路由进行管理
  // 这里采用简单的方式：重新加载页面来重置路由
  console.log('重置动态路由')
} 