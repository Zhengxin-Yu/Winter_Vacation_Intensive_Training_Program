# 前端对接指南（含 Apifox 测试）

本指南用于说明如何把前端与当前后端接口对接，并提供 Apifox 的完整测试流程。

---

## 1. 启动后端服务

在 Windows CMD 中执行：

```bat
cd /d C:\Users\32660\workspace\Winter_Vacation_Intensive_Training_Program\hotel_luggage
set "DB_DSN=root:你的密码@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd
```

看到 `Listening and serving HTTP on :8080` 说明启动成功。

---

## 2. 前端对接核心概念

- **接口地址（Base URL）**：`http://localhost:8080`
- **请求方式**：GET/POST/PUT/DELETE
- **请求体格式**：JSON
- **响应格式**：JSON

前端通过 `fetch` / `axios` 访问接口即可。

---

## 3. 常用接口速查

### 登录
```
POST /login
```
Body：
```json
{
  "username": "admin",
  "password": "123456"
}
```

### 创建用户
```
POST /users
```
Body：
```json
{
  "username": "staff_user",
  "password": "123456",
  "role": "staff"
}
```

### 行李寄存（必须 guest_name + staff_name）
```
POST /storage
```
Body：
```json
{
  "guest_name": "guest_user",
  "staff_name": "staff_user",
  "contact_phone": "13800000000",
  "description": "黑色行李箱",
  "quantity": 1,
  "storeroom_id": 1
}
```

### 取件（retrieved_by 传用户名）
```
POST /storage/retrieve
```
Body：
```json
{
  "code": "取件码",
  "retrieved_by": "staff_user"
}
```

### 寄存单列表（按用户名）
```
GET /storage/list?username=staff_user
```

### 取件码列表（按用户名）
```
GET /pickup-codes?username=staff_user
```

### 取件历史（按客人姓名/手机号）
```
GET /storage/history/by-guest?guest_name=张三&contact_phone=13800000000
```

---

## 4. 前端调用示例（fetch）

```js
fetch("http://localhost:8080/storage", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    guest_name: "guest_user",
    staff_name: "staff_user",
    contact_phone: "13800000000",
    description: "黑色行李箱",
    quantity: 1,
    storeroom_id: 1
  })
})
  .then(res => res.json())
  .then(console.log);
```

---

## 5. Apifox 测试流程（详细）

### 5.1 新建项目
1. 打开 Apifox
2. 点击左侧 **新建项目**
3. 项目名填 `hotel_luggage`

### 5.2 创建接口
以 “创建行李寄存” 为例：

1. 点击 **新建接口**
2. 名称：创建行李寄存  
3. 方法：`POST`  
4. URL：`http://localhost:8080/storage`

### 5.3 配置请求体
1. 选择 **Body**
2. 选择 **JSON**
3. 填写：
```json
{
  "guest_name": "guest_user",
  "staff_name": "staff_user",
  "contact_phone": "13800000000",
  "description": "黑色行李箱",
  "quantity": 1,
  "storeroom_id": 1
}
```

### 5.4 发送请求
点击右上角 **发送**，成功响应示例：
```json
{
  "message": "create luggage success",
  "luggage_id": 1,
  "retrieval_code": "Z75BDSRH",
  "qrcode_url": "/qr/Z75BDSRH"
}
```

### 5.5 其他接口测试示例

#### 登录
```
POST http://localhost:8080/login
```
Body：
```json
{
  "username": "admin",
  "password": "123456"
}
```

#### 取件
```
POST http://localhost:8080/storage/retrieve
```
    Body：
    ```json
    {
    "code": "取件码",
    "retrieved_by": "staff_user"
    }
    ```

#### 查询寄存单列表
```
GET http://localhost:8080/storage/list?username=staff_user
```

#### 查询取件历史
```
GET http://localhost:8080/storage/history/by-guest?guest_name=张三&contact_phone=13800000000
```

---

## 6. 常见问题

### 6.1 返回 `EOF` 或 `invalid request`
- 未传 JSON Body
- Body 格式错误
- 没设置 `Content-Type: application/json`

### 6.2 返回 404
- 服务未重启
- 接口路径写错
- 服务未启动

### 6.3 登录失败
- 用户不存在
- 密码不匹配（bcrypt 必须用接口创建）
