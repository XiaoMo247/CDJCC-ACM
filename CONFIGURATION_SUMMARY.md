# 项目配置管理总结

本文档汇总了前后端项目的配置管理方案。

## 概述

为了规范项目的配置管理，提高代码安全性和可维护性，我们统一了前后端项目的配置方式，将所有硬编码的配置迁移到 `.env` 文件进行集中管理。

## 前端配置

**位置**：`fronted-acm/`

### 核心配置文件

| 文件 | 说明 |
|------|------|
| `.env` | 开发环境配置（不提交） |
| `.env.example` | 配置模板（提交到版本控制） |
| `CONFIG_GUIDE.md` | 详细的配置指南 |
| `vite.config.js` | Vite 构建配置 |
| `src/utils/request.js` | API 请求工具 |
| `src/utils/tokenManager.js` | Token 管理工具 |

### 前端配置项

```env
# API 配置
VITE_API_BASE_URL=http://localhost:8081/api

# 上传服务器配置
VITE_UPLOAD_TARGET=http://localhost:8081
VITE_UPLOAD_PATH=/uploads

# Token 配置
VITE_TOKEN_REFRESH_INTERVAL=3000000
VITE_TOKEN_KEY=auth_token
VITE_USER_INFO_KEY=user_info

# 外部链接
VITE_QQ_GROUP_URL=https://jq.qq.com/?_wv=1027&k=KoVRvLvL
VITE_BILIBILI_URL=https://space.bilibili.com/...
VITE_OJ_URL=https://hydro.ac/d/cdjcc_acm_2333/

# 默认头像
VITE_DEFAULT_AVATAR=https://cube.elemecdn.com/...

# 应用信息
VITE_APP_NAME=CDJCC ACM 集训队官网
```

### 修改的前端文件

1. **vite.config.js**
   - 上传代理地址改为使用 `process.env.VITE_UPLOAD_TARGET`

2. **src/utils/tokenManager.js**
   - Token 存储 key 改为使用环境变量
   - Token 刷新间隔改为使用环境变量

3. **src/pages/HonorWall.vue**
   - 默认头像 URL 改为使用环境变量

4. **src/pages/KnowledgePage.vue**
   - Bilibili 和 OJ 链接改为使用环境变量

5. **src/pages/ForgotPassword.vue**
   - QQ 群链接改为使用环境变量

## 后端配置

**位置**：`acm-site/`

### 核心配置文件

| 文件 | 说明 |
|------|------|
| `.env` | 开发环境配置（不提交） |
| `.env.example` | 配置模板（提交到版本控制） |
| `BACKEND_CONFIG_GUIDE.md` | 详细的配置指南 |
| `config/config.go` | 配置加载代码 |
| `config/config.yml` | YAML 配置文件（可选） |
| `main.go` | 应用入口 |

### 后端配置项

```env
# 数据库配置
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=acm_site
DB_CHARSET=utf8mb4
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100

# JWT 配置
JWT_SECRET_KEY=your_secret_key_here
JWT_EXPIRATION_HOURS=24

# 服务器配置
SERVER_PORT=8081
GIN_MODE=debug

# CORS 配置
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

# 文件上传配置
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=52428800

# 日志配置
LOG_LEVEL=info
LOG_FILE=./logs/app.log

# 安全配置
SUPER_ADMIN_IDS=1
ENABLE_HTTPS=false
SSL_CERT_FILE=
SSL_KEY_FILE=

# 邮件配置（可选）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
MAIL_FROM=noreply@your-domain.com

# 其他配置
DEFAULT_PAGE_SIZE=20
CACHE_EXPIRATION_SECONDS=3600
```

### 修改的后端文件

1. **config/config.go**
   - 扩展配置结构体，支持 JWT、CORS、安全等多个模块
   - 改进环境变量加载和类型转换
   - 添加配置访问辅助函数

2. **main.go**
   - 增强环境变量检查和验证
   - 添加生产环境安全提示

3. **.env.example**
   - 补充所有新增配置项
   - 添加详细的配置说明

## 配置优先级

两个项目的配置优先级一致（高到低）：

### 前端
1. 环境变量（通过 `import.meta.env`）
2. `.env` 文件
3. 代码中的硬编码默认值

### 后端
1. 环境变量（通过 `os.Getenv()`）
2. `.env` 文件
3. `config/config.yml` 文件
4. 代码中的硬编码默认值

