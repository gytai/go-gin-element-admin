import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 配置 NProgress
NProgress.configure({ showSpinner: false })

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    redirect: '/dashboard',
    component: () => import('@/layout/index.vue'),
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '首页',
          icon: 'House'
        }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          title: '个人中心',
          icon: 'User'
        }
      }
    ]
  },
  {
    path: '/system',
    name: 'System',
    component: () => import('@/layout/index.vue'),
    meta: {
      title: '系统设置',
      icon: 'Setting',
      requiresAuth: true
    },
    children: [
      {
        path: 'user',
        name: 'SystemUser',
        component: () => import('@/views/system/user/index.vue'),
        meta: {
          title: '人员管理',
          icon: 'User'
        }
      },
      {
        path: 'role',
        name: 'SystemRole',
        component: () => import('@/views/system/role/index.vue'),
        meta: {
          title: '角色管理',
          icon: 'UserFilled'
        }
      },
      {
        path: 'role/:id/permission',
        name: 'RolePermission',
        component: () => import('@/views/system/role/permission.vue'),
        meta: {
          title: '权限管理',
          icon: 'Lock',
          hidden: true
        }
      },
      {
        path: 'menu',
        name: 'SystemMenu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: {
          title: '菜单管理',
          icon: 'Menu'
        }
      },
      {
        path: 'operation-log',
        name: 'SystemOperationLog',
        component: () => import('@/views/system/operation-log/index.vue'),
        meta: {
          title: '操作日志',
          icon: 'Document'
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  
  const userStore = useUserStore()
  const token = userStore.token
  
  if (to.path === '/login') {
    if (token) {
      next('/')
    } else {
      next()
    }
  } else {
    if (token) {
      if (!userStore.userInfo.id) {
        try {
          await userStore.getUserInfo()
          // 获取用户信息成功后，获取用户菜单
          await userStore.getUserMenus()
          next()
        } catch (error) {
          userStore.logout()
          next('/login')
        }
      } else {
        // 如果用户信息存在但菜单为空，尝试获取菜单
        if (!userStore.userMenus || userStore.userMenus.length === 0) {
          try {
            await userStore.getUserMenus()
          } catch (error) {
            console.warn('获取用户菜单失败:', error)
          }
        }
        next()
      }
    } else {
      next('/login')
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})

export default router 