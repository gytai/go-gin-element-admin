<template>
  <div class="sidebar">
    <div class="logo">
      <h2>Admin脚手架</h2>
    </div>
    <el-menu
      :default-active="activeMenu"
      class="el-menu-vertical"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#ffffff"
      :router="true"
    >
      <template v-for="menu in menuList" :key="menu.ID">
        <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
          <template #title>
            <el-icon v-if="menu.icon">
              <component :is="menu.icon" />
            </el-icon>
            <span>{{ menu.title }}</span>
          </template>
          <template v-for="child in menu.children" :key="child.ID">
            <el-menu-item :index="child.path">
              <el-icon v-if="child.icon">
                <component :is="child.icon" />
              </el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </template>
        </el-sub-menu>
        <el-menu-item v-else :index="menu.path">
          <el-icon v-if="menu.icon">
            <component :is="menu.icon" />
          </el-icon>
          <span>{{ menu.title }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, ref } from 'vue'
import { useRoute } from 'vue-router'
import { House, User, Setting, UserFilled, Menu } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()
const menuLoading = ref(false)

const activeMenu = computed(() => {
  const { path } = route
  return path
})

// 从store中获取菜单
const menuList = computed(() => {
  return userStore.userMenus || []
})

// 安全地获取菜单
const safeGetUserMenus = async () => {
  if (menuLoading.value) return // 避免重复请求
  
  try {
    menuLoading.value = true
    await userStore.getUserMenus()
  } catch (error) {
    console.warn('获取用户菜单失败:', error)
  } finally {
    menuLoading.value = false
  }
}

// 监听路由变化，确保菜单始终存在
watch(() => route.path, () => {
  // 路由变化时检查菜单是否存在
  if (userStore.token && (!userStore.userMenus || userStore.userMenus.length === 0)) {
    console.log('路由变化检测到菜单缺失，重新获取菜单')
    safeGetUserMenus()
  }
}, { immediate: true })

// 监听用户登录状态
watch(() => userStore.token, (newToken) => {
  if (newToken && (!userStore.userMenus || userStore.userMenus.length === 0)) {
    console.log('检测到用户已登录但菜单为空，获取菜单')
    safeGetUserMenus()
  }
})

onMounted(() => {
  // 组件挂载时检查菜单
  if (userStore.token && (!userStore.userMenus || userStore.userMenus.length === 0)) {
    console.log('侧边栏组件挂载时检测到菜单缺失，获取菜单')
    safeGetUserMenus()
  }
})
</script>

<style lang="scss" scoped>
.sidebar {
  height: 100%;
  background-color: #304156;
  transition: width 0.28s;
  min-height: 100vh;
  overflow: hidden;
  width: 210px;
  min-width: 210px;
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

.logo {
  height: 60px;
  background-color: #2b3a4b;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #3c4b5c;
  
  h2 {
    color: #ffffff;
    font-size: 24px;
    font-weight: bold;
    margin: 0;
  }
}

.el-menu-vertical {
  border-right: none;
  height: calc(100vh - 60px);
  overflow-y: auto;
}

:deep(.el-menu-item) {
  &.is-active {
    background-color: #409eff !important;
    color: #ffffff !important;
  }
}

:deep(.el-sub-menu__title) {
  color: #bfcbd9 !important;
  
  &:hover {
    background-color: #48576a !important;
    color: #ffffff !important;
  }
}

:deep(.el-menu-item) {
  &:hover {
    background-color: #48576a !important;
    color: #ffffff !important;
  }
}
</style> 