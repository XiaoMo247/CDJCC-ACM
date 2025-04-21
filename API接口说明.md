## 📌 管理员模块（Admin）

### `/api/admin/login`

描述：管理员登录
参数：

- `username`: 用户名
- `password`: 密码

------

### `/api/admin/register`

描述：管理员注册用户（批量注册）
权限：需要登录为管理员
参数：

- `users`: 批量注册的用户名和密码列表

------

### `/api/admin/me`

描述：获取当前管理员信息
权限：需要登录
返回：当前管理员的详细信息

------

### `/api/admin/add`

描述：新增管理员
权限：需要登录且有权限
参数：用户名、密码等

------

### `/api/admin/delete/:id`

描述：删除管理员
权限：需要登录且有权限
路径参数：

- `id`: 管理员 ID

------

### `/api/admin/list`

描述：获取所有管理员列表
权限：需要登录

------

## 👤 用户模块（User）

### `/api/user/login`

描述：用户登录（学号 + 密码）
 参数：

- `student_number`: 学号
- `password`: 密码

------

## 📢 公告模块（Announcement）

### `/api/announcement/list`

描述：获取公告列表
返回：公告数组

------

### `/api/announcement/create`

描述：创建公告
权限：管理员登录
参数：标题、内容、发布时间等

------

### `/api/announcement/update/:id`

描述：更新指定公告
权限：管理员登录
路径参数：公告 ID
参数：标题、内容等

------

### `/api/announcement/delete/:id`

描述：删除公告
权限：管理员登录
路径参数：公告 ID

------

## 📝 加入我们模块（Join Apply）

### `/api/join/send`

描述：提交加入申请（普通用户提交）
权限：登录后自动带上 `user_id`
参数：申请表信息（姓名、年级、专业、介绍等）

------

### `/api/join/my`

描述：获取当前用户自己的申请记录
权限：登录后查看自己提交的内容

------

### `/api/join/get`

描述：管理员查看所有申请
权限：管理员登录

------

### `/api/join/:id`

描述：查看指定申请信息（管理员查看）
路径参数：申请 ID

------

### `/api/join/ac/:id`

描述：通过申请
权限：管理员
路径参数：申请 ID

------

### `/api/join/wa/:id`

描述：拒绝申请
权限：管理员
路径参数：申请 ID

