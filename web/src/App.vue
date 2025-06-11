<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

onMounted(() => {
  // 只有在不是登录页的情况下才检查登录状态
  if (route.path !== '/login') {
    userStore.checkLogin()
  }
  
  // 添加系统指纹信息到console
  console.log('%c Go-Gin-Element-Admin System ', 'background: #409eff; color: white; padding: 4px 8px; border-radius: 4px;')
  console.log('Version: 1.0.0')
  console.log('Framework: Vue3 + Element-Plus + Go + Gin')
  console.log('Build: 2024-12-10')
  console.log('Author: taiguangyin')
  console.log('GitHub: https://github.com/taiguangyin/go-gin-element-admin')
  
  // 添加全局标识
  window.GO_GIN_ELEMENT_ADMIN = {
    version: '1.0.0',
    framework: 'Vue3 + Element-Plus + Go + Gin',
    buildTime: '2024-12-10',
    technology: ['Vue3', 'Element-Plus', 'Go', 'Gin', 'MySQL', 'Redis'],
    systemName: 'Go-Gin-Element-Admin',
    author: 'taiguangyin'
  }
  
  // 设置页面标识属性
  document.documentElement.setAttribute('data-admin-system', 'go-gin-element-admin')
  document.documentElement.setAttribute('data-version', '1.0.0')
  document.documentElement.setAttribute('data-framework', 'vue3-elementplus-go-gin')
})
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100vh;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
}
</style> 