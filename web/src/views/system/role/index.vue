<template>
  <div class="role-management">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" v-permission="'role:create'" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </div>
      </template>

      <!-- 表格区域 -->
      <div class="table-container">
        <el-table :data="tableData" v-loading="loading" stripe class="data-table">
          <el-table-column prop="authorityId" label="角色ID" width="100" />
          <el-table-column prop="authorityName" label="角色名称" />
          <el-table-column prop="authorityCode" label="角色编码" />
          <el-table-column prop="parentName" label="父角色" width="150">
            <template #default="{ row }">
              {{ getParentRoleName(row.parentId) || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="defaultRouter" label="默认路由" />
          <el-table-column prop="createdAt" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="{ row }">
              <el-tooltip content="编辑" placement="top">
                <el-button type="primary" size="small" v-permission="'role:update'" :icon="Edit" circle @click="handleEdit(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="分配菜单权限" placement="top">
                <el-button type="success" size="small" v-permission="'role:assign_menu'" :icon="Lock" circle @click="handleSetPermission(row)" style="margin-right: 8px;" />
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button type="danger" size="small" v-permission="'role:delete'" :icon="Delete" circle @click="handleDelete(row)" />
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 角色表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新增角色'"
      width="500px"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleFormRules"
        label-width="100px"
      >
        <el-form-item label="角色ID" v-if="isEdit">
          <el-input v-model="roleForm.authorityId" disabled />
        </el-form-item>
        <el-form-item label="角色名称" prop="authorityName">
          <el-input v-model="roleForm.authorityName" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="authorityCode">
          <el-input v-model="roleForm.authorityCode" placeholder="请输入角色编码，如：admin、user" />
        </el-form-item>
        <el-form-item label="父角色">
          <el-select v-model="roleForm.parentId" placeholder="请选择父角色" clearable>
            <el-option
              v-for="role in availableParentRoles"
              :key="role.authorityId"
              :label="role.authorityName"
              :value="role.authorityId"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="默认路由" prop="defaultRouter">
          <el-input v-model="roleForm.defaultRouter" placeholder="请输入默认路由" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 权限设置对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="设置权限"
      width="600px"
    >
      <div v-if="currentRole && currentRole.parentId" class="permission-notice">
        <el-alert
          title="提示：子角色的权限范围不能超过父角色"
          type="warning"
          :closable="false"
          show-icon
        />
      </div>
      <el-tree
        ref="permissionTreeRef"
        :data="menuTreeData"
        :props="treeProps"
        show-checkbox
        node-key="ID"
        :default-checked-keys="checkedKeys"
      />
      
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermission">
          保存权限
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Lock, Delete } from '@element-plus/icons-vue'
import { getRoleList, createRole, updateRole, deleteRole, getAllRoles } from '@/api/system/role'
import { getRoleMenus, assignRoleMenus } from '@/api/system/permission'
import { getMenuTree } from '@/api/system/menu'

const router = useRouter()
const loading = ref(false)
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const roleFormRef = ref()
const permissionTreeRef = ref()
const currentRole = ref(null)

const roleForm = reactive({
  authorityId: '',
  authorityName: '',
  authorityCode: '',
  parentId: null,
  defaultRouter: 'dashboard'
})

const roleFormRules = {
  authorityName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ],
  authorityCode: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '角色编码只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  defaultRouter: [
    { required: true, message: '请输入默认路由', trigger: 'blur' }
  ]
}

const tableData = ref([])
const allRoles = ref([])
const menuTreeData = ref([])
const checkedKeys = ref([])

const treeProps = {
  children: 'children',
  label: 'title'
}

// 可选择的父角色列表（排除自己和子级角色）
const availableParentRoles = computed(() => {
  if (!isEdit.value) {
    return allRoles.value
  }
  
  // 编辑时排除自己和自己的子角色
  return allRoles.value.filter(role => {
    if (role.authorityId === roleForm.authorityId) {
      return false
    }
    // 检查是否是当前角色的子角色
    return !isChildRole(role.authorityId, roleForm.authorityId)
  })
})

