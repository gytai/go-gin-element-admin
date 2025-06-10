<template>
  <div class="permission-test">
    <el-card>
      <template #header>
        <h2>权限调试测试页面</h2>
      </template>
      
      <div class="test-content">
        <!-- 用户信息 -->
        <el-descriptions title="用户信息" :column="2" border>
          <el-descriptions-item label="用户ID">{{ userInfo.id || '未获取' }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ userInfo.username || '未获取' }}</el-descriptions-item>
          <el-descriptions-item label="昵称">{{ userInfo.nickName || '未获取' }}</el-descriptions-item>
          <el-descriptions-item label="角色ID">{{ userInfo.authority?.authorityId || '未获取' }}</el-descriptions-item>
          <el-descriptions-item label="角色名">{{ userInfo.authority?.authorityName || '未获取' }}</el-descriptions-item>
        </el-descriptions>

        <!-- 权限信息 -->
        <div class="section">
          <h3>权限列表</h3>
          <div class="permission-list">
            <el-tag 
              v-for="permission in userPermissions" 
              :key="permission"
              type="success"
              class="permission-tag"
            >
              {{ permission }}
            </el-tag>
            <el-tag v-if="userPermissions.length === 0" type="warning">
              暂无权限
            </el-tag>
          </div>
          <p>权限总数: {{ userPermissions.length }}</p>
        </div>

        <!-- 菜单信息 -->
        <div class="section">
          <h3>用户菜单</h3>
          <el-tree :data="userMenus" :props="{ children: 'children', label: 'title' }" />
        </div>

        <!-- 权限测试按钮 -->
        <div class="section">
          <h3>权限测试</h3>
          <div class="test-buttons">
            <el-button 
              v-permission="'user:create'" 
              type="primary"
              @click="testMessage('user:create')"
            >
              新增用户按钮 (user:create)
            </el-button>
            
            <el-button 
              v-permission="'user:update'" 
              type="success"
              @click="testMessage('user:update')"
            >
              编辑用户按钮 (user:update)
            </el-button>
            
            <el-button 
              v-permission="'user:delete'" 
              type="danger"
              @click="testMessage('user:delete')"
            >
              删除用户按钮 (user:delete)
            </el-button>
            
            <el-button 
              v-permission="'role:create'" 
              type="warning"
              @click="testMessage('role:create')"
            >
              新增角色按钮 (role:create)
            </el-button>
          </div>
        </div>

        <!-- 手动测试 -->
        <div class="section">
          <h3>手动权限测试</h3>
          <el-input 
            v-model="testPermission" 
            placeholder="输入权限编码，如: user:create"
            style="width: 300px; margin-right: 10px;"
          />
          <el-button @click="manualTest" type="primary">测试权限</el-button>
          <p v-if="testResult !== null" :class="testResult ? 'text-green' : 'text-red'">
            测试结果: {{ testResult ? '有权限' : '无权限' }}
          </p>
        </div>

        <!-- 刷新按钮 -->
        <div class="section">
          <el-button @click="refreshData" type="info">刷新数据</el-button>
          <el-button @click="clearConsole" type="warning">清空控制台</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { hasPermission } from '@/utils/permission'

const userStore = useUserStore()
const testPermission = ref('')
const testResult = ref(null)

// 计算属性
const userInfo = computed(() => userStore.userInfo || {})
const userPermissions = computed(() => userStore.userPermissions || [])
const userMenus = computed(() => userStore.userMenus || [])

// 测试消息
const testMessage = (permission) => {
  ElMessage.success(`${permission} 按钮可见，说明有权限`)
}

// 手动测试
const manualTest = () => {
  if (!testPermission.value) {
    ElMessage.warning('请输入权限编码')
    return
  }
  
  testResult.value = hasPermission(testPermission.value)
}

// 刷新数据
const refreshData = async () => {
  try {
    await userStore.getUserInfo()
    await userStore.getUserMenus()
    ElMessage.success('数据刷新成功')
  } catch (error) {
    ElMessage.error('数据刷新失败')
  }
}

// 清空控制台
const clearConsole = () => {
  console.clear()
  ElMessage.info('控制台已清空')
}
</script>

<style scoped>
.permission-test {
  padding: 20px;
}

.test-content {
  padding: 20px 0;
}

.section {
  margin: 30px 0;
}

.section h3 {
  margin-bottom: 15px;
  color: #303133;
}

.permission-list {
  margin-bottom: 10px;
}

.permission-tag {
  margin: 5px;
}

.test-buttons {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.text-green {
  color: #67c23a;
  font-weight: bold;
}

.text-red {
  color: #f56c6c;
  font-weight: bold;
}
</style> 