## 使用方式

### 开发环境

**前端**：
```bash
cd fronted-acm
cp .env.example .env
# 编辑 .env 文件
npm install
npm run dev
```

**后端**：
```bash
cd acm-site
cp .env.example .env
# 编辑 .env 文件
go run main.go
```

### 生产环境部署

**前端构建**：
```bash
cd fronted-acm
# 设置生产环境变量
export VITE_API_BASE_URL=https://your-domain.com/api
export VITE_UPLOAD_TARGET=https://your-domain.com
# 构建
npm run build
```

**后端部署**：
```bash
cd acm-site
# 加载 .env 文件中的配置
export $(cat .env | xargs)
# 启动应用
go run main.go
```

## 安全建议

### 版本控制

- ✅ 提交：`.env.example`、配置指南文档、配置加载代码
- ❌ 不提交：`.env` 文件

**添加到 .gitignore**：
```
.env
.env.local
.env.*.local
```

### 敏感信息管理

前端：
- API 密钥（如果有）存放在 `.env` 中
- 不在组件中硬编码任何密钥

后端：
- JWT_SECRET_KEY 必须在生产环境修改
- 数据库密码必须使用强密码
- 不在日志中输出敏感信息

### 生产环境配置检查清单

前端：
- ✅ `VITE_API_BASE_URL` 指向生产服务器
- ✅ `VITE_UPLOAD_TARGET` 指向生产服务器
- ✅ 外部链接指向正确的生产地址

后端：
- ✅ `JWT_SECRET_KEY` 已修改为安全密钥（使用 `openssl rand -base64 32` 生成）
- ✅ 数据库凭据已更新
- ✅ `ALLOWED_ORIGINS` 限制为生产域名
- ✅ `GIN_MODE` 设置为 `release`
- ✅ 日志文件路径已配置
- ✅ 如需要，`ENABLE_HTTPS` 已启用且证书路径正确

## 配置管理工作流

### 新增配置项

1. 在 `.env.example` 中添加新配置项（带说明注释）
2. 在 `.env` 中添加新配置项（开发环境值）
3. 在相应的代码中读取该配置
4. 更新对应的配置指南文档
5. 提交更改

### 修改配置值

1. 编辑 `.env` 文件
2. 如果涉及多个开发者，更新 `.env.example`
3. 重启应用使配置生效

### 环境差异处理

对于不同环境的配置差异：

1. 在 CI/CD 系统中设置环境变量（推荐）
2. 或为每个环境创建独立的 `.env` 文件（如 `.env.production`）

## 文档导航

### 前端配置
详见：`fronted-acm/CONFIG_GUIDE.md`

关键点：
- API 基础 URL 管理
- Token 刷新间隔和存储 key
- 外部链接配置
- 环境变量访问方式

### 后端配置
详见：`acm-site/BACKEND_CONFIG_GUIDE.md`

关键点：
- 数据库连接配置
- JWT 密钥和过期时间
- CORS 跨域配置
- 文件上传限制
- 日志和安全配置

## 常见问题

### Q: 如何只修改某个特定环境的配置？
A: 使用环境变量覆盖。例如在 Docker 中：
```yaml
environment:
  - DB_HOST=production-mysql.example.com
  - DB_PASSWORD=secure_password
```

### Q: 修改配置后需要重新部署吗？
A:
- 前端：需要重新构建（`npm run build`）
- 后端：只需重启应用（`go run main.go`）

### Q: 如何在不同的开发者之间同步配置？
A: 使用 `.env.example` 作为模板，每个开发者拷贝后根据本地环境修改。

### Q: 密钥泄露了怎么办？
A:
1. 立即生成新密钥
2. 更新 `.env` 文件
3. 重启应用
4. 如果是 JWT_SECRET_KEY，建议要求用户重新登录

## 总结

通过统一的配置管理方案，我们实现了：

✅ **代码与配置分离**：敏感信息不在代码中硬编码
✅ **环境灵活性**：同一份代码可以部署到多个环境
✅ **安全性提升**：配置文件易于加入 .gitignore，防止泄露
✅ **可维护性**：集中管理，修改简单，文档完善
✅ **易于上手**：新成员只需复制 .env.example 就能快速启动开发

---

**最后更新**：2026-01-13
**前端配置版本**：1.0
**后端配置版本**：1.0
