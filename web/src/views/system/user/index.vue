<template>
  <div class="user-management">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>人员管理</span>
          <el-button type="primary" v-permission="'user:create'" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增用户
          </el-button>
        </div>
      </template>

      <!-- 搜索区域 -->
      <div class="search-area">
        <el-form :model="searchForm" inline>
          <el-form-item label="用户名">
            <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
          </el-form-item>
          <el-form-item label="昵称">
            <el-input v-model="searchForm.nickName" placeholder="请输入昵称" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetSearch">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格区域 -->
      <div class="table-container">
        <el-table :data="tableData" v-loading="loading" stripe class="data-table">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="nickName" label="昵称" />
          <el-table-column prop="email" label="邮箱" />
          <el-table-column prop="phone" label="手机号" />
          <el-table-column prop="authorityId" label="角色" />
          <el-table-column prop="CreatedAt" label="创建时间">
            <template #default="{ row }">
              {{ formatDate(row.CreatedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-tooltip content="编辑" placement="top">
                <el-button type="primary" size="small" v-permission="'user:update'" :icon="Edit" circle @click="handleEdit(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="重置密码" placement="top">
                <el-button type="warning" size="small" v-permission="'user:reset_password'" :icon="Refresh" circle @click="handleResetPassword(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button type="danger" size="small" v-permission="'user:delete'" :icon="Delete" circle @click="handleDelete(row)" />
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="600px"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        label-width="80px"
      >
        <el-form-item label="用户名">
          <el-input v-model="userForm.username" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="userForm.nickName" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="userForm.email" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="userForm.phone" />
        </el-form-item>
        <el-form-item label="密码" v-if="!isEdit">
          <el-input v-model="userForm.password" type="password" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="userForm.authorityId" placeholder="请选择角色">
            <el-option
              v-for="role in roleOptions"
              :key="role.authorityId"
              :label="role.authorityName"
              :value="role.authorityId"
            />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Edit, Delete } from '@element-plus/icons-vue'
import { getUserList, createUser, updateUser, deleteUser, resetPassword } from '@/api/system/user'
import { getAllRoles } from '@/api/system/role'
import { useUserStore } from '@/stores/user'
import { hasPermission } from '@/utils/permission'

const userStore = useUserStore()

const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const userFormRef = ref()

const searchForm = reactive({
  username: '',
  nickName: ''
})

const userForm = reactive({
  ID: null,
  username: '',
  nickName: '',
  email: '',
  phone: '',
  password: '',
  authorityId: ''
})

const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})
const roleOptions = ref([])

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    const response = await getUserList(params)
    tableData.value = response.data.list || []
    pagination.total = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  try {
    const response = await getAllRoles()
    roleOptions.value = response.data || []
  } catch (error) {
    console.error('加载角色失败:', error)
  }
}

const handleSearch = () => {
  loadData()
}

const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    nickName: ''
  })
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(userForm, {
    ID: null,
    username: '',
    nickName: '',
    email: '',
    phone: '',
    password: '',
    authorityId: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(userForm, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteUser(row.ID)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const handleSubmit = async () => {
  try {
    if (isEdit.value) {
      await updateUser(userForm.ID, userForm)
      ElMessage.success('更新成功')
    } else {
      await createUser(userForm)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

const handleResetPassword = (row) => {
  ElMessageBox.confirm(`确定要重置用户 "${row.username}" 的密码吗？重置后密码为：123456`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await resetPassword(row.ID)
      ElMessage.success('密码重置成功，新密码为：123456')
    } catch (error) {
      ElMessage.error('密码重置失败')
    }
  })
}

onMounted(() => {
  loadData()
  loadRoles()
})
</script>

<style lang="scss" scoped>
.user-management {
  padding: 20px;
  
  .page-card {
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    
    :deep(.el-card__header) {
      background: #ffffff;
      border-bottom: 1px solid #e4e7ed;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 500;
    color: #303133;
    font-size: 16px;
  }

  .search-area {
    margin-bottom: 20px;
    padding: 20px;
    background: #ffffff;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    
    .el-form {
      margin-bottom: 0;
    }
  }

  .table-container {
    .data-table {
      :deep(.el-table__header) {
        background: #ffffff;
        
        th {
          background: #ffffff !important;
          color: #606266;
          font-weight: 500;
        }
      }
      
      :deep(.el-table__body) {
        tr:hover > td {
          background-color: #f5f7fa;
        }
      }
    }
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    padding: 0 20px;
  }
}
</style> 