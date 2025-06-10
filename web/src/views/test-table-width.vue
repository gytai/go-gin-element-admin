<template>
  <div class="test-container">
    <h2>表格宽度测试</h2>
    
    <!-- 容器宽度指示器 -->
    <div class="width-indicator">
      <div class="indicator-bar">容器宽度 100%</div>
    </div>
    
    <!-- 测试表格 -->
    <el-card>
      <template #header>
        <span>测试表格 - 应该与容器宽度一致</span>
      </template>
      
      <div class="table-container">
        <el-table
          :data="testData"
          class="full-width-table"
          border
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="名称" width="150" />
          <el-table-column prop="type" label="类型" width="120" />
          <el-table-column prop="status" label="状态" width="100" />
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column prop="createTime" label="创建时间" width="180" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default>
              <el-button size="small" type="primary">编辑</el-button>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    
    <!-- 宽度说明 -->
    <div class="width-info">
      <h3>宽度设置说明：</h3>
      <ul>
        <li>容器设置：width: 100%</li>
        <li>表格设置：width: 100% !important</li>
        <li>表格容器：width: 100%, overflow-x: auto</li>
        <li>响应式：大屏幕完全展开，小屏幕允许横向滚动</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const testData = ref([
  {
    id: 1,
    name: '测试项目1',
    type: '类型A',
    status: '正常',
    description: '这是一个测试描述，用来验证表格宽度是否正确适应容器宽度',
    createTime: '2024-01-01 10:00:00'
  },
  {
    id: 2,
    name: '测试项目2',
    type: '类型B',
    status: '异常',
    description: '另一个测试描述，内容稍长一些，用来测试文本溢出处理',
    createTime: '2024-01-02 11:00:00'
  },
  {
    id: 3,
    name: '测试项目3',
    type: '类型C',
    status: '正常',
    description: '第三个测试项目的描述信息',
    createTime: '2024-01-03 12:00:00'
  }
])
</script>

<style scoped>
.test-container {
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}

.width-indicator {
  margin-bottom: 20px;
  border: 2px dashed #409eff;
  padding: 10px;
  background-color: #f0f9ff;
}

.indicator-bar {
  background-color: #409eff;
  color: white;
  padding: 8px;
  text-align: center;
  font-weight: bold;
}

.table-container {
  margin-bottom: 20px;
  width: 100%;
  overflow-x: auto;
  box-sizing: border-box;
  border: 1px solid #e6e6e6; /* 添加边框以便观察 */
}

.full-width-table {
  width: 100% !important;
  min-width: 100%;
  table-layout: auto;
}

/* 确保表格在容器内完全展开 */
.full-width-table :deep(.el-table__body-wrapper) {
  width: 100%;
}

.full-width-table :deep(.el-table__header-wrapper) {
  width: 100%;
}

.full-width-table :deep(.el-table__row) {
  width: 100%;
}

.width-info {
  margin-top: 20px;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.width-info h3 {
  margin-top: 0;
  color: #333;
}

.width-info ul {
  margin: 10px 0;
  padding-left: 20px;
}

.width-info li {
  margin-bottom: 5px;
  color: #666;
}

/* 确保卡片组件也占满容器宽度 */
.test-container :deep(.el-card) {
  width: 100%;
  box-sizing: border-box;
}

.test-container :deep(.el-card__body) {
  width: 100%;
  box-sizing: border-box;
  padding: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .test-container {
    padding: 10px;
  }
  
  .table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }
  
  .full-width-table {
    min-width: 600px; /* 在小屏幕上设置最小宽度，允许横向滚动 */
  }
}

@media (min-width: 769px) {
  .full-width-table {
    width: 100% !important;
    min-width: 100% !important;
  }
}
</style>
