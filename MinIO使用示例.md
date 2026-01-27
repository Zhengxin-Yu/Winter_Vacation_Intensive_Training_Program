# MinIO 使用示例

## 📝 完整使用流程

### 步骤1：启动 MinIO 服务

```powershell
# Windows PowerShell
C:\minio\minio.exe server C:\minio\data --console-address ":9001"
```

**看到以下日志表示启动成功**：
```
MinIO Object Storage Server
Copyright: 2015-2024 MinIO, Inc.
License: GNU AGPLv3 <https://www.gnu.org/licenses/agpl-3.0.html>
Version: RELEASE.2024-XX-XX

API: http://192.168.1.100:9000  http://127.0.0.1:9000
Console: http://192.168.1.100:9001  http://127.0.0.1:9001

Documentation: https://min.io/docs/minio/linux/index.html
Warning: The standard parity is set to 0. This can lead to data loss.
```

### 步骤2：访问 MinIO 控制台（可选）

1. 浏览器打开：http://localhost:9001
2. 登录：
   - Username: `minioadmin`
   - Password: `minioadmin`
3. 查看是否有 `hotel-luggage` 桶（项目启动时会自动创建）

### 步骤3：启动后端服务

```bash
cd hotel_luggage
go run ./cmd/main.go
```

**看到以下日志表示成功**：
```
✅ 数据库连接成功
✅ Redis初始化成功: 127.0.0.1:6379
✅ MinIO初始化成功: localhost:9000/hotel-luggage
✅ MinIO bucket 'hotel-luggage' 创建成功
[GIN-debug] Listening and serving HTTP on 10.154.101.161:8080
```

### 步骤4：测试上传图片

#### 方法1：使用 curl

```bash
# 1. 先登录获取 token
curl -X POST http://10.154.101.161:8080/api/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"staff_user\",\"password\":\"123456\"}"

# 响应示例
{
  "message": "login success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "staff_user",
    "role": "staff",
    "hotel_id": 1
  }
}

# 2. 上传图片（替换 YOUR_JWT_TOKEN）
curl -X POST http://10.154.101.161:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@C:\Users\你的用户名\Pictures\test.jpg"
```

**成功响应（MinIO）**：
```json
{
  "message": "upload success (MinIO)",
  "url": "http://localhost:9000/hotel-luggage/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg",
  "object_name": "uploads/2026/01/a1b2c3d4e5f6g7h8.jpg",
  "content_type": "image/jpeg",
  "size": 102400,
  "file_name": "a1b2c3d4e5f6g7h8.jpg",
  "max_size_byte": 5242880,
  "storage": "minio"
}
```

#### 方法2：使用 Apifox/Postman

1. **新建请求**
   - Method: POST
   - URL: `http://10.154.101.161:8080/api/upload`

2. **设置 Headers**
   - Key: `Authorization`
   - Value: `Bearer YOUR_JWT_TOKEN`

3. **设置 Body**
   - 选择 `form-data`
   - 添加字段：
     - Key: `file` (类型选择 File)
     - Value: 选择一张图片文件

4. **发送请求**

5. **查看响应**
   - `storage: "minio"` 表示使用 MinIO
   - `storage: "local"` 表示降级到本地存储

### 步骤5：验证图片已上传

#### 方法1：通过 MinIO 控制台

1. 访问 http://localhost:9001
2. 进入 `Buckets` → `hotel-luggage` → `uploads`
3. 查看文件列表，找到刚上传的图片
4. 点击图片可以预览

#### 方法2：通过浏览器直接访问

复制响应中的 `url` 字段，粘贴到浏览器地址栏：
```
http://localhost:9000/hotel-luggage/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg
```

### 步骤6：在行李寄存中使用图片URL

```bash
curl -X POST http://10.154.101.161:8080/api/luggage \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"guest_name\": \"张三\",
    \"contact_phone\": \"13800000000\",
    \"storeroom_id\": 1,
    \"description\": \"黑色行李箱\",
    \"photo_url\": \"http://localhost:9000/hotel-luggage/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg\",
    \"photo_urls\": [
      \"http://localhost:9000/hotel-luggage/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg\",
      \"http://localhost:9000/hotel-luggage/uploads/2026/01/another_photo.jpg\"
    ]
  }"
```

---

## 🔍 常见场景

### 场景1：批量上传多张图片

```bash
# 上传第一张
curl -X POST http://10.154.101.161:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@photo1.jpg"
# 得到 url1

# 上传第二张
curl -X POST http://10.154.101.161:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@photo2.jpg"
# 得到 url2

# 创建寄存单时使用这些URL
curl -X POST http://10.154.101.161:8080/api/luggage \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"guest_name\": \"张三\",
    \"storeroom_id\": 1,
    \"photo_urls\": [\"url1\", \"url2\"]
  }"
```

### 场景2：MinIO 未启动时的降级

