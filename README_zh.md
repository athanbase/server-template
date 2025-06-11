# 服务器模板

基于 [Kratos](https://github.com/go-kratos/kratos) 框架的微服务项目模板，集成了常用的中间件和工具，提供了完整的开发流程和最佳实践。

## 项目特性

- 基于 Kratos v2 框架，提供完整的微服务架构
- 使用 Protocol Buffers 定义 API，支持 gRPC 和 HTTP 协议
- 集成 MySQL 主从读写分离
- 集成 Redis 缓存
- 使用 Wire 进行依赖注入
- 使用 sqlc 生成类型安全的数据库访问代码
- 使用 buf 工具管理 Protocol Buffers
- 支持数据库事务和错误处理
- 集成日志系统
- 提供 Docker 部署支持

## 技术栈

- Go 1.24+
- Kratos v2
- Protocol Buffers
- gRPC/HTTP
- MySQL
- Redis
- Wire
- sqlc
- buf
- Docker

## 项目结构

```
.
├── api                 # API 定义
│   └── server          # 服务 API 定义
├── cmd                 # 应用入口
│   └── server          # 服务入口
├── internal            # 内部代码
│   ├── biz             # 业务逻辑层
│   ├── conf            # 配置定义
│   ├── data            # 数据访问层
│   │   ├── migration   # 数据库迁移
│   │   ├── queries     # 生成的查询代码
│   │   └── query       # SQL 查询定义
│   ├── server          # 服务实现
│   └── service         # 服务接口实现
├── pkg                 # 公共包
├── Dockerfile          # Docker 构建文件
├── Makefile            # 构建脚本
├── config.yaml         # 配置文件
└── sqlc.yaml           # sqlc 配置
```

## 开发环境准备

### 安装依赖工具

```bash
# 初始化开发环境，安装所需工具
make init
```

这将安装以下工具：
- kratos CLI
- protoc-gen-go-http
- wire
- sqlc
- buf
- protoc-gen-openapi
- protoc-gen-go
- protoc-gen-go-grpc
- migrate

## 开发指南

### 生成代码

```bash
# 生成 API 相关代码
make api

# 生成配置相关代码
make config

# 生成所有代码
make all

# 生成依赖注入代码
make wire

# 生成数据库访问代码
make sqlc
```

### 创建数据库迁移

```bash
# 创建新的迁移文件
make new_migration name=migration_name
```

### 运行服务

```bash
# 运行服务
make server
```

### 构建项目

```bash
# 构建项目
make build
```

## Docker 部署

```bash
# 构建 Docker 镜像
docker build -t <your-docker-image-name> .

# 运行 Docker 容器
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## API 文档

项目自动生成 OpenAPI 规范文件 `openapi.yaml`，可以使用 Swagger UI 等工具查看 API 文档。

## 数据库设计

项目包含以下数据表：

- `user`: 用户基本信息
- `user_detail`: 用户详细信息

## 配置说明

配置文件位于 `config.yaml`，包含以下主要配置：

- 服务配置（HTTP/gRPC）
- 数据库配置（主从）
- Redis 配置
- 日志配置

## 许可证

[MIT License](LICENSE)