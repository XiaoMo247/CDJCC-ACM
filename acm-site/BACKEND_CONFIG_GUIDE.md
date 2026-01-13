# 后端配置指南

本文档说明了后端项目中所有配置项的使用方法。所有配置项都集中在 `.env` 文件中管理。

## 快速开始

1. 复制 `.env.example` 为 `.env`：
   ```bash
   cp .env.example .env
   ```

2. 根据实际环境修改 `.env` 文件中的配置

3. 重启应用使配置生效（如果应用正在运行）

## 配置优先级

配置采用分层加载策略，优先级从高到低：

1. **环境变量** （优先级最高）
2. **.env 文件**
3. **config/config.yml**
4. **代码中的默认值** （优先级最低）

## 配置详解

### 1. 数据库配置

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=acm_site
DB_CHARSET=utf8mb4
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
```

- **DB_HOST**：MySQL 服务器地址
- **DB_PORT**：MySQL 服务器端口（默认 3306）
- **DB_USER**：数据库用户名
- **DB_PASSWORD**：数据库密码
- **DB_NAME**：数据库名称
- **DB_CHARSET**：字符集（推荐 utf8mb4）
- **DB_MAX_IDLE_CONNS**：最大闲置连接数
- **DB_MAX_OPEN_CONNS**：最大打开连接数

**读取位置**：`config/config.go` 中的 `InitConfig()` 函数

**使用位置**：`database/init.go` 中初始化数据库连接

### 2. JWT 配置

```env
JWT_SECRET_KEY=your_super_secret_jwt_key_here_change_this_in_production
JWT_EXPIRATION_HOURS=24
```

- **JWT_SECRET_KEY**：JWT 签名密钥（生产环境必须修改）
  - 生成安全密钥：`openssl rand -base64 32`
- **JWT_EXPIRATION_HOURS**：Token 过期时间（小时），默认 24 小时

**读取位置**：`config/config.go`

**使用位置**：认证中间件和 Token 生成服务

### 3. 服务器配置

```env
SERVER_PORT=8081
GIN_MODE=debug
```

- **SERVER_PORT**：应用监听端口
- **GIN_MODE**：运行模式
  - `debug`：调试模式，输出详细日志
  - `release`：发布模式，优化性能

**使用位置**：`main.go` 中的 `r.Run(":8081")`

### 4. CORS 配置

```env
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000
```

- **说明**：允许跨域请求的源地址（多个用逗号分隔）
- **开发环境**：`http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000`
- **生产环境示例**：`https://your-domain.com,https://www.your-domain.com`

**读取位置**：`config/config.go`

**使用位置**：`middleware/cors.go`

### 5. 文件上传配置

```env
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=52428800
```

- **UPLOAD_DIR**：上传文件保存目录
- **MAX_UPLOAD_SIZE**：最大上传文件大小（字节）
  - 默认 50MB（52428800 字节）
  - 1GB = 1073741824 字节

**读取位置**：使用时从环境变量直接读取

**使用位置**：`api/folder.go` 等文件上传相关的处理函数

### 6. 日志配置

```env
LOG_LEVEL=info
LOG_FILE=./logs/app.log
```

- **LOG_LEVEL**：日志级别
  - `debug`：调试信息（最详细）
  - `info`：一般信息
  - `warn`：警告信息
  - `error`：错误信息（最简略）
- **LOG_FILE**：日志文件路径
  - 空值或不设置：输出到控制台
  - 指定路径：输出到文件

**读取位置**：日志初始化模块

### 7. 安全配置

```env
SUPER_ADMIN_IDS=1
ENABLE_HTTPS=false
SSL_CERT_FILE=
SSL_KEY_FILE=
```

- **SUPER_ADMIN_IDS**：超级管理员用户 ID（逗号分隔）
- **ENABLE_HTTPS**：是否启用 HTTPS（生产环境建议启用）
- **SSL_CERT_FILE**：SSL 证书文件路径
- **SSL_KEY_FILE**：SSL 密钥文件路径

**读取位置**：`config/config.go`

**使用位置**：权限校验中间件

### 8. 邮件配置（可选）

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
MAIL_FROM=noreply@your-domain.com
```

- **SMTP_HOST**：SMTP 服务器地址
- **SMTP_PORT**：SMTP 端口
- **SMTP_USER**：SMTP 用户名/邮箱
- **SMTP_PASSWORD**：SMTP 密码/应用密码
- **MAIL_FROM**：发件人邮箱地址

**说明**：用于系统发送通知邮件，如密码重置、消息通知等

### 9. 其他配置

```env
DEFAULT_PAGE_SIZE=20
CACHE_EXPIRATION_SECONDS=3600
```

- **DEFAULT_PAGE_SIZE**：分页默认页码大小
- **CACHE_EXPIRATION_SECONDS**：缓存过期时间（秒）

## 环境变量访问方式

在 Go 代码中访问环境变量：

```go
// 直接读取环境变量
jwtKey := os.Getenv("JWT_SECRET_KEY")
port := os.Getenv("DB_PORT")

