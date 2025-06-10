# go-gin-element-admin

基于 Go + Gin + MySQL8 + GORM + Redis + Vue3 + Element-Plus 构建的极简后端管理系统脚手架

## 技术栈

### 后端
- **Go 1.21** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL 8** - 数据库
- **Redis** - 缓存
- **Casbin** - 权限管理
- **JWT** - 身份认证
- **Viper** - 配置管理

### 前端
- **Vue 3** - 前端框架
- **Element Plus** - UI 组件库
- **Vite** - 构建工具
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Axios** - HTTP 客户端

## 功能特性

- 🔐 **用户认证** - JWT 登录认证
- 👤 **个人中心** - 个人信息管理、密码修改、头像上传
- 🏠 **仪表板** - 系统概览、数据统计、快捷操作
- 👥 **用户管理** - 用户增删改查、状态管理
- 🛡️ **角色管理** - 角色权限管理、基于 RBAC 的权限控制
- 📋 **菜单管理** - 动态菜单配置、权限分配
- 📊 **操作日志** - 详细的操作审计、实时监控、数据统计
- 🎨 **现代化 UI** - 响应式设计、美观的界面

## 项目结构

```
go-gin-element-admin/
├── server/                 # 后端代码
│   ├── main.go            # 程序入口
│   ├── config/            # 配置文件
│   ├── core/              # 核心功能
│   ├── global/            # 全局变量
│   ├── initialize/        # 初始化
│   └── model/             # 数据模型
├── web/                   # 前端代码
│   ├── src/
│   │   ├── api/          # API 接口
│   │   ├── components/   # 组件
│   │   ├── layout/       # 布局
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   ├── utils/        # 工具函数
│   │   └── views/        # 页面
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- Redis 6.0+

### 后端启动

1. 进入后端目录
```bash
cd server
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库
编辑 `server/config/config.yaml` 文件，修改数据库连接信息：
```yaml
mysql:
  path: '127.0.0.1:3306'
  db-name: 'go_gin_element_admin'
  username: 'root'
  password: 'your_password'
```

4. 启动服务
```bash
go run main.go
```

后端服务将在 `http://localhost:8888` 启动

### 前端启动

1. 进入前端目录
```bash
cd web
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

前端服务将在 `http://localhost:8080` 启动

### 默认账号

- 用户名：`admin`
- 密码：`123456`

> 注意：请修改 `server/config/config.yaml` 中的数据库密码为你的实际密码

## 开发说明

### 后端开发

- 配置文件位于 `server/config/config.yaml`
- 数据模型定义在 `server/model/` 目录
- API 路由配置在 `server/router/` 目录
- 业务逻辑在 `server/service/` 目录

### 前端开发

- 页面组件在 `web/src/views/` 目录
- API 接口定义在 `web/src/api/` 目录
- 路由配置在 `web/src/router/index.js`
- 状态管理在 `web/src/stores/` 目录

## 部署

### 后端部署

1. 编译
```bash
cd server
go build -o app main.go
```

2. 运行
```bash
./app
```

### 前端部署

1. 构建
```bash
cd web
npm run build
```

2. 部署 `dist` 目录到 Web 服务器

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
