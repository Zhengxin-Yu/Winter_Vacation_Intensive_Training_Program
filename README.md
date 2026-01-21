# 酒店行李寄存系统（后端）

基于 Go + Gin + GORM 的酒店行李寄存后端服务，提供用户登录、行李寄存、取件码管理、寄存室管理、行李迁移等功能接口。

## 功能概览
- 用户登录（bcrypt 校验）
- 创建用户（bcrypt 哈希）
- 行李寄存（生成取件码 + 可选二维码 URL）
- 行李查询（按姓名/手机号/取件码）
- 寄存单列表（按用户/客人）
- 寄存单详情（按 ID/取件码/手机号）
- 取件功能（更新状态/取件人/取件时间）
- 寄存室管理（列表/创建/删除/状态更新）
- 行李迁移（更新寄存室 + 迁移日志）
- 二维码生成与展示（PNG）
- 首页功能入口（接口清单）
- 修改寄存信息/取件码
- 行李绑定（将行李绑定到用户）

## 环境依赖
- Go 1.20+
- MySQL 5.7/8.0

## 目录结构
```
hotel_luggage/
├── cmd/                # 程序入口
│   ├── main.go
│   └── create_user/    # 命令行创建用户工具
├── configs/            # 配置
├── internal/
│   ├── models/         # 数据模型
│   ├── handlers/       # 接口处理
│   ├── services/       # 业务逻辑
│   ├── repositories/   # 数据访问
│   └── middleware/
├── router/             # 路由
└── utils/              # 工具
```

## 快速开始

### 1) 创建数据库并导入表结构
在 MySQL 中执行：
```sql
CREATE DATABASE IF NOT EXISTS hotel_luggage DEFAULT CHARSET utf8mb4;
USE hotel_luggage;
-- 导入你的 hotel_luggage_system.sql
```

### 2) 配置数据库连接
在 Windows CMD 中设置环境变量（注意使用引号）：
```bat
set "DB_DSN=root:123456@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3) 启动服务
```bat
cd /d C:\Users\32660\workspace\Winter_Vacation_Intensive_Training_Program\hotel_luggage
go run ./cmd
```

看到 `Listening and serving HTTP on :8080` 即启动成功。

## 常用接口

### 基础
- `GET /ping` 健康检查
- `GET /home` 首页功能入口

### 用户
- `POST /users` 创建用户
- `POST /login` 登录

### 行李寄存
- `POST /storage` 创建寄存记录
- `POST /storage/retrieve` 取件
- `GET /storage/search` 按姓名/手机号查询
- `GET /storage/by-code` 按取件码查询
- `PUT /storage/:id` 修改寄存信息
- `PUT /storage/:id/code` 修改取件码
- `POST /storage/bind` 行李绑定到用户

### 寄存单
- `GET /storage/list` 按用户查询列表
- `GET /storage/list/by-guest` 按客人姓名/手机号查询列表
- `GET /storage/detail` 按 ID 查询详情
- `GET /storage/detail/by-code` 按取件码查询详情
- `GET /storage/detail/by-phone` 按手机号查询详情

### 取件码
- `GET /pickup-codes` 按用户查询取件码列表
- `GET /pickup-codes/by-phone` 按手机号查询取件码列表

### 寄存室管理
- `GET /storerooms` 列表
- `POST /storerooms` 创建
- `DELETE /storerooms/:id` 删除（有行李不能删）
- `PUT /storerooms/:id/status` 更新状态
- `POST /storerooms/migrate` 行李迁移

### 二维码
- `GET /qr/:code` 生成并返回二维码 PNG

### 前端测试页
- `frontend_testing/index.html` 用于简单接口测试

## 测试示例

### 创建用户
```bat
curl -X POST http://localhost:8080/users ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"admin\",\"password\":\"123456\",\"role\":\"admin\"}"
```

### 创建寄存室
```bat
curl -X POST http://localhost:8080/storerooms ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"A区-1号\",\"location\":\"一楼A区\",\"capacity\":50,\"is_active\":true}"
```

### 行李寄存
```bat
curl -X POST http://localhost:8080/storage ^
  -H "Content-Type: application/json" ^
  -d "{\"guest_name\":\"张三\",\"contact_phone\":\"13800000000\",\"description\":\"黑色行李箱\",\"quantity\":1,\"storeroom_id\":1,\"stored_by\":1}"
```

### 取件
```bat
curl -X POST http://localhost:8080/storage/retrieve ^
  -H "Content-Type: application/json" ^
  -d "{\"code\":\"取件码\",\"retrieved_by\":1}"
```

### 查看二维码
```bat
curl http://localhost:8080/qr/取件码 --output qrcode.png
```

## 备注
- 若接口返回 404，请确认服务已重启并包含最新路由。
- 若数据库连接失败，请检查 `DB_DSN` 配置和 MySQL 账号密码。
