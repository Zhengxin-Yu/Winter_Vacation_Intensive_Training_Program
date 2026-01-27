# 前端对接指南

本文档面向前端同学，描述接口的用途、请求体/响应体字段含义与示例。以实际后端实现为准。

---

## 0. 基础信息

- **Base URL**：`http://10.154.101.161:8080`
- **数据格式**：
  - JSON 接口：`Content-Type: application/json`
  - 上传接口：`multipart/form-data`
- **认证方式**：
  - 除 `GET /ping`、`POST /api/login` 外，其余 `/api/**` 均需要 Header：

```
Authorization: Bearer <token>
```

- **通用响应字段**：
  - `message`：字符串，表示本次请求结果（成功/失败）
  - `error`：字符串（可选），失败原因

---

## 1. 健康检查

### 1.1 GET `/ping`

**用途**：检查服务是否运行。  
**请求**：无

**响应（200）**：
```json
{ "message": "pong" }
```

---

## 2. 登录

### 2.1 POST `/api/login`

**用途**：登录获取 JWT token。

**请求体（JSON）**：
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `username` | string | 是 | 登录用户名 |
| `password` | string | 是 | 登录密码 |

**响应（200）**：
```json
{
  "message": "login success",
  "user": { "id": 1, "username": "staff_user", "role": "staff", "hotel_id": 1 },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**失败示例（401）**：
```json
{ "message": "login failed", "error": "invalid username or password" }
```

---

## 3. 图片上传（用于 `photo_urls` / `photo_url`）

> 说明：行李表里存的是图片地址（URL），图片本体不进数据库。  
> 推荐把上传接口返回的 `relative_url` 写入 `photo_urls`（数组），便于“一个寄存单多张图”。

### 3.1 POST `/api/upload`（需要登录）

**用途**：上传图片，返回可访问的图片 URL。

**请求（multipart/form-data）**：
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `file` | file | 是 | 图片文件（支持 jpg/png/webp，最大 5MB） |

**响应（200）**：
```json
{
  "message": "upload success",
  "url": "http://10.154.101.161:8080/uploads/2026/01/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "relative_url": "/uploads/2026/01/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "content_type": "image/jpeg",
  "size": 123456,
  "file_name": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "max_size_byte": 5242880
}
```

**失败示例（400）**：
```json
{ "message": "upload failed", "error": "missing file" }
```

### 3.2 图片访问（公开）

上传后的图片可通过静态资源访问：

- `GET /uploads/...`

例如：

- `http://10.154.101.161:8080/uploads/2026/01/xxx.jpg`

---

## 4. 行李寄存单

### 4.1 POST `/api/luggage`（创建寄存单，需要登录）

**用途**：创建寄存单并生成取件码。

**请求体（JSON）**：
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 是 | 客人姓名 |
| `staff_name` | string | 否 | 经办人姓名；不传则后端自动用当前登录账号 |
| `contact_phone` | string | 否 | 联系电话 |
| `contact_email` | string | 否 | 联系邮箱 |
| `description` | string | 否 | 行李描述（单件模式） |
| `quantity` | number | 否 | 数量（默认 1，单件模式） |
| `special_notes` | string | 否 | 备注（单件模式） |
| `photo_urls` | string[] | 否 | 图片地址数组（建议用 `/api/upload` 返回的 `relative_url` 组成数组） |
| `photo_url` | string | 否 | 单图兼容字段（如果只传它，后端会自动转成 `photo_urls=[photo_url]`） |
| `storeroom_id` | number | 否 | 寄存室 ID（单件模式必填） |
| `items` | object[] | 否 | 多件模式（同一单多件可不同寄存室） |

`items` 内每个元素支持字段：`storeroom_id`（必填）、`description`、`quantity`、`special_notes`、`photo_url`、`photo_urls`。

**响应（200）**：
```json
{
  "message": "create luggage success",
  "luggage_id": 1,
  "retrieval_code": "Z75BDSRH",
  "qrcode_url": "/qr/Z75BDSRH",
  "photo_url": "/uploads/2026/01/xxx.jpg",
  "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
}
```

