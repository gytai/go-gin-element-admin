<template>
  <div class="login-container">
    <div class="login-form">
      <div class="title-container">
        <h3 class="title">Go-Gin-Element-Admin</h3>
        <p class="subtitle">后台管理系统脚手架</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form-content"
        auto-complete="on"
        label-position="left"
      >
        <el-form-item prop="username">
          <span class="svg-container">
            <el-icon><User /></el-icon>
          </span>
          <el-input
            ref="username"
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            tabindex="1"
            auto-complete="on"
          />
        </el-form-item>

        <el-form-item prop="password">
          <span class="svg-container">
            <el-icon><Lock /></el-icon>
          </span>
          <el-input
            :key="passwordType"
            ref="password"
            v-model="loginForm.password"
            :type="passwordType"
            placeholder="密码"
            name="password"
            tabindex="2"
            auto-complete="on"
            @keyup.enter="handleLogin"
          />
          <span class="show-pwd" @click="showPwd">
            <el-icon>
              <component :is="passwordType === 'password' ? 'View' : 'Hide'" />
            </el-icon>
          </span>
        </el-form-item>

        <el-button
          :loading="loading"
          type="primary"
          style="width: 100%; margin-bottom: 30px"
          @click.prevent="handleLogin"
        >
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginForm = reactive({
  username: 'admin',
  password: '123456'
})

const loginRules = {
  username: [{ required: true, trigger: 'blur', message: '请输入用户名' }],
  password: [{ required: true, trigger: 'blur', message: '请输入密码' }]
}

const passwordType = ref('password')
const loading = ref(false)
const loginFormRef = ref()

const showPwd = () => {
  if (passwordType.value === 'password') {
    passwordType.value = ''
  } else {
    passwordType.value = 'password'
  }
  nextTick(() => {
    loginFormRef.value.password.focus()
  })
}

const handleLogin = () => {
  loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(loginForm)
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        ElMessage.error(error.message || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background: 
    linear-gradient(135deg, rgba(102, 126, 234, 0.25) 0%, rgba(118, 75, 162, 0.35) 100%),
    url('https://images.unsplash.com/photo-1519389950473-47ba0277781c?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80') center/cover no-repeat;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  
  /* 添加动态粒子效果 */
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.15) 0%, transparent 50%),
      radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.08) 0%, transparent 50%),
      radial-gradient(circle at 40% 40%, rgba(120, 119, 198, 0.1) 0%, transparent 50%);
    animation: floating 6s ease-in-out infinite;
    pointer-events: none;
  }
  
  @keyframes floating {
    0%, 100% { opacity: 0.7; transform: translateY(0px); }
    50% { opacity: 1; transform: translateY(-10px); }
  }

  .login-form {
    width: 520px;
    max-width: 100%;
    padding: 60px 40px 40px;
    margin: 0 auto;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.25);
    backdrop-filter: blur(6px);
    border-radius: 24px;
    box-shadow: 
      0 8px 32px 0 rgba(31, 38, 135, 0.4),
      0 0 0 1px rgba(255, 255, 255, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.4);
    position: relative;
    z-index: 10;
    
    /* 添加微妙的动画效果 */
    animation: formSlideIn 0.8s ease-out;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(145deg, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0.08) 100%);
      border-radius: 24px;
      pointer-events: none;
    }
  }
  
  @keyframes formSlideIn {
    0% {
      opacity: 0;
      transform: translateY(30px) scale(0.95);
    }
    100% {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  .tips {
    margin-bottom: 20px;
    text-align: center;
    z-index: 20;
    position: relative;

    .tip-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 16px;
      margin-bottom: 8px;
      background: rgba(255, 255, 255, 0.15);
      border-radius: 8px;
      border: 1px solid rgba(255, 255, 255, 0.2);
      transition: all 0.3s ease;

      &:hover {
        background: rgba(255, 255, 255, 0.22);
        transform: translateY(-1px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
      }

      &:last-child {
        margin-bottom: 0;
      }

      .tip-label {
        font-size: 14px;
        color: rgba(255, 255, 255, 0.8);
        font-weight: 400;
      }

      .tip-value {
        font-size: 14px;
        color: #fff;
        font-weight: 600;
        font-family: 'Courier New', monospace;
        letter-spacing: 1px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: #889aa4;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;
    text-align: center;
    margin-bottom: 40px;
    z-index: 20;

    .title {
      font-size: 28px;
      color: #fff;
      margin: 0px auto 15px auto;
      font-weight: 700;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
      background: linear-gradient(135deg, #fff 0%, #f0f0f0 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      letter-spacing: 1px;
    }

    .subtitle {
      font-size: 16px;
      color: rgba(255, 255, 255, 0.9);
      margin: 0;
      font-weight: 400;
      text-shadow: 0 1px 5px rgba(0, 0, 0, 0.2);
    }
  }

  .show-pwd {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
    display: flex;
    align-items: center;
    height: 20px;
    transition: color 0.3s ease;
    
    &:hover {
      color: #667eea;
    }
  }

  .login-form-content {
    z-index: 20;
    position: relative;
    
    .el-form-item {
      border: 1px solid rgba(255, 255, 255, 0.3);
      background: rgba(255, 255, 255, 0.15);
      border-radius: 12px;
      color: #454545;
      margin-bottom: 22px;
      transition: all 0.3s ease;
      
      &:hover {
        border-color: rgba(255, 255, 255, 0.4);
        background: rgba(255, 255, 255, 0.22);
        transform: translateY(-2px);
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
      }

      .el-input {
        display: inline-block;
        height: 52px;
        width: 85%;

        :deep(.el-input__wrapper) {
          padding: 0;
          background: transparent;
          box-shadow: none;

          .el-input__inner {
            background: transparent;
            border: 0px;
            border-radius: 0px;
            padding: 15px 5px 15px 15px;
            color: #fff;
            height: 52px;
            caret-color: #fff;
            font-size: 15px;

            &::placeholder {
              color: rgba(255, 255, 255, 0.6);
            }

            &:-webkit-autofill {
              box-shadow: 0 0 0px 1000px transparent inset !important;
              -webkit-text-fill-color: #fff !important;
            }
          }
        }
      }
    }

    .el-button {
      width: 100%;
      height: 52px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
      border-radius: 12px;
      font-size: 16px;
      font-weight: 600;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(102, 126, 234, 0.6);
        background: linear-gradient(135deg, #7c8eeb 0%, #8b6bb1 100%);
      }
      
      &:active {
        transform: translateY(0);
        box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
      }
    }
  }
}
</style> 