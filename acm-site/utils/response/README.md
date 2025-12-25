# 统一响应格式工具

## 说明

这个包提供了统一的 API 响应格式，帮助保持代码一致性。

## 使用方法

```go
import "acm-site/utils/response"

// 成功响应
func GetUserInfo(c *gin.Context) {
    user := getUserFromDB()
    response.Success(c, user)
}

// 成功响应（自定义消息）
func UpdateUser(c *gin.Context) {
    // ... 更新逻辑
    response.SuccessWithMessage(c, "用户信息更新成功", nil)
}

// 错误响应
func Login(c *gin.Context) {
    if err := validateCredentials(); err != nil {
        response.Unauthorized(c, "用户名或密码错误")
        return
    }
}
```

## 响应格式

所有响应都遵循统一格式：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

## 快捷方法

- `Success(c, data)` - 成功响应
- `SuccessWithMessage(c, msg, data)` - 成功响应（自定义消息）
- `BadRequest(c, msg)` - 400 错误
- `Unauthorized(c, msg)` - 401 未授权
- `Forbidden(c, msg)` - 403 禁止访问
- `NotFound(c, msg)` - 404 未找到
- `InternalServerError(c, msg)` - 500 服务器错误

## 注意事项

⚠️ **当前项目中的旧 API 仍使用不同的响应格式**

为了保持向后兼容，现有 API 暂不修改。新增 API 建议使用此统一格式。

### 现有格式

当前项目中混用了多种格式：

```go
// 格式 1
gin.H{"message": "成功"}

// 格式 2
gin.H{"msg": "成功"}

// 格式 3
gin.H{"code": 200, "msg": "成功", "data": ...}
```

### 迁移建议

如果要迁移到统一格式，需要：
1. 逐个更新 API 接口
2. 同步更新前端代码
3. 进行充分测试

建议在下一个大版本中进行统一迁移。
