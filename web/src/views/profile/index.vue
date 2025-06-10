<template>
  <div class="profile" v-loading="pageLoading" element-loading-text="加载用户信息中..." style="padding: 20px;">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <div class="avatar-section">
            <el-avatar :size="120" :src="userInfo.headerImg || defaultAvatar" />
            <div class="avatar-info">
              <h3>{{ userInfo.nickName || userInfo.username }}</h3>
              <p>{{ userInfo.email }}</p>
              <el-tag type="success">{{ userInfo.authority?.authorityName || '普通用户' }}</el-tag>
            </div>
          </div>
        </el-col>
        
        <el-col :span="16">
          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="100px"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="profileForm.username" disabled />
            </el-form-item>
            
            <el-form-item label="昵称" prop="nickName">
              <el-input v-model="profileForm.nickName" />
            </el-form-item>
            
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            
            <el-form-item label="头像" prop="headerImg">
              <div class="avatar-upload">
                <el-upload
                  class="avatar-uploader"
                  :show-file-list="false"
                  :before-upload="beforeAvatarUpload"
                  :http-request="customUpload"
                  accept="image/jpeg,image/jpg,image/png,image/gif"
                  :disabled="avatarUploading"
                >
                  <div class="avatar-container" v-loading="avatarUploading" element-loading-text="上传中...">
                    <img v-if="profileForm.headerImg" :src="getAvatarUrl(profileForm.headerImg)" class="avatar" />
                    <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
                  </div>
                </el-upload>
                <div class="avatar-tips">
                  <p>支持 JPG、PNG、GIF 格式</p>
                  <p>文件大小不超过 10MB</p>
                  <p>点击头像上传新图片</p>
                </div>
              </div>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="updateProfile" :loading="loading">
                更新信息
              </el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-card>

    <el-card style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>修改密码</span>
        </div>
      </template>
      
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
        style="max-width: 500px;"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="changeUserPassword" :loading="passwordLoading">
            修改密码
          </el-button>
          <el-button @click="resetPasswordForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getUserInfo, updateUserInfo, changePassword } from '@/api/system/user'
import { uploadAvatar } from '@/api/upload'
import { Plus } from '@element-plus/icons-vue'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

const profileFormRef = ref()
const passwordFormRef = ref()
const loading = ref(false)
const passwordLoading = ref(false)
const pageLoading = ref(false)
const avatarUploading = ref(false)

const profileForm = reactive({
  username: '',
  nickName: '',
  email: '',
  phone: '',
  headerImg: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const profileRules = {
  nickName: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    pageLoading.value = true
    const response = await getUserInfo()
    if (response.code === 0) {
      const userData = response.data
      // 映射后端数据到前端期望的格式
      const mappedUserData = {
        id: userData.ID,
        username: userData.username,
        nickName: userData.nickName,
        email: userData.email || '',
        phone: userData.phone || '',
        headerImg: userData.headerImg || '',
        authority: {
          authorityId: userData.authorityId,
          authorityName: userData.authorityId === 888 ? '管理员' : '普通用户'
        }
      }

      // 更新store中的用户信息
      Object.assign(userStore.userInfo, mappedUserData)
      // 初始化表单
      initForm()
    } else {
      ElMessage.error(response.msg || '获取用户信息失败')
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    pageLoading.value = false
  }
}

const initForm = () => {
  Object.assign(profileForm, {
    username: userInfo.value.username || '',
    nickName: userInfo.value.nickName || '',
    email: userInfo.value.email || '',
    phone: userInfo.value.phone || '',
    headerImg: userInfo.value.headerImg || ''
  })
}

const updateProfile = () => {
  profileFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        console.log('发送更新请求，数据:', profileForm)
        const response = await updateUserInfo(profileForm)
        console.log('更新响应:', response)
        ElMessage.success('更新成功')
        // 更新store中的用户信息
        Object.assign(userStore.userInfo, profileForm)
        // 重新获取用户信息以确保数据同步
        await fetchUserInfo()
      } catch (error) {
        console.error('更新失败:', error)
        ElMessage.error(error.message || '更新失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const changeUserPassword = () => {
  passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      try {
        await changePassword(passwordForm)
        ElMessage.success('密码修改成功')
        resetPasswordForm()
      } catch (error) {
        ElMessage.error(error.message || '密码修改失败')
      } finally {
        passwordLoading.value = false
      }
    }
  })
}

const resetForm = () => {
  initForm()
}

const resetPasswordForm = () => {
  passwordFormRef.value?.resetFields()
  Object.assign(passwordForm, {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
}

const beforeAvatarUpload = (file) => {
  const isValidType = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'].includes(file.type)
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isValidType) {
    ElMessage.error('上传头像图片只能是 JPG、PNG、GIF 格式!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('上传头像图片大小不能超过 10MB!')
    return false
  }
  return true
}

const customUpload = async (options) => {
  try {
    avatarUploading.value = true
    const response = await uploadAvatar(options.file)
    
    if (response.code === 0) {
      // 构建完整的URL
      const baseUrl = window.location.origin
      const avatarUrl = response.data.url.startsWith('http') 
        ? response.data.url 
        : `${baseUrl}${response.data.url}`
      
      profileForm.headerImg = avatarUrl
      ElMessage.success('头像上传成功')
      
      // 调用成功回调
      options.onSuccess(response)
    } else {
      ElMessage.error(response.msg || '上传失败')
      options.onError(new Error(response.msg || '上传失败'))
    }
  } catch (error) {
    console.error('上传头像失败:', error)
    ElMessage.error('上传失败: ' + (error.message || '网络错误'))
    options.onError(error)
  } finally {
    avatarUploading.value = false
  }
}

const getAvatarUrl = (url) => {
  if (!url) return defaultAvatar
  return url.startsWith('http') ? url : defaultAvatar
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style lang="scss" scoped>
.profile {
  .avatar-section {
    text-align: center;

    .avatar-info {
      margin-top: 20px;

      h3 {
        margin: 10px 0;
        color: #303133;
      }

      p {
        color: #909399;
        margin: 5px 0;
      }
    }
  }

  .card-header {
    font-weight: 500;
    color: #303133;
  }

  .avatar-upload {
    display: flex;
    align-items: center;
    gap: 20px;

    .avatar-uploader {
      .avatar-container {
        width: 120px;
        height: 120px;
        border: 2px dashed #dcdfe6;
        border-radius: 50%;
        cursor: pointer;
        position: relative;
        overflow: hidden;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;

        &:hover {
          border-color: #409eff;
        }

        .avatar {
          width: 100%;
          height: 100%;
          object-fit: cover;
          border-radius: 50%;
        }

        .avatar-uploader-icon {
          font-size: 28px;
          color: #8c939d;
        }
      }
    }

    .avatar-tips {
      color: #909399;
      font-size: 12px;
      line-height: 1.5;

      p {
        margin: 0 0 4px 0;
      }
    }
  }
}
</style> 