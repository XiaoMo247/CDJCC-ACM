## 项目部署总结

### 前端创建 Dockerfile

```dockerfile
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
```

### 后端创建 Dockerfile

```dockerfile
# 使用官方 Go 镜像
FROM golang:1.24.1

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制项目所有文件
COPY . .

# 构建项目
RUN go build -o main .

# 开放端口（替换为你的后端端口）
EXPOSE 8081

# 启动命令
CMD ["./main"]
```

### 配置前端 Nginx.conf

```conf
server {
  listen 80;
  server_name localhost;

  root /usr/share/nginx/html;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }
}
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  backend:
    build: ./acm-site
    ports:
      - "8081:8081"
    restart: always
    depends_on:
      - mysql
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=123456
      - DB_NAME=acm_site

  frontend:
    build: ./jingcheng-acm-web
    ports:
      - "80:80"
    restart: always

  mysql:
    image: mysql:8.0
    container_name: acm-mysql
    restart: always
    ports:
      - "3307:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: acm_site
    volumes:
      - mysql-data:/var/lib/mysql

volumes:
  mysql-data:
```

### 本地打包测试（Windows）

```bash
# 切换到根目录
cd CDJCC-ACM

# 构建所有服务
docker-compose build

# 启动
docker-compose up -d
```

### 远端服务器部署

1. **导出 Docker 镜像为 tar 文件**

   使用 `docker save` 命令将本地的 Docker 镜像导出为一个 `.tar` 文件：

   ```bash
   docker save -o cdjccacm-backend.tar cdjccacm-backend
   docker save -o cdjccacm-frontend.tar cdjccacm-frontend
   docker save -o mysql.tar mysql
   ```

2. **将 tar 文件传输到远端服务器**

   使用 `scp` 或 `rsync` 命令将 `.tar` 文件传输到远端服务器：

   ```bash
   scp cdjccacm-backend.tar username@your_server_ip:/path/to/your/folder
   scp cdjccacm-frontend.tar username@your_server_ip:/path/to/your/folder
   scp mysql.tar username@your_server_ip:/path/to/your/folder
   ```

3. **在远端服务器上加载 Docker 镜像**

   登录远端服务器后，使用 `docker load` 命令加载上传的 `.tar` 文件：

   ```bash
   docker load -i /path/to/your/folder/cdjccacm-backend.tar
   docker load -i /path/to/your/folder/cdjccacm-frontend.tar
   docker load -i /path/to/your/folder/mysql.tar
   ```

4. **启动 Docker 容器**

   加载镜像后，可以用 `docker run` 或者 `docker-compose` 启动容器，和本地一样进行部署：

   远端docker-compose：

   ```yaml
   version: '3'
   
   services:
     mysql:
       image: mysql:8.0
       container_name: acm-mysql
       environment:
         MYSQL_ROOT_PASSWORD: 密码
         MYSQL_DATABASE: acm_site
       networks:
         - app-net
       ports:
         - "3306:3306"
       restart: always
   
     backend:
       image: cdjccacm-backend
       container_name: acm-backend
       environment:
         DB_HOST: 服务器IP地址
         DB_PORT: 3306
         DB_USER: root
         DB_PASSWORD: 密码
         DB_NAME: acm_site
       networks:
         - app-net
       ports:
         - "8081:8081"
       depends_on:
         - mysql
       restart: always
   
     frontend:
       image: cdjccacm-frontend
       container_name: acm-frontend
       networks:
         - app-net
       ports:
         - "80:80"
       restart: always
   
   networks:
     app-net:
       driver: bridge
   ```

   启动：

   ```bash
   docker-compose up -d
   ```

------

### docker 和 docker-compose 常用指令

- **查看所有镜像**

  ```bash
  docker images
  ```

- **查看容器列表**

  ```bash
  docker ps
  ```

- **停止容器**

  ```bash
  docker stop <container_id>
  ```

- **启动容器**

  ```bash
  docker start <container_id>
  ```

- **删除容器**

  ```bash
  docker rm <container_id>
  ```

- **删除镜像**

  ```bash
  docker rmi <image_id>
  ```

- **查看容器日志**

  ```bash
  docker logs <container_id>
  ```

- **进入容器内部**

  ```bash
  docker exec -it <container_id> /bin/sh
  ```

- **构建镜像**

  ```bash
  docker-compose build
  ```

- **启动容器**

  ```bash
  docker-compose up -d
  ```

- **查看 Docker Compose 服务状态**

  ```bash
  docker-compose ps
  ```