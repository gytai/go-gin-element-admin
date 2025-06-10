<template>
  <div class="permission-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>角色权限管理</span>
          <div class="header-buttons">
            <el-button @click="goBack" :icon="ArrowLeft">
              返回
            </el-button>
            <el-button type="primary" @click="savePermissions" :loading="loading">
              保存权限
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="permission-content">
        <div class="role-info">
          <el-descriptions title="角色信息" :column="2" border>
            <el-descriptions-item label="角色ID">{{ roleInfo.authorityId }}</el-descriptions-item>
            <el-descriptions-item label="角色名称">{{ roleInfo.authorityName }}</el-descriptions-item>
            <el-descriptions-item label="默认路由">{{ roleInfo.defaultRouter }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(roleInfo.createdAt) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="menu-tree">
          <h3>菜单权限</h3>
          <el-tree
            ref="menuTreeRef"
            :data="menuTree"
            :props="treeProps"
            show-checkbox
            node-key="ID"
            :default-checked-keys="checkedMenuIds"
            @check="handleMenuCheck"
          >
            <template #default="{ node, data }">
              <span class="custom-tree-node">
                <el-icon v-if="data.icon" class="menu-icon">
                  <component :is="data.icon" />
                </el-icon>
                <span>{{ data.title }}</span>
                <el-tag v-if="data.path" size="small" type="info" class="path-tag">
                  {{ data.path }}
                </el-tag>
              </span>
            </template>
          </el-tree>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getRoleById } from '@/api/system/role'
import { getRoleMenus, assignRoleMenus } from '@/api/system/permission'

const route = useRoute()
const router = useRouter()

const roleId = route.params.id
const loading = ref(false)
const menuTreeRef = ref()

const roleInfo = reactive({
  authorityId: '',
  authorityName: '',
  defaultRouter: '',
  createdAt: ''
})

const menuTree = ref([])
const checkedMenuIds = ref([])

const treeProps = {
  children: 'children',
  label: 'title'
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 获取角色信息
const getRoleInfo = async () => {
  try {
    const response = await getRoleById(roleId)
    Object.assign(roleInfo, response.data)
  } catch (error) {
    ElMessage.error('获取角色信息失败')
    console.error(error)
  }
}

// 获取角色权限
const getRolePermissions = async () => {
  try {
    const response = await getRoleMenus(roleId)
    checkedMenuIds.value = response.data.menuIds || []
    menuTree.value = response.data.menus || []
  } catch (error) {
    ElMessage.error('获取角色权限失败')
    console.error(error)
  }
}

// 处理菜单选择
const handleMenuCheck = (data, checked) => {
  // 这里可以添加额外的逻辑，比如父子节点联动
}

// 保存权限
const savePermissions = async () => {
  try {
    await ElMessageBox.confirm('确定要保存权限设置吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    loading.value = true
    const checkedKeys = menuTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
    const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]

    await assignRoleMenus(roleId, allCheckedKeys)
    ElMessage.success('权限保存成功')
    
    // 重新获取权限数据
    await getRolePermissions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('权限保存失败')
      console.error(error)
    }
  } finally {
    loading.value = false
  }
}

// 返回角色列表
const goBack = () => {
  router.push('/system/role')
}

onMounted(() => {
  getRoleInfo()
  getRolePermissions()
})
</script>

<style lang="scss" scoped>
.permission-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-buttons {
  display: flex;
  gap: 10px;
}

.permission-content {
  padding: 20px 0;
}

.role-info {
  margin-bottom: 30px;
}

.menu-tree {
  margin-top: 20px;
}

.menu-tree h3 {
  margin-bottom: 15px;
  color: #303133;
}

.custom-tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-icon {
  font-size: 16px;
  color: #409eff;
}

.path-tag {
  margin-left: 8px;
}

:deep(.el-tree-node__content) {
  height: 40px;
}
</style> 