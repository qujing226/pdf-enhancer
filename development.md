# AI增强报告系统开发文档

## 1. 项目概述

本项目是一个基于Gin+Vue+MySQL的RESTful API系统，用于实现用户登录、报告查询和AI摘要生成功能。系统通过Docker部署，并集成DeepSeek客户端用于AI摘要生成。

## 2. 系统架构

### 2.1 整体架构

```
+------------------+      +------------------+      +------------------+
|                  |      |                  |      |                  |
|  前端 (Vue.js)   | <--> |  后端 (Gin)      | <--> |  数据库 (MySQL)  |
|                  |      |                  |      |                  |
+------------------+      +--------+---------+      +------------------+
                                   |
                                   v
                          +------------------+      +------------------+
                          |                  |      |                  |
                          |  DeepSeek API    | <--> |  文件服务器      |
                          |                  |      |                  |
                          +------------------+      +------------------+
```

### 2.2 技术栈选择

- **前端**：Vue.js + Element UI
- **后端**：Go + Gin框架
- **数据库**：MySQL
- **AI服务**：DeepSeek API
- **文件服务器**：MinIO
- **容器化**：Docker + Docker Compose
- **API文档**：Swagger

## 3. 数据库设计

### 3.1 用户表 (users)

```sql
CREATE TABLE `users` (
  `id` varchar(64) NOT NULL COMMENT '用户ID',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `salt` varchar(64) NOT NULL COMMENT '盐值',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 3.2 报告表 (reports)

```sql
CREATE TABLE `reports` (
  `id` varchar(64) NOT NULL COMMENT '报告ID',
  `user_id` varchar(64) NOT NULL COMMENT '用户ID',
  `title` varchar(255) NOT NULL COMMENT '报告标题',
  `content` text COMMENT '报告内容',
  `summary` text COMMENT 'AI生成的摘要',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `pdf_path` varchar(255) DEFAULT NULL COMMENT 'PDF文件路径',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_reports_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报告表';
```

## 4. API设计

### 4.1 用户登录

```
POST /api/v1/login
```

**请求参数**：

```json
{
  "email": "alice@example.com",
  "password": "password123"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "u123",
      "name": "Alice",
      "email": "alice@example.com"
    }
  }
}
```

### 4.2 查询报告列表

```
GET /api/v1/reports?user_id=xxx
```

**请求头**：

```
Authorization: Bearer {token}
```

**响应**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "reports": [
      {
        "report_id": "rpt001",
        "title": "Q1 2025 Investment Report",
        "created_at": "2025-04-01",
        "has_summary": true
      },
      {
        "report_id": "rpt002",
        "title": "Q4 2024 Investment Report",
        "created_at": "2025-01-01",
        "has_summary": false
      }
    ]
  }
}
```

### 4.3 查询单份报告

```
GET /api/v1/report/:report_id
```

**请求头**：

```
Authorization: Bearer {token}
```

**响应**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "report_id": "rpt001",
    "user_id": "u123",
    "title": "Q1 2025 Investment Report",
    "content": "这是完整报告内容，可以很长很长...",
    "summary": "这是由 GPT 模拟产生的摘要内容。",
    "created_at": "2025-04-01",
    "pdf_url": "/api/v1/report/rpt001/pdf"
  }
}
```

### 4.4 生成报告摘要

```
POST /api/v1/report/:report_id/summary
```

**请求头**：

```
Authorization: Bearer {token}
```

**响应**：

```json
{
  "code": 200,
  "message": "摘要生成成功",
  "data": {
    "summary": "这是由DeepSeek AI生成的报告摘要内容..."
  }
}
```

### 4.5 获取报告PDF

```
GET /api/v1/report/:report_id/pdf
```

**请求头**：

```
Authorization: Bearer {token}
```

**响应**：

PDF文件二进制流，带有适当的Content-Type和Content-Disposition头。

## 5. 安全方案

### 5.1 用户信息加密

1. **密码安全**
   - 使用Argon2id算法进行密码哈希（比bcrypt和PBKDF2更安全）
   - 为每个用户生成唯一的盐值(salt)
   - 设置适当的内存、迭代次数和并行度参数

2. **JWT认证**
   - 使用RS256算法（非对称加密）签名JWT令牌
   - 设置合理的过期时间（如1小时）
   - 实现令牌刷新机制
   - 在Redis中维护黑名单，处理已注销的令牌

3. **敏感数据保护**
   - 数据库中的敏感字段使用AES-256-GCM加密
   - 加密密钥使用KMS（密钥管理服务）管理

### 5.2 PDF文件保护

1. **存储安全**
   - PDF文件存储在MinIO对象存储服务中
   - 使用服务端加密(SSE)保护存储的文件
   - 为每个文件生成唯一的加密密钥

2. **访问控制**
   - 实现基于用户身份的访问控制(RBAC)
   - PDF文件URL使用签名URL，有效期短（如15分钟）
   - 实现文件访问审计日志

3. **传输安全**
   - 所有API通信使用TLS 1.3加密
   - 实现HTTP安全头（如HSTS, CSP等）
   - 针对PDF文件传输启用内容加密

## 6. DeepSeek客户端集成

### 6.1 DeepSeek API配置

```go
type DeepSeekConfig struct {
    APIKey      string
    BaseURL     string
    ModelName   string
    MaxTokens   int
    Temperature float64
}

