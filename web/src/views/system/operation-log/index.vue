<template>
  <div class="operation-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :model="searchForm" :inline="true" class="search-form">
        <el-form-item label="用户名">
          <el-input 
            v-model="searchForm.username" 
            placeholder="请输入用户名" 
            clearable 
            style="width: 180px;"
          />
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select 
            v-model="searchForm.operationType" 
            placeholder="请选择操作类型" 
            clearable
            style="width: 160px;"
          >
            <el-option label="创建" value="CREATE" />
            <el-option label="更新" value="UPDATE" />
            <el-option label="删除" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select 
            v-model="searchForm.method" 
            placeholder="请选择请求方法" 
            clearable
            style="width: 150px;"
          >
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态码">
          <el-select 
            v-model="searchForm.status" 
            placeholder="请选择状态码" 
            clearable
            style="width: 140px;"
          >
            <el-option label="200" :value="200" />
            <el-option label="400" :value="400" />
            <el-option label="401" :value="401" />
            <el-option label="403" :value="403" />
            <el-option label="404" :value="404" />
            <el-option label="500" :value="500" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作时间">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleDateChange"
            style="width: 350px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="loading"
          :data="tableData"
          @selection-change="handleSelectionChange"
          class="full-width-table"
        >
          <el-table-column type="selection" width="55" />
        <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" width="120" />
          <el-table-column prop="operationType" label="操作类型" width="100">
            <template #default="{ row }">
              <el-tag 
                :type="getOperationTypeTagType(row.operationType)"
                size="small"
              >
                {{ getOperationTypeText(row.operationType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="method" label="请求方法" width="100">
            <template #default="{ row }">
              <el-tag 
                :type="getMethodTagType(row.method)"
                size="small"
              >
                {{ row.method }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="path" label="请求路径" width="200" show-overflow-tooltip />
          <el-table-column prop="description" label="操作描述" show-overflow-tooltip />
          <el-table-column prop="status" label="状态码" width="100">
            <template #default="{ row }">
              <el-tag 
                :type="getStatusTagType(row.status)"
                size="small"
              >
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="ip" label="IP地址" width="120" />
          <el-table-column prop="latency" label="耗时(ms)" width="100" />
          <el-table-column prop="operationTime" label="操作时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.operationTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button 
                type="primary" 
                size="small" 
                :icon="View" 
                @click="handleView(row)"
              >
                查看
              </el-button>
            </template>
          </el-table-column>
        </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; text-align: right;"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="操作日志详情" width="80%">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ currentLog.ID }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentLog.username }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getOperationTypeTagType(currentLog.operationType)">
            {{ getOperationTypeText(currentLog.operationType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">
          <el-tag :type="getMethodTagType(currentLog.method)">
            {{ currentLog.method }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ currentLog.path }}</el-descriptions-item>
        <el-descriptions-item label="操作描述">{{ currentLog.description }}</el-descriptions-item>
        <el-descriptions-item label="状态码">
          <el-tag :type="getStatusTagType(currentLog.status)">
            {{ currentLog.status }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ currentLog.latency }}ms</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatDateTime(currentLog.operationTime) }}</el-descriptions-item>
        <el-descriptions-item label="用户代理" :span="2">{{ currentLog.userAgent }}</el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentLog.errorMessage">
          <el-text type="danger">{{ currentLog.errorMessage }}</el-text>
        </el-descriptions-item>
      </el-descriptions>
      
      <el-divider>请求参数</el-divider>
      <el-input
        v-model="currentLog.requestBody"
        type="textarea"
        :rows="6"
        readonly
        placeholder="无请求参数"
      />
      
      <el-divider>响应结果</el-divider>
      <el-input
        v-model="currentLog.responseBody"
        type="textarea"
        :rows="6"
        readonly
        placeholder="无响应数据"
      />
    </el-dialog>

    <!-- 清理日志对话框 -->
    <el-dialog v-model="clearDialogVisible" title="清理日志" width="400px">
      <el-form :model="clearForm" label-width="120px">
        <el-form-item label="保留天数">
          <el-input-number 
            v-model="clearForm.days" 
            :min="1" 
            :max="365" 
            placeholder="请输入保留天数"
            style="width: 100%;"
          />
          <div class="form-tip">删除指定天数之前的日志数据</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="clearDialogVisible = false">取消</el-button>
          <el-button type="danger" @click="confirmClearByDays">确定清理</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 统计信息对话框 -->
    <el-dialog v-model="statsDialogVisible" title="操作日志统计" width="600px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-statistic title="今日操作数" :value="statsData.todayCount" />
        </el-col>
        <el-col :span="12">
          <el-statistic title="总操作数" :value="statsData.totalCount" />
        </el-col>
      </el-row>
      
      <el-divider>操作类型分布</el-divider>
      <div class="stat-types">
        <div v-for="item in statsData.operationTypes" :key="item.type" class="type-item">
          <span>{{ getOperationTypeText(item.type) }}</span>
          <span class="count">{{ item.count }}</span>
        </div>
      </div>
      
      <el-divider>最近7天趋势</el-divider>
      <div class="daily-stats">
        <div v-for="item in statsData.dailyStats" :key="item.date" class="daily-item">
          <div class="date">{{ item.date }}</div>
          <div class="count">{{ item.count }}</div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, View, Refresh, Download, DataAnalysis } from '@element-plus/icons-vue'
import {
  getOperationLogList,
  getOperationLogById,
  deleteOperationLog,
  deleteOperationLogsByIds,
  clearOperationLogs,
  clearOperationLogsByDays,
  getOperationStats,
  exportOperationLogs
} from '@/api/system/operation-log'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const selectedRows = ref([])
const detailDialogVisible = ref(false)
const clearDialogVisible = ref(false)
const statsDialogVisible = ref(false)
const currentLog = ref({})
const dateRange = ref([])
const statsData = ref({})

// 搜索表单
const searchForm = reactive({
  username: '',
  operationType: '',
  method: '',
  status: null,
  startTime: '',
  endTime: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 清理表单
const clearForm = reactive({
  days: 30
})

// 获取操作日志列表
const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    const response = await getOperationLogList(params)
    if (response.code === 0) {
      tableData.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.msg || '获取操作日志列表失败')
    }
  } catch (error) {
    console.error('获取操作日志列表失败:', error)
    ElMessage.error('获取操作日志列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  getList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    username: '',
    operationType: '',
    method: '',
    status: null,
    startTime: '',
    endTime: ''
  })
  dateRange.value = []
  pagination.page = 1
  getList()
}

// 刷新
const handleRefresh = () => {
  getList()
}

// 日期范围变化
const handleDateChange = (dates) => {
  if (dates && dates.length === 2) {
    searchForm.startTime = dates[0]
    searchForm.endTime = dates[1]
  } else {
    searchForm.startTime = ''
    searchForm.endTime = ''
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 查看详情
const handleView = async (row) => {
  try {
    const response = await getOperationLogById(row.ID)
    if (response.code === 0) {
      currentLog.value = response.data
      detailDialogVisible.value = true
    } else {
      ElMessage.error(response.msg || '获取日志详情失败')
    }
  } catch (error) {
    console.error('获取日志详情失败:', error)
    ElMessage.error('获取日志详情失败')
  }
}

// 删除单个日志
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条操作日志吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await deleteOperationLog(row.id)
    if (response.code === 0) {
      ElMessage.success('删除成功')
      getList()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除操作日志失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的日志')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 条操作日志吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const ids = selectedRows.value.map(row => row.id)
    const response = await deleteOperationLogsByIds(ids)
    if (response.code === 0) {
      ElMessage.success('删除成功')
      getList()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除操作日志失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 清理日志
const handleClearByDays = () => {
  clearDialogVisible.value = true
}

// 确认清理
const confirmClearByDays = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除 ${clearForm.days} 天前的所有操作日志吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await clearOperationLogsByDays(clearForm.days)
    if (response.code === 0) {
      ElMessage.success('清理成功')
      clearDialogVisible.value = false
      getList()
    } else {
      ElMessage.error(response.msg || '清理失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清理操作日志失败:', error)
      ElMessage.error('清理失败')
    }
  }
}

// 清空所有日志
const handleClearAll = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有操作日志吗？此操作不可恢复！', '危险操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    })

    const response = await clearOperationLogs()
    if (response.code === 0) {
      ElMessage.success('清空成功')
      getList()
    } else {
      ElMessage.error(response.msg || '清空失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清空操作日志失败:', error)
      ElMessage.error('清空失败')
    }
  }
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  getList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  getList()
}

