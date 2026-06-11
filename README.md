# GoForum

> 一个基于 **Go + Gin** 构建的轻量级论坛后端服务，支持帖子发布、社区分类、用户投票、JWT 认证等核心功能。可作为中小型社区/论坛系统的后端模板使用。

## 技术栈

| 层 | 技术 |
|------|------|
| 框架 | Gin |
| 数据库 | MySQL + Redis |
| 认证 | JWT (bcrypt) |
| 文档 | Swagger |
| 部署 | Docker |
| 测试 | k6 + Go testing |

## 项目结构

```
.
├── controller/       # HTTP 请求处理器
├── dao/
│   ├── mysql/        # MySQL 数据访问
│   └── redis/        # Redis 数据访问
├── logic/            # 业务逻辑层
├── models/           # 数据模型与参数定义
├── pkg/              # 工具包（JWT、雪花ID、错误码等）
├── route/            # 路由配置
├── middlewire/       # 中间件（JWT 认证、限流）
├── docs/             # Swagger 自动生成文档
├── logger/           # 日志模块
├── settings/         # 配置管理
├── bluebell_frontend/dist/  # 前端静态页面
└── k6_test.js        # k6 压测脚本
```

## 快速开始

### 前置要求

- Go 1.25+
- MySQL 8.0+
- Redis 7+

### 本地运行

```bash
# 1. 克隆项目
git clone <repo-url>
cd GoForum

# 2. 修改配置
#    编辑 config.yaml，确认 MySQL 和 Redis 连接信息

# 3. 初始化数据库
mysql -u root -p goforum < models/init.sql

# 4. 启动服务
go run .
```

### Docker 运行

```bash
# 构建并启动所有服务
docker compose up -d

# 查看日志
docker compose logs -f
```

### 访问

| 地址 | 说明 |
|------|------|
| http://localhost:8084 | 应用首页 |
| http://localhost:8084/swagger/index.html | Swagger API 文档 |

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/signup | 用户注册 | 否 |
| POST | /api/v1/login | 用户登录 | 否 |
| GET | /api/v1/community | 社区列表 | 否 |
| GET | /api/v1/community/:id | 社区详情 | 否 |
| GET | /api/v1/posts2 | 帖子列表（排序/筛选） | 否 |
| GET | /api/v1/post/:id | 帖子详情 | 否 |
| POST | /api/v1/post | 创建帖子 | 是 |
| POST | /api/v1/vote | 帖子投票 | 是 |

## 测试

```bash
# 运行所有单元测试
go test ./pkg/... ./controller/... -count=1

# 运行压测
k6 run k6_test.js
```

## 项目亮点

- ✅ Swagger 自动生成 API 文档
- ✅ 单元测试覆盖 pkg 和 controller 层（46+ 用例）
- ✅ k6 压测脚本，混合场景下 p95 < 50ms
- ✅ Docker 一键部署
- ✅ JWT 认证 + 接口限流