```bash
# MinIO 未启动或连接失败时
curl -X POST http://10.154.101.161:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@photo.jpg"

# 响应（自动降级到本地存储）
{
  "message": "upload success (Local)",
  "url": "http://10.154.101.161:8080/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg",
  "relative_url": "/uploads/2026/01/a1b2c3d4e5f6g7h8.jpg",
  "content_type": "image/jpeg",
  "size": 102400,
  "file_name": "a1b2c3d4e5f6g7h8.jpg",
  "max_size_byte": 5242880,
  "storage": "local"  ← 注意这里
}
```

### 场景3：更换为云端 MinIO

```bash
# 修改环境变量
set "MINIO_ENDPOINT=minio.yourdomain.com:9000"
set "MINIO_ACCESS_KEY=your_access_key"
set "MINIO_SECRET_KEY=your_secret_key"
set "MINIO_USE_SSL=true"

# 重启服务
go run ./cmd/main.go

# 上传的图片会存储到云端 MinIO
```

---

## 🎨 前端集成示例

### React + Axios

```javascript
// 1. 上传图片
async function uploadImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await axios.post('http://10.154.101.161:8080/api/upload', formData, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`,
      'Content-Type': 'multipart/form-data'
    }
  });
  
  return response.data.url; // 返回图片URL
}

// 2. 创建寄存单时使用图片URL
async function createLuggage(data) {
  const imageUrls = [];
  
  // 上传所有选中的图片
  for (const file of data.images) {
    const url = await uploadImage(file);
    imageUrls.push(url);
  }
  
  // 创建寄存单
  const response = await axios.post('http://10.154.101.161:8080/api/luggage', {
    guest_name: data.guestName,
    contact_phone: data.phone,
    storeroom_id: data.storeroomId,
    description: data.description,
    photo_urls: imageUrls
  }, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`,
      'Content-Type': 'application/json'
    }
  });
  
  return response.data;
}
```

### Vue 3 + Element Plus

```vue
<template>
  <el-upload
    action="http://10.154.101.161:8080/api/upload"
    :headers="uploadHeaders"
    :on-success="handleSuccess"
    :before-upload="beforeUpload"
    list-type="picture-card"
    multiple
  >
    <el-icon><Plus /></el-icon>
  </el-upload>
</template>

<script setup>
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';

const photoUrls = ref([]);
const token = localStorage.getItem('token');

const uploadHeaders = computed(() => ({
  'Authorization': `Bearer ${token}`
}));

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/');
  const isLt5M = file.size / 1024 / 1024 < 5;
  
  if (!isImage) {
    ElMessage.error('只能上传图片文件！');
    return false;
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB！');
    return false;
  }
  return true;
};

const handleSuccess = (response) => {
  if (response.url) {
    photoUrls.value.push(response.url);
    ElMessage.success('上传成功！');
  }
};
</script>
```

---

## 🛠️ 故障排查

### 问题1：上传后显示 "storage": "local"

**原因**：MinIO 未启动或连接失败

**解决**：
```bash
# 1. 检查 MinIO 是否启动
curl http://localhost:9000/minio/health/live

# 2. 查看后端日志
# 如果看到以下信息，说明 MinIO 未成功连接：
⚠️  MinIO初始化失败: connection refused (将使用本地文件存储)

# 3. 启动 MinIO
C:\minio\minio.exe server C:\minio\data --console-address ":9001"

# 4. 重启后端服务
```

### 问题2：图片无法访问（403 Forbidden）

**原因**：Bucket 策略未设置为公开

**解决**：
```bash
# 方法1：代码已自动设置，重启服务即可
go run ./cmd/main.go

# 方法2：在 MinIO 控制台手动设置
# 访问 http://localhost:9001
# Buckets → hotel-luggage → Access Policy → 选择 "public"
```

### 问题3：上传成功但URL无法访问

**原因**：返回的URL中的host可能不正确

**解决**：
```bash
# 修改 MinIO 配置，使用公网IP或域名
set "MINIO_ENDPOINT=你的公网IP:9000"
# 或
set "MINIO_ENDPOINT=minio.yourdomain.com:9000"

# 重启服务
```

---

## 📊 性能对比

| 项目 | 本地存储 | MinIO |
|------|----------|-------|
| 上传速度 | ⚡ 快（直接写磁盘） | ⚡ 快（网络传输） |
| 访问速度 | 🐢 中等（经过Gin转发） | ⚡ 快（直接访问） |
| 扩展性 | ❌ 受限 | ✅ 易扩展 |
| 高可用 | ❌ 单点故障 | ✅ 支持集群 |
| 适用场景 | 开发测试 | 生产环境 |

---

## ✅ 总结

通过集成 MinIO，你的项目获得了：

1. **专业的对象存储**：图片存储更规范
2. **更好的性能**：图片访问不经过后端转发
3. **易于扩展**：可以轻松扩展到云端
4. **高可用性**：支持分布式部署
5. **优雅降级**：MinIO不可用时自动使用本地存储

这是一个**生产级别**的解决方案！🎉