**多件模式示例**：
```json
{
  "guest_name": "张三",
  "staff_name": "staff_user",
  "contact_phone": "13800000000",
  "items": [
    {
      "storeroom_id": 1,
      "description": "行李箱A",
      "quantity": 1,
      "photo_urls": ["/uploads/2026/01/a.jpg"]
    },
    {
      "storeroom_id": 2,
      "description": "行李箱B",
      "quantity": 1,
      "photo_urls": ["/uploads/2026/01/b.jpg", "/uploads/2026/01/c.jpg"]
    }
  ]
}
```

**多件模式响应示例**（同一寄存单共用一个取件码）：
```json
{
  "message": "create luggage success",
  "retrieval_code": "123456",
  "items": [
    {
      "luggage_id": 1,
      "storeroom_id": 1,
      "photo_url": "/uploads/2026/01/a.jpg",
      "photo_urls": ["/uploads/2026/01/a.jpg"]
    },
    {
      "luggage_id": 2,
      "storeroom_id": 2,
      "photo_url": "/uploads/2026/01/b.jpg",
      "photo_urls": ["/uploads/2026/01/b.jpg", "/uploads/2026/01/c.jpg"]
    }
  ]
}
```

**失败示例（400）**：
```json
{ "message": "create luggage failed", "error": "storeroom not found" }
```

### 4.2 GET `/api/luggage/by_code`（按取件码查询，需要登录）

**请求参数（Query）**：
| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `code` | string | 是 | 取件码（`retrieval_code`） |

**响应（200）**：
```json
{
  "message": "query luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "contact_phone": "13800000000",
      "storeroom_id": 1,
      "retrieval_code": "Z75BDSRH",
      "status": "stored",
      "photo_url": "/uploads/2026/01/xxx.jpg",
      "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
    }
  ]
}
```

**失败示例（400）**：
```json
{ "message": "query luggage failed", "error": "code is empty" }
```

### 4.3 POST `/api/luggage/{code}/checkout`（取件，需要登录）

> Path 参数名在路由里叫 `:id`，实际传取件码即可：`/api/luggage/Z75BDSRH/checkout`

**响应（200）**：
```json
{
  "message": "checkout success",
  "retrieval_code": "Z75BDSRH",
  "retrieved_count": 2,
  "luggage_ids": [1, 2],
  "luggage_id": null
}
```

**失败示例（400）**：
```json
{ "message": "checkout failed", "error": "luggage is not in stored status" }
```

### 4.4 GET `/api/luggage/{any}/checkout`（获取当前酒店“在存”客人名单，需要登录）

> Path 参数占位即可，例如：`/api/luggage/any/checkout`

**响应（200）**：
```json
{ "message": "get checkout info success", "items": ["张三", "李四"] }
```

**失败示例（401）**：
```json
{ "message": "missing user info" }
```

### 4.5 GET `/api/luggage/list/by_guest_name`（按客人姓名查在存行李，需要登录）

**Query 参数**：
| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 是 | 客人姓名 |

**响应（200）**：
```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored",
      "photo_url": "/uploads/2026/01/xxx.jpg",
      "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
    }
  ]
}
```

**失败示例（400）**：
```json
{ "message": "list luggage failed", "error": "guest_name is empty" }
```

### 4.6 PUT `/api/luggage/{id}`（修改寄存信息，需要登录）

**用途**：修改寄存信息（包含 `photo_url` / `photo_urls`）。

**请求体（JSON，字段可选）**：
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 否 | 客人姓名 |
| `contact_phone` | string | 否 | 联系电话 |
| `description` | string | 否 | 行李描述 |
| `special_notes` | string | 否 | 备注 |
| `photo_urls` | string[] | 否 | 图片地址数组（推荐） |
| `photo_url` | string | 否 | 单图兼容字段 |