func NewDeepSeekClient(config DeepSeekConfig) *DeepSeekClient {
    // 初始化DeepSeek客户端
}
```

### 6.2 摘要生成实现

```go
func (c *DeepSeekClient) GenerateSummary(reportContent string) (string, error) {
    // 构建提示词
    prompt := fmt.Sprintf("请为以下报告生成一个简洁的摘要（不超过200字）:\n\n%s", reportContent)
    
    // 调用DeepSeek API
    response, err := c.client.CreateCompletion(context.Background(), &deepseek.CompletionRequest{
        Model:       c.config.ModelName,
        Prompt:      prompt,
        MaxTokens:   c.config.MaxTokens,
        Temperature: c.config.Temperature,
    })
    
    if err != nil {
        return "", fmt.Errorf("调用DeepSeek API失败: %w", err)
    }
    
    return response.Choices[0].Text, nil
}
```

## 7. 文件服务器配置

选择MinIO作为文件服务器，理由如下：

1. **兼容S3 API**：MinIO兼容Amazon S3 API，便于迁移和集成
2. **轻量级**：适合容器化部署，资源消耗低
3. **高性能**：支持高并发读写，适合PDF文件存储和检索
4. **安全特性**：支持服务端加密、对象锁定和版本控制
5. **可扩展性**：支持分布式部署，可随业务增长扩展

### 7.1 MinIO Docker配置

```yaml
services:
  minio:
    image: minio/minio:RELEASE.2023-09-30T07-02-29Z
    container_name: minio
    ports:
      - "9000:9000"  # API端口
      - "9001:9001"  # 控制台端口
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

volumes:
  minio_data:
```

## 8. 部署方案

### 8.1 Docker Compose配置

```yaml
version: '3.8'

services:
  # 前端服务
  frontend:
    build: ./frontend
    container_name: report-frontend
    ports:
      - "80:80"
    depends_on:
      - backend

  # 后端API服务
  backend:
    build: ./backend
    container_name: report-backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=reports_db
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=minioadmin
      - MINIO_SECRET_KEY=minioadmin
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
      - JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY}
      - JWT_PUBLIC_KEY=${JWT_PUBLIC_KEY}
    depends_on:
      mysql:
        condition: service_healthy
      minio:
        condition: service_healthy

  # MySQL数据库
  mysql:
    image: mysql:8.0
    container_name: report-mysql
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=reports_db
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  # MinIO对象存储
  minio:
    image: minio/minio:RELEASE.2023-09-30T07-02-29Z
    container_name: report-minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

volumes:
  mysql_data:
  minio_data:
```

### 8.2 部署步骤

1. **环境准备**
   - 安装Docker和Docker Compose
   - 准备.env文件，包含必要的环境变量

2. **构建和启动**
   ```bash
   # 克隆代码库
   git clone https://github.com/yourusername/ai-enhance-test.git
   cd ai-enhance-test
   
   # 创建.env文件
   cp .env.example .env
   # 编辑.env文件，填入必要的密钥
   
   # 启动服务
   docker-compose up -d
   ```

3. **初始化数据**
   ```bash
   # 导入示例数据
   docker exec -i report-mysql mysql -uroot -ppassword reports_db < sample_data.sql
   ```

4. **验证部署**
   - 访问 http://localhost 查看前端界面
   - 访问 http://localhost:8080/swagger/index.html 查看API文档
   - 访问 http://localhost:9001 查看MinIO控制台

## 9. 开发指南

### 9.1 本地开发环境设置

1. **后端开发**
   ```bash
   # 进入后端目录
   cd backend
   
   # 安装依赖
   go mod tidy
   
   # 启动开发服务器
   go run main.go
   ```

2. **前端开发**
   ```bash
   # 进入前端目录
   cd frontend
   
   # 安装依赖
   npm install
   
   # 启动开发服务器
   npm run dev
   ```

### 9.2 测试

1. **单元测试**
   ```bash
   # 运行后端单元测试
   cd backend
   go test ./...
   
   # 运行前端单元测试
   cd frontend
   npm run test:unit
   ```

2. **API测试**
   - 使用Postman导入提供的集合文件进行API测试
   - 或使用Swagger UI进行交互式API测试

## 10. 项目结构

```
.
├── backend/                 # 后端Go代码
│   ├── api/                 # API处理器
│   ├── config/              # 配置文件
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   ├── services/            # 业务逻辑
│   ├── utils/               # 工具函数
│   └── main.go              # 入口文件
├── frontend/                # 前端Vue代码
│   ├── public/              # 静态资源
│   ├── src/                 # 源代码
│   └── package.json         # 依赖配置
├── docker/                  # Docker相关文件
├── scripts/                 # 脚本文件
├── .env.example            # 环境变量示例
├── docker-compose.yml      # Docker Compose配置
├── init.sql                # 数据库初始化SQL
└── README.md               # 项目说明
```

## 11. 总结

本项目实现了一个完整的报告管理系统，包括用户认证、报告查询和AI摘要生成功能。通过使用Gin+Vue+MySQL技术栈，结合DeepSeek AI和MinIO文件服务器，构建了一个安全、高效的应用系统。项目采用Docker容器化部署，便于在各种环境中快速启动和运行。

在安全方面，项目采用了多层次的保护措施，包括密码哈希、JWT认证、数据加密和文件访问控制，确保用户数据和报告内容的安全性。

未来可以考虑添加更多功能，如报告协作编辑、版本控制、更高级的AI分析等，进一步提升系统的实用性和用户体验。