const isChildRole = (roleId, parentId) => {
  const role = allRoles.value.find(r => r.authorityId === roleId)
  if (!role || !role.parentId) return false
  
  if (role.parentId === parentId) return true
  return isChildRole(role.parentId, parentId)
}

const getParentRoleName = (parentId) => {
  if (!parentId) return null
  const parentRole = allRoles.value.find(role => role.authorityId === parentId)
  return parentRole ? parentRole.authorityName : null
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleString()
}

// 获取叶子节点ID（避免父子节点选中冲突）
const getLeafNodeIds = (treeData, assignedIds) => {
  const leafIds = []
  const hasParent = new Set()
  
  // 递归遍历树，收集所有有子节点的节点ID
  const collectParentIds = (nodes) => {
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        hasParent.add(node.ID)
        collectParentIds(node.children)
      }
    })
  }
  
  collectParentIds(treeData)
  
  // 只返回叶子节点（没有子节点的节点）的ID
  assignedIds.forEach(id => {
    if (!hasParent.has(id)) {
      leafIds.push(id)
    }
  })
  
  return leafIds
}

const loadData = async () => {
  loading.value = true
  try {
    const response = await getRoleList()
    tableData.value = response.data.list || []
  } catch (error) {
    console.error('加载角色数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const loadAllRoles = async () => {
  try {
    const response = await getAllRoles()
    allRoles.value = response.data || []
  } catch (error) {
    console.error('加载所有角色失败:', error)
  }
}

const loadMenuTree = async () => {
  try {
    const response = await getMenuTree()
    menuTreeData.value = response.data || []
  } catch (error) {
    console.error('加载菜单树失败:', error)
    ElMessage.error('加载菜单树失败')
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(roleForm, {
    authorityId: null,
    authorityName: '',
    authorityCode: '',
    parentId: null,
    defaultRouter: 'dashboard'
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(roleForm, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除角色 "${row.authorityName}" 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteRole(row.authorityId)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const handleSetPermission = async (row) => {
  currentRole.value = row
  
  try {
    // 获取完整的菜单树数据（用于显示所有菜单选项）
    const menuResponse = await getMenuTree()
    menuTreeData.value = menuResponse.data || []
    
    // 获取角色已有的权限ID列表
    const permissionResponse = await getRoleMenus(row.authorityId)
    const permissionData = permissionResponse.data || {}
    
    // 设置默认选中的菜单ID（只设置叶子节点，避免父子节点冲突）
    const assignedMenuIds = permissionData.menuIds || []
    checkedKeys.value = getLeafNodeIds(menuTreeData.value, assignedMenuIds)
    
    permissionDialogVisible.value = true
  } catch (error) {
    console.error('加载权限数据失败:', error)
    ElMessage.error('加载权限数据失败')
  }
}

const handleSubmit = async () => {
  if (!roleFormRef.value) return
  
  const valid = await roleFormRef.value.validate().catch(() => false)
  if (!valid) return
  
  try {
    if (isEdit.value) {
      await updateRole(roleForm.authorityId, roleForm)
      ElMessage.success('更新成功')
    } else {
      const createData = {
        authorityName: roleForm.authorityName,
        authorityCode: roleForm.authorityCode,
        parentId: roleForm.parentId,
        defaultRouter: roleForm.defaultRouter
      }
      await createRole(createData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
    loadAllRoles() // 重新加载所有角色数据
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

const handleSavePermission = async () => {
  if (!currentRole.value) {
    ElMessage.error('未选择角色')
    return
  }
  
  try {
    const checkedNodes = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedNodes = permissionTreeRef.value.getHalfCheckedKeys()
    
    // 合并全选和半选的节点
    const allSelectedMenuIds = [...checkedNodes, ...halfCheckedNodes]
    
    await assignRoleMenus(currentRole.value.authorityId, allSelectedMenuIds)
    ElMessage.success('权限设置成功')
    permissionDialogVisible.value = false
  } catch (error) {
    console.error('保存权限失败:', error)
    ElMessage.error('保存权限失败')
  }
}

onMounted(() => {
  loadData()
  loadAllRoles()
})
</script>

<style lang="scss" scoped>
.role-management {
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

  .permission-notice {
    margin-bottom: 20px;
  }
}
</style> 