**响应（200）**：
```json
{ "message": "update luggage success" }
```

**失败示例（400）**：
```json
{ "message": "update luggage failed", "error": "invalid luggage id" }
```

---

## 5. 寄存室

### 5.1 GET `/api/luggage/storerooms`（需要登录）

**响应（200）**：返回寄存室列表（后端计算并返回容量信息）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 寄存室 ID |
| `hotel_id` | number | 酒店 ID |
| `name` | string | 寄存室名称 |
| `location` | string | 位置 |
| `capacity` | number | 容量（最大可存数量） |
| `is_active` | boolean | 是否启用 |
| `stored_count` | number | 已存数量（后端计算字段） |
| `remaining_capacity` | number | 剩余容量（后端计算字段） |

**失败示例（400）**：
```json
{ "message": "list storerooms failed", "error": "hotel_id is missing" }
```

### 5.2 POST `/api/luggage/storerooms`（需要登录）

**请求体（JSON）**：
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `name` | string | 是 | 名称 |
| `location` | string | 否 | 位置 |
| `capacity` | number | 是 | 容量 |
| `is_active` | boolean | 否 | 是否启用（不传时按后端默认值处理） |

**响应（200）**：
```json
{ "message": "create storeroom success", "item": { "id": 1, "hotel_id": 1, "name": "A区-1号", "location": "一楼A区", "capacity": 50, "is_active": true } }
```

**失败示例（400）**：
```json
{ "message": "create storeroom failed", "error": "invalid request" }
```

### 5.3 PUT `/api/luggage/storerooms/{id}`（需要登录）

**用途**：启用/停用寄存室（软停用）。

**请求体（JSON）**：
```json
{ "is_active": false }
```

**响应（200）**：
```json
{ "message": "update storeroom status success" }
```

**失败示例（400）**：
```json
{ "message": "update storeroom status failed", "error": "invalid storeroom id" }
```

### 5.4 GET `/api/luggage/storerooms/{id}/orders`（需要登录）

**Query 参数**：
| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `status` | string | 否 | 过滤状态（如 `stored`） |

**响应（200）**：
```json
{
  "message": "list luggage success",
  "items": [
    { "id": 1, "guest_name": "张三", "retrieval_code": "Z75BDSRH", "status": "stored" }
  ]
}
```

**失败示例（400）**：
```json
{ "message": "invalid storeroom id" }
```

---

## 6. 日志

### 6.1 GET `/api/luggage/logs/stored`（需要登录）

**响应（200）**：寄存记录列表（结构由后端模型返回，常用字段如下）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `guest_name` | string | 客人姓名 |
| `status` | string | 一般为 `stored` |
| `stored_at` | string | 时间（ISO8601） |

**失败示例（400）**：
```json
{ "message": "list logs failed", "error": "hotel_id is missing" }
```

### 6.2 GET `/api/luggage/logs/updated`（需要登录）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `updated_by` | string | 修改人 |
| `old_data` | string | 修改前快照（JSON 字符串） |
| `new_data` | string | 修改后快照（JSON 字符串） |
| `updated_at` | string | 修改时间 |

**失败示例（400）**：
```json
{ "message": "list logs failed", "error": "hotel_id is missing" }
```

### 6.3 GET `/api/luggage/logs/retrieved`（需要登录）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `guest_name` | string | 客人姓名 |
| `retrieved_by` | string | 取件人（登录账号） |
| `retrieved_at` | string | 取件时间 |

**失败示例（400）**：
```json
{ "message": "list logs failed", "error": "hotel_id is missing" }
```

---

## 7. 前端最小流程（建议照这个跑通）

1) `POST /api/login` 获取 `token`  
2) `POST /api/upload` 上传图片（可多次），收集 `relative_url[]`  
3) `POST /api/luggage` 创建寄存单，把 `photo_urls = relative_url[]`（或单图用 `photo_url`）  
4) `GET /api/luggage/by_code?code=...` 查询并展示图片（`<img src="BaseURL + photo_url">`）
