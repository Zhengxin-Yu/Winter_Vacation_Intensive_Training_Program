# 前端对接指南

本指南用于说明如何把前端与当前后端接口对接。

---

## 1. 启动后端服务

在 Windows CMD 中执行：

```bat
cd /d C:\Users\32660\workspace\Winter_Vacation_Intensive_Training_Program\hotel_luggage
set "DB_DSN=root:你的密码@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd
```

看到 `Listening and serving HTTP on 10.154.101.161:8080` 说明启动成功。

---

## 2. 前端对接核心概念

- **接口地址（Base URL）**：`http://10.154.101.161:8080`
- **请求方式**：GET/POST/PUT/DELETE
- **请求体格式**：JSON
- **响应格式**：JSON

前端通过 `fetch` / `axios` 访问接口即可。
需要登录后携带 `Authorization: Bearer <token>` 才能访问 `/api` 下的受保护接口。

---

## 3. 接口清单与对接说明

说明：
- `GET /ping` 和 `POST /api/login` 不需要登录。
- 其余 `/api/luggage/*` 需要请求头 `Authorization: Bearer <token>`。

### 3.1 GET /ping
作用：健康检查。
请求：无参数。  
响应示例：
```json
{
  "message": "pong"
}
```

### 3.2 POST /api/login
作用：登录获取 token。  
Body 示例：
```json
{
  "username": "staff_user",
  "password": "123456"
}
```
响应示例：
```json
{
  "message": "login success",
  "user": {
    "id": 1,
    "username": "staff_user",
    "role": "staff",
    "hotel_id": 1
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```
响应示例：
```json
{
  "message": "login failed",
  "error": "invalid username or password"
}
```

### 3.3 POST /api/luggage
作用：创建寄存单。  
Body 示例：
```json
{
  "guest_name": "张三",
  "staff_name": "staff_user",
  "contact_phone": "13800000000",
  "contact_email": "guest@example.com",
  "description": "黑色行李箱",
  "quantity": 1,
  "special_notes": "易碎",
  "photo_url": "http://example.com/photo.jpg",
  "storeroom_id": 1
}
```
响应示例：
```json
{
  "message": "create luggage success",
  "luggage_id": 1,
  "retrieval_code": "Z75BDSRH",
  "qrcode_url": "/qr/Z75BDSRH",
  "photo_url": "http://example.com/photo.jpg"
}
```
失败响应示例：
```json
{
  "message": "create luggage failed",
  "error": "storeroom not found"
}
```

### 3.4 GET /api/luggage/by_code
作用：按取件码查询寄存单。  
Query 参数：`code`（必填，取件码）。  
请求示例：
```
GET /api/luggage/by_code?code=Z75BDSRH
```
响应示例：
```json
{
  "message": "query luggage success",
  "item": {
    "id": 1,
    "guest_name": "张三",
    "contact_phone": "13800000000",
    "storeroom_id": 1,
    "retrieval_code": "Z75BDSRH",
    "status": "stored",
    "photo_url": "http://example.com/photo.jpg"
  }
}
```
失败响应示例：
```json
{
  "message": "query luggage failed",
  "error": "code is empty"
}
```

### 3.5 POST /api/luggage/{id}/checkout
作用：取件（id 为取件码，取件人自动使用登录账号）。  
Path 参数：`id`（必填，取件码）。  
请求示例：
```
POST /api/luggage/Z75BDSRH/checkout
```
响应示例：
```json
{
  "message": "checkout success",
  "luggage_id": 1
}
```
失败响应示例：
```json
{
  "message": "checkout failed",
  "error": "luggage is not in stored status"
}
```

### 3.6 GET /api/luggage/{id}/checkout
作用：获取当前酒店有行李在存的客人名单。  
Path 参数：`id`（必填，占位即可）。  
请求示例：
```
GET /api/luggage/any/checkout
```
响应示例：
```json
{
  "message": "get checkout info success",
  "items": ["张三", "李四"]
}
```
失败响应示例：
```json
{
  "message": "missing user info"
}
```

