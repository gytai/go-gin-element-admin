<template>
  <div class="menu-management">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" v-permission="'menu:create'" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增菜单
          </el-button>
        </div>
      </template>

      <!-- 表格区域 -->
      <div class="table-container">
        <el-table
          :data="tableData"
          v-loading="loading"
          row-key="ID"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          :default-expand-all="false"
          class="data-table full-width"
        >
          <el-table-column prop="title" label="菜单名称" min-width="180" />
          <el-table-column label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="getMenuTypeColor(row)">
                {{ getMenuTypeText(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="路由名称" min-width="120" />
          <el-table-column prop="path" label="路由路径" min-width="150" />
          <el-table-column prop="component" label="组件路径" min-width="180">
            <template #default="{ row }">
              <span v-if="row.component" class="component-path">{{ row.component }}</span>
              <span v-else class="empty-text">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="permissionCode" label="权限编码" min-width="150">
            <template #default="{ row }">
              <span v-if="row.permissionCode" class="permission-code">{{ row.permissionCode }}</span>
              <span v-else class="empty-text">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="icon" label="图标" width="80">
            <template #default="{ row }">
              <el-icon v-if="row.icon">
                <component :is="row.icon" />
              </el-icon>
              <span v-else class="empty-text">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="80" />
          <el-table-column prop="hidden" label="是否隐藏" width="100">
            <template #default="{ row }">
              <el-tag :type="row.hidden ? 'danger' : 'success'">
                {{ row.hidden ? '隐藏' : '显示' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-tooltip content="编辑" placement="top">
                <el-button type="primary" size="small" v-permission="'menu:update'" :icon="Edit" circle @click="handleEdit(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="添加子菜单" placement="top">
                <el-button type="success" size="small" v-permission="'menu:create_child'" :icon="Plus" circle @click="handleAddChild(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button type="danger" size="small" v-permission="'menu:delete'" :icon="Delete" circle @click="handleDelete(row)" />
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 菜单表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="menuFormRef"
        :model="menuForm"
        label-width="100px"
      >
        <el-form-item label="菜单类型">
          <el-radio-group v-model="menuForm.menuType">
            <el-radio value="menu">菜单</el-radio>
            <el-radio value="button">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="父级菜单" v-if="!isAddChild">
          <el-tree-select
            v-model="menuForm.parentId"
            :data="menuTreeData"
            :props="treeSelectProps"
            placeholder="请选择父级菜单"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="菜单名称">
          <el-input v-model="menuForm.title" />
        </el-form-item>
        
        <el-form-item label="权限编码" v-if="menuForm.menuType === 'button'">
          <el-input 
            v-model="menuForm.permissionCode" 
            placeholder="如：user:create, user:delete, user:update"
          />
          <div class="form-tip">
            按钮权限编码，用于前端权限控制，建议格式：模块:操作
          </div>
        </el-form-item>
        
        <el-form-item label="路由名称" v-if="menuForm.menuType === 'menu'">
          <el-input v-model="menuForm.name" />
        </el-form-item>
        
        <el-form-item label="路由路径" v-if="menuForm.menuType === 'menu'">
          <el-input v-model="menuForm.path" />
        </el-form-item>
        
        <el-form-item label="组件路径" v-if="menuForm.menuType === 'menu'">
          <el-input v-model="menuForm.component" />
        </el-form-item>
        
        <el-form-item label="图标">
          <IconSelector v-model="menuForm.icon" />
        </el-form-item>
        
        <el-form-item label="排序">
          <el-input-number v-model="menuForm.sort" :min="0" />
        </el-form-item>
        
        <el-form-item label="是否隐藏">
          <el-switch v-model="menuForm.hidden" />
        </el-form-item>
        
        <el-form-item label="是否缓存" v-if="menuForm.menuType === 'menu'">
          <el-switch v-model="menuForm.keepAlive" />
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
import { reactive, ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getMenuList, createMenu, updateMenu, deleteMenu, getMenuTree } from '@/api/system/menu'
import IconSelector from '@/components/IconSelector.vue'

const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const isAddChild = ref(false)
const menuFormRef = ref()

const menuForm = reactive({
  id: null,
  parentId: null,
  title: '',
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  hidden: false,
  keepAlive: false,
  menuType: 'menu',
  permissionCode: ''
})

const tableData = ref([])
const menuTreeData = ref([])

const treeSelectProps = {
  children: 'children',
  label: 'title',
  value: 'ID'
}

const dialogTitle = computed(() => {
  if (isAddChild.value) return '添加子菜单'
  return isEdit.value ? '编辑菜单' : '新增菜单'
})

const loadData = async () => {
  loading.value = true
  try {
    const response = await getMenuTree()
    
    // 直接使用后端返回的数据
    const data = response.data || []
    
    tableData.value = data
    menuTreeData.value = buildTreeSelectData(data)
  } catch (error) {
    console.error('加载菜单数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const buildTreeSelectData = (data) => {
  return data.map(item => ({
    ID: item.ID,
    title: item.title,
    children: item.children ? buildTreeSelectData(item.children) : []
  }))
}

const handleAdd = () => {
  isEdit.value = false
  isAddChild.value = false
  resetForm()
  dialogVisible.value = true
}

const handleAddChild = (row) => {
  isEdit.value = false
  isAddChild.value = true
  resetForm()
  menuForm.parentId = row.ID
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  isAddChild.value = false
  Object.assign(menuForm, {
    ...row,
    id: row.ID,
    // 根据menuType字段判断类型，如果没有则根据组件路径判断
    menuType: row.menuType || ((row.component || row.children?.length) ? 'menu' : 'button'),
    permissionCode: row.permissionCode || ''
  })
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除菜单 "${row.title}" 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteMenu(row.ID)
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
      await updateMenu(menuForm.id, menuForm)
      ElMessage.success('更新成功')
    } else {
      await createMenu(menuForm)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

const resetForm = () => {
  Object.assign(menuForm, {
    id: null,
    parentId: null,
    title: '',
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    hidden: false,
    keepAlive: false,
    menuType: 'menu',
    permissionCode: ''
  })
}

// 获取菜单类型颜色
const getMenuTypeColor = (row) => {
  // 优先使用menuType字段，没有则根据组件路径判断
  const type = row.menuType || (!row.component && !row.children?.length ? 'button' : 'menu')
  return type === 'button' ? 'info' : 'success'
}

// 获取菜单类型文本
const getMenuTypeText = (row) => {
  // 优先使用menuType字段，没有则根据组件路径判断
  const type = row.menuType || (!row.component && !row.children?.length ? 'button' : 'menu')
  return type === 'button' ? '按钮' : '菜单'
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.menu-management {
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

  .table-container {
    .data-table {
      width: 100%;
      
      &.full-width {
        :deep(.el-table__body-wrapper) {
          overflow-x: auto;
        }
      }
      
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
      
      .component-path {
        font-size: 12px;
        color: #606266;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      }
      
      .permission-code {
        font-size: 12px;
        color: #409eff;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        background: #ecf5ff;
        padding: 2px 6px;
        border-radius: 3px;
      }
      
      .empty-text {
        color: #c0c4cc;
        font-size: 12px;
      }
    }
  }
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style> 