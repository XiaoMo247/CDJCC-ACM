# 构建阶段
FROM node:22.12.0-alpine AS build

# 忽略证书错误
RUN npm config set strict-ssl false

# 安装 cnpm
RUN npm install -g cnpm --registry=https://registry.npm.taobao.org

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到容器中的工作目录
COPY . .

# 清理并重新安装依赖，确保所有依赖包正确安装
RUN rm -rf node_modules && cnpm install

# 构建应用
RUN cnpm run build

# 生产环境阶段
FROM nginx:alpine

# 将构建出来的前端文件复制到 nginx 的目录中
COPY --from=build /app/dist /usr/share/nginx/html

# 将自定义的 nginx 配置文件复制到容器中
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# 暴露 80 端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]