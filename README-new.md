# AI增强报告系统

## 项目简介

本项目是基于README.md中的笔试任务要求开发的一个完整的报告管理系统，实现了用户登录、报告查询和AI摘要生成功能。系统采用Gin+Vue+MySQL技术栈，通过Docker部署，并集成DeepSeek客户端用于AI摘要生成。

## 功能特点

- 用户登录认证（JWT Token）
- 报告列表查询
- 单份报告内容查询
- AI自动生成报告摘要（DeepSeek API）
- PDF文件安全存储与访问
- 完整的API文档（Swagger）

## 技术栈

- **前端**：Vue.js + Element UI
- **后端**：Go + Gin框架
- **数据库**：MySQL
- **AI服务**：DeepSeek API
- **文件服务器**：MinIO
- **容器化**：Docker + Docker Compose

## 安全特性

- 用户密码使用Argon2id算法加密
- JWT令牌采用RS256非对称加密
- PDF文件存储采用服务端加密
- 所有API通信使用TLS 1.3加密
- 基于用户身份的访问控制

## 快速开始

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+

### 部署步骤

1. 克隆代码库

```bash
git clone https://github.com/qujing226/pdf-enhancer.git
cd pdf-enhancer
```

2. 配置环境变量

```bash
cp .env.example .env
# 编辑.env文件，填入必要的密钥
```

3. 启动服务

```bash
docker-compose up -d
```

4. 访问服务

- 前端界面：http://localhost
- API文档：http://localhost:8080/swagger/index.html
- MinIO控制台：http://localhost:9001

## 项目结构

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

## API文档

详细的API文档可通过Swagger UI访问：http://localhost:8080/swagger/index.html

主要API包括：

- `POST /api/v1/login` - 用户登录
- `GET /api/v1/reports` - 查询报告列表
- `GET /api/v1/report/:report_id` - 查询单份报告
- `POST /api/v1/report/:report_id/summary` - 生成报告摘要
- `GET /api/v1/report/:report_id/pdf` - 获取报告PDF

## 开发文档

详细的开发文档请参考 [development.md](./development.md)，其中包含：

- 系统架构设计
- 数据库设计
- API详细说明
- 安全方案
- DeepSeek客户端集成
- 文件服务器配置
- 部署方案
- 开发指南

## 测试

项目提供了完整的测试套件：

```bash
# 运行后端单元测试
cd backend
go test ./...

# 运行前端单元测试
cd frontend
npm run test:unit
```

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

如有任何问题，请联系项目维护者。