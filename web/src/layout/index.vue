<template>
  <div class="app-wrapper">
    <div class="sidebar-container">
      <Sidebar />
    </div>
    <div class="main-container">
      <div class="navbar">
        <Navbar />
      </div>
      <div class="app-main">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
</script>

<style lang="scss" scoped>
.app-wrapper {
  position: relative;
  height: 100vh;
  width: 100vw;
  max-width: 100%;
  display: flex;
  overflow-x: hidden;

  .sidebar-container {
    width: 210px;
    height: 100vh;
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 1001;
    overflow: hidden;
    background: #304156;
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  }

  .main-container {
    min-height: 100vh;
    flex: 1;
    display: flex;
    flex-direction: column;
    margin-left: 210px;
    width: calc(100vw - 210px);
    max-width: calc(100% - 210px);

    .navbar {
      height: 50px;
      overflow: hidden;
      position: relative;
      background: #fff;
      box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
      width: 100%;
    }

    .app-main {
      flex: 1;
      padding: 0px;
      background: #f0f2f5;
      min-height: calc(100vh - 50px);
      width: 100%;
    }
  }
}
</style> 