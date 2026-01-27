# MinIO 对象存储配置说明

## 📦 MinIO 简介

MinIO 是一个高性能的分布式对象存储服务，兼容 Amazon S3 API。

### 为什么使用 MinIO？

1. **专业的对象存储**：专为存储图片、视频等非结构化数据设计
2. **易于扩展**：可以轻松扩展存储容量
3. **高可用性**：支持分布式部署
4. **S3兼容**：可以无缝迁移到云服务（阿里云OSS、AWS S3等）

---

## 🚀 快速开始

### 1. 安装 MinIO（Windows）

#### 方法1：使用二进制文件（推荐）

```powershell
# 下载 MinIO
Invoke-WebRequest -Uri "https://dl.min.io/server/minio/release/windows-amd64/minio.exe" -OutFile "C:\minio\minio.exe"

# 创建数据目录
mkdir C:\minio\data

# 启动 MinIO 服务器
C:\minio\minio.exe server C:\minio\data --console-address ":9001"
```

#### 方法2：使用 Docker（如果已安装Docker）

```bash
docker run -d \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio \
  -v C:\minio\data:/data \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  quay.io/minio/minio server /data --console-address ":9001"
```

### 2. 访问 MinIO 控制台

启动后，访问：
- **API地址**：http://localhost:9000
- **控制台地址**：http://localhost:9001

**默认账号密码**：
- 用户名：`minioadmin`
- 密码：`minioadmin`

### 3. 配置环境变量

在项目根目录创建 `.env` 文件（或在系统中设置环境变量）：

```bash
# MinIO 配置
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=hotel-luggage
```

### 4. 启动项目

```bash
cd hotel_luggage
go run ./cmd/main.go
```

**启动日志**：
```
✅ 数据库连接成功
✅ Redis初始化成功
✅ MinIO初始化成功: localhost:9000/hotel-luggage
✅ MinIO bucket 'hotel-luggage' 创建成功
[GIN-debug] Listening and serving HTTP on 10.154.101.161:8080
```

---

## 📝 使用说明

### 上传图片

```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@image.jpg"
```

**成功响应（MinIO）**：
```json
{
  "message": "upload success (MinIO)",
  "url": "http://localhost:9000/hotel-luggage/uploads/2026/01/abc123def456.jpg",
  "object_name": "uploads/2026/01/abc123def456.jpg",
  "content_type": "image/jpeg",
  "size": 102400,
  "file_name": "abc123def456.jpg",
  "max_size_byte": 5242880,
  "storage": "minio"
}
```

**降级响应（本地存储）**：
```json
{
  "message": "upload success (Local)",
  "url": "http://localhost:8080/uploads/2026/01/abc123def456.jpg",
  "relative_url": "/uploads/2026/01/abc123def456.jpg",
  "content_type": "image/jpeg",
  "size": 102400,
  "file_name": "abc123def456.jpg",
  "max_size_byte": 5242880,
  "storage": "local"
}
```

---

## 🔧 高级配置

### 1. 生产环境配置

**修改默认密码**（强烈建议）：

```bash
# 设置环境变量
export MINIO_ROOT_USER=your_username
export MINIO_ROOT_PASSWORD=your_strong_password

# 启动 MinIO
minio server /data --console-address ":9001"
```

**在项目中使用**：
```bash
MINIO_ACCESS_KEY=your_username
MINIO_SECRET_KEY=your_strong_password
```

### 2. 启用 HTTPS

```bash
# 生成证书
mkdir -p ~/.minio/certs
# 将 private.key 和 public.crt 放到 ~/.minio/certs/ 目录

# 启动时自动启用 HTTPS
minio server /data --console-address ":9001"

# 项目配置
MINIO_USE_SSL=true
MINIO_ENDPOINT=yourdomain.com:9000
```

### 3. 设置 Bucket 策略

在 MinIO 控制台中：
1. 登录 http://localhost:9001
2. 进入 Buckets → hotel-luggage
3. 点击 "Access Policy"
4. 设置为 "public" 或 "download"

**或通过命令行**（已在代码中自动设置）：
```go
// 代码已自动将 bucket 设置为公开读取
// 见 repositories/minio.go
```

---

## 🛡️ 降级策略

### 自动降级机制

项目实现了**优雅降级**：

1. **MinIO 可用**：上传到 MinIO
2. **MinIO 不可用**：自动降级到本地文件存储
3. **不影响核心功能**：上传功能始终可用

**降级触发条件**：
- MinIO 服务未启动
- MinIO 连接失败
- MinIO 上传失败

**日志提示**：
```
⚠️  MinIO初始化失败: connection refused (将使用本地文件存储)
⚠️  MinIO上传失败，降级到本地存储: timeout
```

---

## 📊 对比：MinIO vs 本地存储

| 特性 | MinIO | 本地存储 |
|------|-------|----------|
| **扩展性** | ✅ 易扩展 | ❌ 受限于磁盘 |
| **高可用** | ✅ 支持分布式 | ❌ 单点故障 |
| **访问速度** | ⚡ 快（专用服务） | 🐢 中等 |
| **成本** | 💰 需要额外服务 | 💵 无额外成本 |
| **迁移云端** | ✅ 兼容S3协议 | ❌ 需要改造 |
| **适用场景** | 生产环境 | 开发/测试环境 |

---

## 🔍 故障排查

### 1. MinIO 启动失败

**问题**：端口被占用
```
Error: listen tcp :9000: bind: address already in use
```

**解决**：
```bash
# 查找占用端口的进程
netstat -ano | findstr :9000

# 杀死进程（记下 PID）
taskkill /PID <进程ID> /F
```

### 2. 上传失败

**问题**：权限不足
```
Access Denied
```

**解决**：
1. 检查 Access Key 和 Secret Key 是否正确
2. 检查 Bucket 策略是否允许上传
3. 查看 MinIO 控制台日志

### 3. 无法访问上传的图片

**问题**：Bucket 未设置为公开
```
403 Forbidden
```

**解决**：
```bash
# 方法1：在控制台设置 Bucket 为 public
# 方法2：代码已自动设置（见 repositories/minio.go）
```

---

## 📚 参考资料

- MinIO 官网：https://min.io/
- MinIO 文档：https://min.io/docs/minio/linux/index.html
- Go SDK 文档：https://min.io/docs/minio/linux/developers/go/minio-go.html

---

## 💡 最佳实践

### 开发环境

```bash
# 使用默认配置即可
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
```

### 生产环境

```bash
# 使用强密码
MINIO_ACCESS_KEY=production_user_12345
MINIO_SECRET_KEY=very_strong_password_here_67890

# 启用 HTTPS
MINIO_USE_SSL=true

# 使用域名
MINIO_ENDPOINT=minio.yourdomain.com:9000

# 定期备份
# 设置存储策略
# 监控存储使用情况
```

---

## ✅ 验证 MinIO 是否正常工作

### 测试脚本

```bash
# 1. 测试 MinIO 是否启动
curl http://localhost:9000/minio/health/live

# 2. 测试上传接口
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@test.jpg"

# 3. 查看返回的 storage 字段
# "storage": "minio"  表示使用 MinIO
# "storage": "local"  表示降级到本地
```

### 在 MinIO 控制台验证

1. 访问 http://localhost:9001
2. 登录
3. 进入 Buckets → hotel-luggage → uploads
4. 查看上传的文件

---

## 🎉 完成！

现在你的项目已经集成了 MinIO 对象存储服务！

- ✅ 专业的图片存储方案
- ✅ 支持自动降级
- ✅ 兼容 S3 协议
- ✅ 易于扩展到云端
