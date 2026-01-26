# 酒店行李寄存系统（后端）

基于 Go + Gin + GORM 的酒店行李寄存后端服务，提供用户登录、行李寄存、取件码管理、寄存室管理、行李迁移等功能接口。

## 功能概览
- 用户登录（bcrypt 校验）
- 创建用户（bcrypt 哈希，角色固定为 staff）
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

如需新增照片 URL 字段，请执行：
```sql
ALTER TABLE luggage_items ADD COLUMN photo_url VARCHAR(255) NULL;
ALTER TABLE luggage_history ADD COLUMN photo_url VARCHAR(255) NULL;
```

新增多图字段，请执行：
```sql
ALTER TABLE luggage_items ADD COLUMN photo_urls TEXT NULL;
ALTER TABLE luggage_history ADD COLUMN photo_urls TEXT NULL;
```

新增“行李寄存信息修改表”（用于 /api/luggage/logs/updated），请执行：
```sql
CREATE TABLE IF NOT EXISTS `行李寄存信息修改表` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `hotel_id` BIGINT NOT NULL,
  `luggage_id` BIGINT NOT NULL,
  `updated_by` VARCHAR(50) NOT NULL,
  `old_data` TEXT NOT NULL,
  `new_data` TEXT NOT NULL,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

可选：配置 Redis（用于缓存取件码查询）
```bat
set "REDIS_ADDR=127.0.0.1:6379"
set "REDIS_PASSWORD="
set "REDIS_DB=0"
```

### 3) 启动服务
```bat
cd /d C:\Users\32660\workspace\Winter_Vacation_Intensive_Training_Program\hotel_luggage
go run ./cmd
```

看到 `Listening and serving HTTP on 10.154.101.161:8080` 即启动成功。

## 常用接口（新结构）

接口统一前缀：`/api`

key-value
redis.get( )
redis.set( , , )


### 基础
- `GET /ping` 健康检查

### public 组（无需认证）
- `POST /api/login` 登录（返回 token）

### auth 组（需要登录，统一前缀 /api/luggage）
- `POST /api/luggage` 行李寄存
- `GET /api/luggage/by_code` 按取件码查询
- `POST /api/luggage/:id/checkout` 确认取件（id 为取件码，取件人自动使用登录账号）
- `GET /api/luggage/:id/checkout` 获取当前酒店有行李在存的客人名单
- `GET /api/luggage/list/by_guest_name` 查询某客人正在寄存的行李
- `PUT /api/luggage/:id` 修改寄存信息
- `GET /api/luggage/storerooms` 获取当前酒店所有寄存室
- `GET /api/luggage/storerooms/:id/orders` 获取该寄存室所有行李订单
- `POST /api/luggage/storerooms` 增加寄存室（自动使用当前用户的 hotel_id）
- `PUT /api/luggage/storerooms/:id` 软删除/停用寄存室
- `GET /api/luggage/logs/stored` 获取当前酒店寄存记录
- `GET /api/luggage/logs/updated` 获取当前酒店寄存信息修改记录
- `GET /api/luggage/logs/retrieved` 获取当前酒店取出记录



## 测试示例

首次创建用户建议使用命令行工具：
```bat
go run ./cmd/create_user -u staff_user -p 123456 -h 1
```

### 登录并获取 Token
```bat
curl -X POST http://localhost:8080/api/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"staff_user\",\"password\":\"123456\"}"
```

登录 JSON 请求体示例：
```json
{
  "username": "staff_user",
  "password": "123456"
}
```

## POST 接口 JSON 结构

### 创建用户（角色固定为 staff）
```json
{
  "username": "staff_user",
  "password": "123456",
  "hotel_id": 1
}
```

### 创建寄存
```json
{
  "guest_name": "张三",
  "staff_name": "staff_user",
  "contact_phone": "13800000000",
  "contact_email": "zhangsan@example.com",
  "description": "黑色行李箱",
  "quantity": 1,
  "special_notes": "易碎品",
  "photo_url": "http://example.com/photo.jpg",
  "storeroom_id": 1,
  "qr_code_url": "/qr/xxxxxx"
}
```

### 修改寄存（PUT /api/luggage/:id）
```json
{
  "description": "黑色行李箱",
  "photo_url": "http://example.com/photo.jpg"
}
```

### 取件
无需 Body，取件人自动使用当前登录账号。

### 行李迁移
```json
{
  "to_storeroom_id": 2,
  "migrated_by": "staff_user"
}
```

### 创建寄存室
```json
{
  "name": "A区-1号",
  "location": "一楼A区",
  "capacity": 50,
  "is_active": true
}
```

返回中包含 `token`，后续请求需携带：
```
Authorization: Bearer <token>
```

### 携带 Token 调用受保护接口
示例（创建寄存室）：
```bat
curl -X POST http://localhost:8080/api/luggage/storerooms ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer <token>" ^
  -d "{\"name\":\"A区-1号\",\"location\":\"一楼A区\",\"capacity\":50,\"is_active\":true}"
```

### 创建寄存室
```bat
curl -X POST http://localhost:8080/api/luggage/storerooms ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer <token>" ^
  -d "{\"name\":\"A区-1号\",\"location\":\"一楼A区\",\"capacity\":50,\"is_active\":true}"
```

### 行李寄存
```bat
  curl -X POST http://localhost:8080/api/luggage ^
    -H "Content-Type: application/json" ^
    -H "Authorization: Bearer <token>" ^
    -d "{\"guest_name\":\"张三\",\"staff_name\":\"staff_user\",\"contact_phone\":\"13800000000\",\"description\":\"黑色行李箱\",\"quantity\":1,\"storeroom_id\":1}"
```

### 取件
```bat
curl -X POST http://localhost:8080/api/luggage/取件码/checkout ^
  -H "Authorization: Bearer <token>"
```

### 行李迁移
```bat
curl -X POST http://localhost:8080/api/luggage/1/transfer ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer <token>" ^
  -d "{\"to_storeroom_id\":2,\"migrated_by\":\"用户名\"}"
```

### 查询迁移历史
```bat
curl "http://localhost:8080/api/luggage/1/transfers" ^
  -H "Authorization: Bearer <token>"
```

### 查看二维码
```bat
curl http://localhost:8080/qr/取件码 --output qrcode.png
```

## 备注
- 若接口返回 404，请确认服务已重启并包含最新路由。
- 若数据库连接失败，请检查 `DB_DSN` 配置和 MySQL 账号密码。
