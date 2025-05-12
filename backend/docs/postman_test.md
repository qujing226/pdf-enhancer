# Postman测试指南

## 文件上传测试

### 1. 配置请求

1. 打开Postman，创建新的POST请求
2. 设置请求URL：`http://localhost:8080/api/v1/reports/upload`
3. 在Headers中添加：
   - Authorization: Bearer {你的JWT令牌}

### 2. 设置请求体

1. 选择请求体类型为 `form-data`
2. 添加文件字段：
   - Key: `file`
   - Value: 选择要上传的PDF文件
   - 确保在Key的右侧下拉菜单中选择 `File` 类型

### 3. 发送请求

1. 点击 "Send" 按钮发送请求
2. 服务器将返回上传结果，包含：
   - 状态码：201（创建成功）
   - 响应体：包含文件ID和其他相关信息

### 4. 常见问题

1. 401 Unauthorized
   - 检查JWT令牌是否正确
   - 确保令牌未过期

2. 400 Bad Request
   - 确保使用了正确的字段名（file）
   - 检查文件格式是否为PDF

3. 413 Payload Too Large
   - 检查文件大小是否超过服务器限制

### 5. 示例响应

成功响应：
```json
{
    "code": 201,
    "message": "上传成功",
    "data": {
        "report_id": "123456789",
        "filename": "example.pdf",
        "size": 1024576
    }
}
```

错误响应：
```json
{
    "code": 400,
    "message": "无效的文件格式",
    "data": null
}
```