// 获取操作类型标签类型
const getOperationTypeTagType = (type) => {
  const typeMap = {
    'CREATE': 'success',
    'UPDATE': 'warning',
    'DELETE': 'danger'
  }
  return typeMap[type] || 'info'
}

// 获取操作类型文本
const getOperationTypeText = (type) => {
  const typeMap = {
    'CREATE': '创建',
    'UPDATE': '更新',
    'DELETE': '删除'
  }
  return typeMap[type] || type
}

// 获取请求方法标签类型
const getMethodTagType = (method) => {
  const methodMap = {
    'POST': 'success',
    'PUT': 'warning',
    'DELETE': 'danger',
    'GET': 'info'
  }
  return methodMap[method] || 'info'
}

// 获取状态码标签类型
const getStatusTagType = (status) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 400 && status < 500) return 'warning'
  if (status >= 500) return 'danger'
  return 'info'
}

// 格式化日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 导出操作日志
const handleExport = async () => {
  try {
    const params = {
      ...searchForm
    }
    const response = await exportOperationLogs(params)
    
    // 创建下载链接
    const blob = new Blob([response], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `operation_logs_${new Date().getTime()}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出操作日志失败:', error)
    ElMessage.error('导出失败')
  }
}

// 查看统计信息
const handleViewStats = async () => {
  try {
    const response = await getOperationStats()
    if (response.code === 0) {
      statsData.value = response.data
      statsDialogVisible.value = true
    } else {
      ElMessage.error(response.msg || '获取统计信息失败')
    }
  } catch (error) {
    console.error('获取统计信息失败:', error)
    ElMessage.error('获取统计信息失败')
  }
}

// 页面加载时获取数据
onMounted(() => {
  getList()
})
</script>

<style scoped>
.operation-log {
  padding: 20px;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
}

/* 确保卡片组件也占满容器宽度 */
.operation-log :deep(.el-card) {
  width: 100%;
  box-sizing: border-box;
}

.operation-log :deep(.el-card__body) {
  width: 100%;
  box-sizing: border-box;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
  width: 100%;
  box-sizing: border-box;
}

/* 搜索表单项自适应布局 */
.search-form :deep(.el-form-item) {
  margin-right: 16px;
  margin-bottom: 16px;
}

.search-form :deep(.el-form-item__content) {
  flex: 1;
}

.operation-buttons {
  margin-bottom: 20px;
}

.operation-buttons .el-button {
  margin-right: 10px;
}

.table-container {
  margin-bottom: 20px;
  width: 100%;
  overflow-x: auto;
  box-sizing: border-box;
}

.full-width-table {
  width: 100% !important;
  table-layout: auto;
}

/* 确保表格在容器内完全展开 */
.full-width-table :deep(.el-table__body-wrapper) {
  width: 100%;
}

.full-width-table :deep(.el-table__header-wrapper) {
  width: 100%;
}

/* 确保表格行宽度一致 */
.full-width-table :deep(.el-table__row) {
  width: 100%;
}

/* 修改表格标题背景色 */
.full-width-table :deep(.el-table__header) {
  background-color: #ffffff;
}

.full-width-table :deep(.el-table__header th) {
  background-color: #ffffff !important;
  color: #606266;
  font-weight: 500;
  border-bottom: 1px solid #ebeef5;
}

.full-width-table :deep(.el-table__header th.el-table__cell) {
  background: #ffffff !important;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
}

.stat-types {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.type-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.type-item .count {
  font-weight: bold;
  color: #333;
}

.daily-stats {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.daily-item {
  flex: 1;
  text-align: center;
  padding: 12px;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  background-color: #f9f9f9;
}

.daily-item .date {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
}

.daily-item .count {
  font-size: 18px;
  font-weight: bold;
  color: #409eff;
}

.form-tip {
  font-size: 0.8em;
  color: #909399;
  margin-top: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .operation-log {
    padding: 10px;
  }

  .search-form {
    flex-direction: column;
  }

  .search-form :deep(.el-form-item) {
    margin-right: 0;
    width: 100%;
  }

  .operation-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .operation-buttons .el-button {
    margin-right: 0;
    flex: 1;
    min-width: 120px;
  }

  .table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .full-width-table {
    min-width: 800px; /* 在小屏幕上设置最小宽度，允许横向滚动 */
  }
}

/* 确保表格在大屏幕上完全展开 */
@media (min-width: 769px) {
  .full-width-table {
    width: 100% !important;
    min-width: 100% !important;
  }
}
</style>