### 3.7 GET /api/luggage/list/by_guest_name
作用：查询某客人正在寄存的行李。  
Query 参数：`guest_name`（必填）。  
请求示例：
```
GET /api/luggage/list/by_guest_name?guest_name=张三
```
响应示例：
```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored"
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "list luggage failed",
  "error": "guest_name is empty"
}
```

### 3.8 GET /api/luggage/storerooms
作用：获取当前酒店所有寄存室（含已存数量与剩余容量）。  
请求：无参数。  
响应示例：
```json
{
  "message": "list storerooms success",
  "items": [
    {
      "id": 1,
      "hotel_id": 1,
      "name": "A区-1号",
      "location": "一楼A区",
      "capacity": 50,
      "is_active": true,
      "stored_count": 12,
      "remaining_capacity": 38
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "list storerooms failed",
  "error": "hotel_id is missing"
}
```

### 3.9 GET /api/luggage/storerooms/{id}/orders
作用：获取某寄存室下的行李订单列表。  
Path 参数：`id`（必填，寄存室ID）。  
Query 参数：`status`（可选，例如 `stored`）。  
请求示例：
```
GET /api/luggage/storerooms/1/orders?status=stored
```
响应示例：
```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored"
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "invalid storeroom id"
}
```

### 3.10 POST /api/luggage/storerooms
作用：创建寄存室。  
Body 示例：
```json
{
  "name": "A区-1号",
  "location": "一楼A区",
  "capacity": 50,
  "is_active": true
}
```
响应示例：
```json
{
  "message": "create storeroom success",
  "item": {
    "id": 1,
    "hotel_id": 1,
    "name": "A区-1号",
    "location": "一楼A区",
    "capacity": 50,
    "is_active": true
  }
}
```
失败响应示例：
```json
{
  "message": "create storeroom failed",
  "error": "invalid request"
}
```

### 3.11 PUT /api/luggage/storerooms/{id}
作用：软删除/停用寄存室。  
Path 参数：`id`（必填，寄存室ID）。  
Body 示例：
```json
{
  "is_active": false
}
```
响应示例：
```json
{
  "message": "update storeroom status success"
}
```
失败响应示例：
```json
{
  "message": "update storeroom status failed",
  "error": "invalid storeroom id"
}
```

### 3.12 GET /api/luggage/logs/stored
作用：获取当前酒店寄存记录。  
请求：无参数。  
响应示例：
```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "status": "stored",
      "stored_at": "2026-01-22T10:00:00+08:00"
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

### 3.13 GET /api/luggage/logs/updated
作用：获取当前酒店寄存信息修改记录（含修改前后快照）。  
请求：无参数。  
响应示例：
```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "hotel_id": 1,
      "luggage_id": 1,
      "updated_by": "staff_user",
      "old_data": "{\"guest_name\":\"张三\"}",
      "new_data": "{\"guest_name\":\"张三\",\"special_notes\":\"易碎\"}",
      "updated_at": "2026-01-22T11:00:00+08:00"
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

### 3.14 GET /api/luggage/logs/retrieved
作用：获取当前酒店取出记录。  
请求：无参数。  
响应示例：
```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieved_by": "staff_user",
      "retrieved_at": "2026-01-22T12:00:00+08:00"
    }
  ]
}
```
失败响应示例：
```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

### 3.15 PUT /api/luggage/{id}
作用：修改寄存信息（非迁移）。  
Path 参数：`id`（必填，寄存单ID）。  
Body 示例（字段可选）：
```json
{
  "guest_name": "张三",
  "contact_phone": "13800000000",
  "description": "黑色行李箱-加锁",
  "special_notes": "易碎",
  "photo_url": "http://example.com/new.jpg"
}
```
响应示例：
```json
{
  "message": "update luggage success"
}
```
失败响应示例：
```json
{
  "message": "update luggage failed",
  "error": "invalid luggage id"
}
```
