# CDJCC-ACM

成都锦城学院 ACM 集训队官网，包含团队介绍、公告管理、比赛信息、排行榜、常见问题等模块。

## 技术栈

- 前端：Vue 3 + Vite + Node.js
- 后端：Golang + Gin
- 数据库：MySQL
- 部署：Docker

## 快速开始

### 前端启动
```bash
cd fronted-acm
npm install
npm run dev
```

### 后端启动
```bash
cd acm-site
go run main.go
```

## 文档导航

### 配置相关
- [配置管理总结](CONFIGURATION_SUMMARY.md) - 前后端配置管理的完整指南
- [前端配置指南](fronted-acm/CONFIG_GUIDE.md) - 前端 .env 配置详细说明
- [后端配置指南](acm-site/BACKEND_CONFIG_GUIDE.md) - 后端 .env 配置详细说明

### 功能设计
- [需求分析](需求分析.md) - 项目功能需求列表
- [API接口说明](API接口说明.md) - API 接口文档

### 技术实施
- [统一 Token 系统实施总结](统一Token系统实施总结.md) - Token 认证系统重构总结
- [Docker 部署指南](关于Docker部署.md) - Docker 容器化部署说明

### 工具库
- [响应格式工具](acm-site/utils/response/README.md) - 后端响应格式统一工具

## 主要功能

- 团队介绍与展示
- 公告管理系统
- 比赛信息发布
- 排行榜展示
- 常见问题解答
- 用户认证与管理
- 申请加入团队功能
