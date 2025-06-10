<template>
  <div class="permission-demo">
    <el-card>
      <template #header>
        <h2>按钮权限控制演示</h2>
      </template>
      
      <div class="demo-content">
        <el-alert
          title="权限控制说明"
          type="info"
          :closable="false"
          show-icon
        >
          <p>本系统实现了基于权限编码的按钮级权限控制，支持以下功能：</p>
          <ul>
            <li>权限指令：使用 <code>v-permission="'权限编码'"</code> 控制按钮显示</li>
            <li>权限函数：使用 <code>hasPermission('权限编码')</code> 进行权限判断</li>
            <li>角色指令：使用 <code>v-role="'角色名'"</code> 控制基于角色的显示</li>
            <li>支持多权限组合判断</li>
          </ul>
        </el-alert>

        <div class="section">
          <h3>当前用户权限</h3>
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
              暂无权限编码
            </el-tag>
          </div>
        </div>

        <div class="section">
          <h3>权限指令演示</h3>
          <div class="demo-group">
            <h4>单个权限判断</h4>
            <div class="button-group">
              <el-button v-permission="'user:create'" type="primary">
                用户创建 (user:create)
              </el-button>
              <el-button v-permission="'user:update'" type="success">
                用户编辑 (user:update)
              </el-button>
              <el-button v-permission="'user:delete'" type="danger">
                用户删除 (user:delete)
              </el-button>
              <el-button v-permission="'admin:all'" type="warning">
                超管权限 (admin:all)
              </el-button>
            </div>
          </div>

          <div class="demo-group">
            <h4>多权限组合判断</h4>
            <div class="button-group">
              <el-button v-permission="['user:create', 'user:update']" type="primary">
                用户创建或编辑
              </el-button>
              <el-button v-permission="['role:create', 'role:update']" type="success">
                角色创建或编辑
              </el-button>
            </div>
          </div>
        </div>

        <div class="section">
          <h3>编程式权限判断</h3>
          <div class="code-demo">
            <pre><code>// 使用权限函数进行判断
if (hasPermission('user:create')) {
  // 有权限时的逻辑
  console.log('可以创建用户')
} else {
  // 无权限时的逻辑
  console.log('无创建用户权限')
}

// 多权限判断 - 任一权限
hasPermission(['user:create', 'user:update'], 'some')

// 多权限判断 - 全部权限
hasPermission(['user:create', 'user:update'], 'every')</code></pre>
          </div>

          <div class="button-group">
            <el-button @click="testPermission('user:create')" type="primary">
              测试用户创建权限
            </el-button>
            <el-button @click="testPermission('user:delete')" type="danger">
              测试用户删除权限
            </el-button>
            <el-button @click="testPermission(['user:create', 'user:update'])" type="success">
              测试多权限组合
            </el-button>
          </div>
        </div>

        <div class="section">
          <h3>菜单权限配置示例</h3>
          <el-table :data="menuExamples" border>
            <el-table-column prop="title" label="菜单名称" width="150" />
            <el-table-column prop="menuType" label="类型" width="80">
              <template #default="{ row }">
                <el-tag :type="row.menuType === 'button' ? 'info' : 'success'">
                  {{ row.menuType === 'button' ? '按钮' : '菜单' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="path" label="路由路径" width="150" />
            <el-table-column prop="permissionCode" label="权限编码" width="150">
              <template #default="{ row }">
                <el-tag v-if="row.permissionCode" type="primary">
                  {{ row.permissionCode }}
                </el-tag>
                <span v-else class="text-gray-400">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" />
          </el-table>
        </div>

        <div class="section">
          <h3>最佳实践</h3>
          <el-alert
            title="权限编码命名建议"
            type="success"
            :closable="false"
            show-icon
          >
            <div class="best-practices">
              <h4>命名规范：</h4>
              <ul>
                <li><strong>模块:操作</strong> 格式，如 <code>user:create</code></li>
                <li><strong>资源:动作</strong> 格式，如 <code>order:export</code></li>
                <li>使用小写字母和冒号分隔</li>
                <li>保持简洁和语义化</li>
              </ul>

              <h4>常用权限编码示例：</h4>
              <div class="example-codes">
                <el-tag class="code-tag">user:create</el-tag>
                <el-tag class="code-tag">user:update</el-tag>
                <el-tag class="code-tag">user:delete</el-tag>
                <el-tag class="code-tag">user:view</el-tag>
                <el-tag class="code-tag">role:assign</el-tag>
                <el-tag class="code-tag">order:export</el-tag>
                <el-tag class="code-tag">system:backup</el-tag>
              </div>
            </div>
          </el-alert>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { hasPermission } from '@/utils/permission'

const userStore = useUserStore()

// 获取用户权限列表
const userPermissions = computed(() => {
  return userStore.userPermissions || []
})

// 菜单配置示例
const menuExamples = [
  {
    title: '用户管理',
    menuType: 'menu',
    path: '/system/user',
    permissionCode: '',
    description: '菜单类型，不需要权限编码'
  },
  {
    title: '新增用户',
    menuType: 'button',
    path: '',
    permissionCode: 'user:create',
    description: '按钮类型，需要权限编码控制显示'
  },
  {
    title: '编辑用户',
    menuType: 'button',
    path: '',
    permissionCode: 'user:update',
    description: '编辑用户按钮权限'
  },
  {
    title: '删除用户',
    menuType: 'button',
    path: '',
    permissionCode: 'user:delete',
    description: '删除用户按钮权限'
  },
  {
    title: '导出用户',
    menuType: 'button',
    path: '',
    permissionCode: 'user:export',
    description: '导出功能按钮权限'
  }
]

// 测试权限函数
const testPermission = (permission) => {
  const hasAuth = hasPermission(permission)
  const permText = Array.isArray(permission) ? permission.join(', ') : permission
  
  if (hasAuth) {
    ElMessage.success(`拥有权限: ${permText}`)
  } else {
    ElMessage.warning(`无权限: ${permText}`)
  }
}
</script>

<style lang="scss" scoped>
.permission-demo {
  .demo-content {
    .section {
      margin-bottom: 32px;
      
      h3 {
        margin-bottom: 16px;
        color: #303133;
        border-bottom: 2px solid #409eff;
        padding-bottom: 8px;
      }
      
      h4 {
        margin: 16px 0 12px 0;
        color: #606266;
      }
    }

    .permission-list {
      .permission-tag {
        margin-right: 8px;
        margin-bottom: 8px;
      }
    }

    .demo-group {
      margin-bottom: 24px;
      
      .button-group {
        .el-button {
          margin-right: 12px;
          margin-bottom: 8px;
        }
      }
    }

    .code-demo {
      background: #f6f8fa;
      border: 1px solid #e1e4e8;
      border-radius: 6px;
      padding: 16px;
      margin-bottom: 16px;
      
      pre {
        margin: 0;
        color: #24292e;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 13px;
        line-height: 1.5;
        
        code {
          background: transparent;
          color: inherit;
          padding: 0;
          font-size: inherit;
        }
      }
    }

    .best-practices {
      ul {
        margin: 8px 0;
        padding-left: 20px;
        
        li {
          margin-bottom: 4px;
          
          code {
            background: #f1f3f4;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: monospace;
            color: #e83e8c;
          }
        }
      }
      
      .example-codes {
        margin-top: 12px;
        
        .code-tag {
          margin-right: 8px;
          margin-bottom: 8px;
          font-family: monospace;
        }
      }
    }
  }
}

.text-gray-400 {
  color: #9ca3af;
}
</style> 