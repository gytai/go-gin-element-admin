<template>
  <div class="dashboard" style="padding: 20px;">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon user">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.userCount }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon role">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.roleCount }}</div>
              <div class="stat-label">角色总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon menu">
              <el-icon><Menu /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.menuCount }}</div>
              <div class="stat-label">菜单总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统信息</span>
            </div>
          </template>
          <div class="system-info">
            <div class="info-item">
              <span class="label">系统版本：</span>
              <span class="value">{{ systemInfo.systemVersion }}</span>
            </div>
            <div class="info-item">
              <span class="label">Go版本：</span>
              <span class="value">{{ systemInfo.goVersion }}</span>
            </div>
            <div class="info-item">
              <span class="label">Gin版本：</span>
              <span class="value">{{ systemInfo.ginVersion }}</span>
            </div>
            <div class="info-item">
              <span class="label">Vue版本：</span>
              <span class="value">{{ systemInfo.vueVersion }}</span>
            </div>
            <div class="info-item">
              <span class="label">Element Plus：</span>
              <span class="value">{{ systemInfo.elementPlus }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>版本发布说明</span>
            </div>
          </template>
          <div class="system-info">
            <div class="info-item">
              <span class="label">20250610</span>
              <span class="value">发布初版</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { User, UserFilled, Menu, Connection, Setting } from '@element-plus/icons-vue'
import { getDashboardStats, getSystemInfo } from '@/api/dashboard'
import { ElMessage } from 'element-plus'

const stats = reactive({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  onlineCount: 0
})

const systemInfo = reactive({
  systemVersion: 'v1.0.0',
  goVersion: '1.21',
  ginVersion: '1.9.1',
  vueVersion: '3.4.21',
  elementPlus: '2.6.3'
})

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await getDashboardStats()
    if (response.code === 0) {
      Object.assign(stats, response.data)
    } else {
      ElMessage.error('获取统计数据失败：' + response.msg)
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  }
}

// 加载系统信息
const loadSystemInfo = async () => {
  try {
    const response = await getSystemInfo()
    if (response.code === 0) {
      Object.assign(systemInfo, response.data)
    } else {
      ElMessage.error('获取系统信息失败：' + response.msg)
    }
  } catch (error) {
    console.error('获取系统信息失败:', error)
    ElMessage.error('获取系统信息失败')
  }
}

onMounted(() => {
  loadStats()
  loadSystemInfo()
})
</script>

<style lang="scss" scoped>
.dashboard {
  .stat-card {
    .stat-item {
      display: flex;
      align-items: center;
      
      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;
        
        .el-icon {
          font-size: 24px;
          color: #fff;
        }
        
        &.user {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        
        &.role {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }
        
        &.menu {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }
        
        &.online {
          background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        }
      }
      
      .stat-content {
        .stat-number {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          line-height: 1;
        }
        
        .stat-label {
          font-size: 14px;
          color: #909399;
          margin-top: 4px;
        }
      }
    }
  }

  .system-info {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .label {
        color: #606266;
        font-weight: 500;
      }
      
      .value {
        color: #303133;
      }
    }
  }

  .quick-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    
    .el-button {
      height: 50px;
      
      .el-icon {
        margin-right: 8px;
      }
    }
  }

  .card-header {
    font-weight: 500;
    color: #303133;
  }
}
</style> 