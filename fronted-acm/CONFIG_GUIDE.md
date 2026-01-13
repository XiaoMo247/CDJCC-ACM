# 项目配置指南

本文档说明了项目中所有配置项的使用方法。所有配置项都集中在 `.env` 文件中管理。

## 快速开始

1. 复制 `.env.example` 为 `.env`：
   ```bash
   cp .env.example .env
   ```

2. 根据实际环境修改 `.env` 文件中的配置

3. 重启开发服务器使配置生效

## 配置说明

### 1. API 配置

```env
VITE_API_BASE_URL=http://localhost:8081/api
```

- **说明**：API 服务器的基础 URL
- **开发环境**：`http://localhost:8081/api`
- **生产环境示例**：`http://your-domain.com/api`
- **使用位置**：`src/utils/request.js`

### 2. 上传文件服务器配置

```env
VITE_UPLOAD_TARGET=http://localhost:8081
VITE_UPLOAD_PATH=/uploads
```

- **说明**：上传文件时的目标服务器和路径
- **VITE_UPLOAD_TARGET**：文件上传服务器地址
- **VITE_UPLOAD_PATH**：上传路径前缀
- **使用位置**：`vite.config.js` 中的代理配置

### 3. Token 配置

```env
VITE_TOKEN_REFRESH_INTERVAL=3000000
VITE_TOKEN_KEY=auth_token
VITE_USER_INFO_KEY=user_info
```

- **VITE_TOKEN_REFRESH_INTERVAL**：Token 自动刷新间隔（单位：毫秒，默认 50 分钟）
- **VITE_TOKEN_KEY**：localStorage 中存储 token 的 key
- **VITE_USER_INFO_KEY**：localStorage 中存储用户信息的 key
- **使用位置**：`src/utils/tokenManager.js`

### 4. 外部链接配置

```env
VITE_QQ_GROUP_URL=https://jq.qq.com/?_wv=1027&k=KoVRvLvL
VITE_BILIBILI_URL=https://space.bilibili.com/3546651937475184?spm_id_from=333.337.search-card.all.click
VITE_OJ_URL=https://hydro.ac/d/cdjcc_acm_2333/
```

- **VITE_QQ_GROUP_URL**：QQ 群链接，用于忘记密码页面
- **VITE_BILIBILI_URL**：Bilibili 频道链接，用于知识库页面
- **VITE_OJ_URL**：在线 OJ 平台链接，用于知识库页面
- **使用位置**：
  - `src/pages/ForgotPassword.vue`
  - `src/pages/KnowledgePage.vue`

### 5. 默认头像配置

```env
VITE_DEFAULT_AVATAR=https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png
```

- **说明**：用户头像加载失败时的默认图片 URL
- **使用位置**：`src/pages/HonorWall.vue`

### 6. 应用信息

```env
VITE_APP_NAME=CDJCC ACM 集训队官网
```

- **说明**：应用名称（预留用于今后扩展）
- **使用位置**：可在各组件中通过 `import.meta.env.VITE_APP_NAME` 访问

## 环境变量访问方式

在 Vue 组件中访问环境变量：

```javascript
// 在 data() 或 computed 中
data() {
  return {
    apiUrl: import.meta.env.VITE_API_BASE_URL,
    qqGroupUrl: import.meta.env.VITE_QQ_GROUP_URL
  }
}

// 或在方法中
methods: {
  someMethod() {
    const apiUrl = import.meta.env.VITE_API_BASE_URL
  }
}
```

在 JavaScript 文件中访问：

```javascript
// 在 request.js 等工具文件中
const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
```

在配置文件中访问：

```javascript
// 在 vite.config.js 中
const uploadTarget = process.env.VITE_UPLOAD_TARGET || 'http://localhost:8081'
```

## 部署指南

### 开发环境

直接使用 `.env` 文件，配置本地开发服务器地址：

```env
VITE_API_BASE_URL=http://localhost:8081/api
VITE_UPLOAD_TARGET=http://localhost:8081
```

### 生产环境

创建 `.env.production` 文件（可选），配置生产服务器地址：

```env
VITE_API_BASE_URL=https://your-domain.com/api
VITE_UPLOAD_TARGET=https://your-domain.com
VITE_QQ_GROUP_URL=https://jq.qq.com/?_wv=1027&k=KoVRvLvL
VITE_BILIBILI_URL=https://space.bilibili.com/3546651937475184
VITE_OJ_URL=https://your-oj.com/path
```

### 自动构建

如果使用 CI/CD 流程，可以在构建时注入环境变量：

```bash
# 构建前设置环境变量
export VITE_API_BASE_URL=https://your-domain.com/api
export VITE_UPLOAD_TARGET=https://your-domain.com

# 执行构建
npm run build
```

## 安全建议

1. **不要在版本控制中提交 `.env` 文件**
   - 将 `.env` 添加到 `.gitignore`
   - 只提交 `.env.example` 作为模板

2. **敏感信息处理**
   - 不要在 `.env` 中存储密钥或密码
   - 生产环境的敏感配置应通过环境变量或密钥管理工具注入

3. **定期检查配置**
   - 确保生产环境配置指向正确的服务器
   - 不要使用开发环境的 URL 进行生产部署

## 常见问题

### Q: 修改 `.env` 文件后，变化不生效？
A: 需要重启开发服务器（`npm run dev`）。Vite 在启动时读取环境变量。

### Q: 如何在不同环境使用不同配置？
A: 创建对应的环境文件：
- `.env.development` - 开发环境
- `.env.production` - 生产环境
- `.env.staging` - 测试环境

### Q: 环境变量前缀为什么是 `VITE_`？
A: 这是 Vite 的安全机制。只有以 `VITE_` 开头的变量才会被暴露给客户端代码，防止敏感信息泄露。

## 相关文件

- **`.env`** - 开发环境配置（不提交到版本控制）
- **`.env.example`** - 配置模板（提交到版本控制）
- **`vite.config.js`** - Vite 构建配置
- **`src/utils/request.js`** - API 请求工具
- **`src/utils/tokenManager.js`** - Token 管理工具