// 通过配置包读取（推荐）
jwtConfig := config.GetJWTConfig()
mysqlConfig := config.GetMySQLConfig()
```

## 部署指南

### 开发环境

1. 复制配置模板：
   ```bash
   cp .env.example .env
   ```

2. 编辑 `.env` 文件，设置本地开发环境：
   ```env
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=123456
   DB_NAME=acm_site
   GIN_MODE=debug
   JWT_SECRET_KEY=dev_key_for_testing_only
   ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
   ```

3. 启动应用：
   ```bash
   go run main.go
   ```

### 生产环境

1. **准备 .env 文件**：
   ```bash
   cp .env.example .env
   ```

2. **生成安全的 JWT 密钥**：
   ```bash
   openssl rand -base64 32
   ```

3. **更新 .env 文件中的敏感配置**：
   ```env
   # 数据库配置
   DB_HOST=your-db-server.com
   DB_PORT=3306
   DB_USER=production_user
   DB_PASSWORD=<strong_password>
   DB_NAME=acm_site_prod

   # JWT 配置（使用生成的密钥）
   JWT_SECRET_KEY=<generated_secure_key>
   JWT_EXPIRATION_HOURS=24

   # 服务器配置
   GIN_MODE=release
   SERVER_PORT=8081

   # CORS 配置
   ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com

   # 日志配置
   LOG_LEVEL=info
   LOG_FILE=./logs/app.log

   # 文件上传配置
   UPLOAD_DIR=/var/uploads/acm-site
   MAX_UPLOAD_SIZE=52428800

   # 安全配置
   ENABLE_HTTPS=true
   SSL_CERT_FILE=/etc/ssl/certs/your-domain.crt
   SSL_KEY_FILE=/etc/ssl/private/your-domain.key
   ```

4. **确保目录权限**：
   ```bash
   mkdir -p logs uploads
   chmod 755 logs uploads
   ```

5. **启动应用**：
   ```bash
   export $(cat .env | xargs)  # 加载 .env 文件
   go run main.go
   ```

### Docker 部署

使用 `docker-compose` 时，在 `docker-compose.yml` 中设置环境变量：

```yaml
services:
  acm-site:
    build: .
    ports:
      - "8081:8081"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=your_secure_password
      - DB_NAME=acm_site
      - JWT_SECRET_KEY=your_secure_jwt_key
      - GIN_MODE=release
      - ALLOWED_ORIGINS=https://your-domain.com
    depends_on:
      - mysql
```

### CI/CD 部署

在 GitLab CI、GitHub Actions 等 CI/CD 系统中，设置环境变量：

```bash
# GitLab CI 或 GitHub Actions 的 Secret 变量
DB_HOST=your-db-server.com
DB_USER=production_user
DB_PASSWORD=<secret>
JWT_SECRET_KEY=<secret>
```

## 常见问题

### Q: 如何在不同环境使用不同的配置？
A: 可以创建多个 `.env` 文件：
- `.env.development` - 开发环境
- `.env.production` - 生产环境
- `.env.staging` - 测试环境

启动时指定文件：`export $(cat .env.production | xargs)`

### Q: 如何在生产环境验证配置是否正确？
A: 启动应用时，会输出警告信息：
- JWT_SECRET_KEY 仍为默认值：⚠️ 安全警告
- ALLOWED_ORIGINS 未设置：⚠️ 安全警告
- 数据库使用默认凭据：⚠️ 安全警告

### Q: 如何修改最大上传文件大小？
A: 编辑 `.env` 文件：
```env
# 修改为 100MB
MAX_UPLOAD_SIZE=104857600
```

### Q: 如何启用 HTTPS？
A:
1. 获取 SSL 证书和密钥
2. 在 `.env` 中配置：
   ```env
   ENABLE_HTTPS=true
   SSL_CERT_FILE=/path/to/cert.crt
   SSL_KEY_FILE=/path/to/key.key
   ```

### Q: JWT_SECRET_KEY 泄露了怎么办？
A:
1. 立即生成新密钥：`openssl rand -base64 32`
2. 更新 `.env` 文件
3. 重启应用
4. 建议要求所有用户重新登录

## 安全建议

1. **不要在版本控制中提交 `.env` 文件**
   - 将 `.env` 添加到 `.gitignore`
   - 只提交 `.env.example` 作为模板

2. **生产环境配置检查清单**
   - ✅ JWT_SECRET_KEY 已修改为安全密钥
   - ✅ 数据库使用强密码
   - ✅ ALLOWED_ORIGINS 限制为生产域名
   - ✅ GIN_MODE 设置为 release
   - ✅ LOG_FILE 已配置
   - ✅ 如需要，ENABLE_HTTPS 已启用

3. **敏感信息处理**
   - 不要在代码中硬编码密钥
   - 不要在日志中输出敏感信息
   - 定期轮换 JWT_SECRET_KEY

4. **权限管理**
   - `.env` 文件权限：`chmod 600 .env`
   - 日志目录权限：`chmod 755 logs`
   - 上传目录权限：`chmod 755 uploads`

## 相关文件

- **`.env`** - 开发环境配置（不提交到版本控制）
- **`.env.example`** - 配置模板（提交到版本控制）
- **`config/config.go`** - 配置加载代码
- **`config/config.yml`** - YAML 配置文件（可选）
- **`main.go`** - 应用入口，包含环境变量检查
- **`docker-compose.yml`** - Docker 环境配置示例
