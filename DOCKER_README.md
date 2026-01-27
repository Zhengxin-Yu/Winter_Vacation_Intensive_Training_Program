# 🐳 Docker 部署指南

## 快速开始

### 前提条件
- 已安装 Docker 和 Docker Compose
- 确保端口 3306、6379、8080 未被占用

### 一键启动

```bash
# 1. 进入项目目录
cd Winter_Vacation_Intensive_Training_Program

# 2. 启动所有服务
docker-compose up -d

# 3. 查看启动日志
docker-compose logs -f app

# 4. 等待服务就绪（约 30 秒）
# 看到 "Listening and serving HTTP on :8080" 表示启动成功
```

### 访问应用

- API 地址：`http://localhost:8080`
- 健康检查：`http://localhost:8080/ping`
- 登录接口：`http://localhost:8080/api/login`

---

## 服务说明

### 服务列表

| 服务名 | 容器名 | 端口 | 说明 |
|--------|--------|------|------|
| mysql | hotel_luggage_mysql | 3306 | MySQL 8.0 数据库 |
| redis | hotel_luggage_redis | 6379 | Redis 7 缓存 |
| app | hotel_luggage_app | 8080 | Go 后端应用 |

### 默认配置

- **MySQL**
  - 用户：root
  - 密码：123456
  - 数据库：hotel_luggage

- **Redis**
  - 无密码
  - 数据库：0

- **JWT 密钥**
  - 默认：your-secret-key-change-in-production
  - ⚠️ 生产环境务必修改！

---

## 常用操作

### 查看服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 查看 MySQL 日志
docker-compose logs -f mysql

# 查看所有日志
docker-compose logs -f
```

### 重启服务

```bash
# 重启应用
docker-compose restart app

# 重启所有服务
docker-compose restart
```

### 停止服务

```bash
# 停止所有服务
docker-compose stop

# 停止并删除容器（数据保留）
docker-compose down

# 停止并删除容器和数据卷（⚠️ 数据会丢失）
docker-compose down -v
```

### 重新构建

```bash
# 代码修改后，重新构建并启动
docker-compose up -d --build

# 强制重新构建（不使用缓存）
docker-compose build --no-cache app
docker-compose up -d
```

---

## 数据库初始化

### 方法 1：手动导入（推荐）

```bash
# 1. 等待 MySQL 启动完成
docker-compose logs mysql | grep "ready for connections"

# 2. 进入 MySQL 容器
docker exec -it hotel_luggage_mysql mysql -uroot -p123456

# 3. 在 MySQL 中执行
mysql> USE hotel_luggage;
mysql> -- 创建表、插入初始数据等
mysql> exit
```

### 方法 2：使用 SQL 文件自动初始化

```bash
# 1. 创建 init.sql 文件在项目根目录
# 2. 取消 docker-compose.yml 中的注释：
#    - ./init.sql:/docker-entrypoint-initdb.d/init.sql

# 3. 重新启动（只在首次启动时执行）
docker-compose down -v
docker-compose up -d
```

### 方法 3：从备份恢复

```bash
# 1. 复制 SQL 备份到容器
docker cp backup.sql hotel_luggage_mysql:/tmp/

# 2. 导入数据
docker exec -it hotel_luggage_mysql mysql -uroot -p123456 hotel_luggage < /tmp/backup.sql
```

---

## 数据持久化

### 数据卷位置

数据存储在 Docker 数据卷中，即使删除容器也不会丢失：

```bash
# 查看数据卷
docker volume ls | grep winter

# 查看数据卷详细信息
docker volume inspect winter_vacation_intensive_training_program_mysql_data
docker volume inspect winter_vacation_intensive_training_program_redis_data
```

### 文件上传

本地上传的文件存储在：
```
hotel_luggage/uploads/
```

该目录已通过 volume 挂载，文件会持久化到主机。

---

## 环境变量配置

### 修改配置

编辑 `docker-compose.yml` 中的 `environment` 部分：

```yaml
app:
  environment:
    # 修改数据库密码
    DB_DSN: "root:新密码@tcp(mysql:3306)/hotel_luggage?..."
    
    # 修改 JWT 密钥
    JWT_SECRET: "your-production-secret-key"
    
    # 禁用 MinIO（删除或注释相关配置）
    # MINIO_ENDPOINT: ...
```

修改后重启：
```bash
docker-compose up -d
```

---

## 故障排查

### 应用启动失败

```bash
# 1. 查看应用日志
docker-compose logs app

# 常见问题：
# - "database connect failed" → MySQL 还未就绪，等待 10 秒后重试
# - "bind: address already in use" → 端口被占用，修改端口映射
```

### MySQL 连接失败

```bash
# 1. 检查 MySQL 是否就绪
docker-compose logs mysql | grep "ready for connections"

# 2. 测试连接
docker exec -it hotel_luggage_mysql mysql -uroot -p123456

# 3. 检查网络
docker network inspect winter_vacation_intensive_training_program_hotel_network
```

### 端口冲突

如果端口被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
services:
  mysql:
    ports:
      - "3307:3306"  # 改为 3307
  
  app:
    ports:
      - "8081:8080"  # 改为 8081
```

---

## 生产环境部署

### 安全建议

1. **修改默认密码**
   ```yaml
   environment:
     MYSQL_ROOT_PASSWORD: 强密码
     JWT_SECRET: 随机生成的强密钥
   ```

2. **限制端口暴露**
   ```yaml
   mysql:
     ports:
       - "127.0.0.1:3306:3306"  # 只允许本地访问
   ```

3. **使用环境变量文件**
   - 创建 `.env` 文件存储敏感信息
   - 将 `.env` 添加到 `.gitignore`

4. **定期备份**
   ```bash
   # 备份数据库
   docker exec hotel_luggage_mysql mysqldump -uroot -p123456 hotel_luggage > backup_$(date +%Y%m%d).sql
   
   # 备份数据卷
   docker run --rm -v winter_vacation_intensive_training_program_mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql_data_backup.tar.gz /data
   ```

### 性能优化

1. **资源限制**
   ```yaml
   app:
     deploy:
       resources:
         limits:
           cpus: '1'
           memory: 512M
   ```

2. **使用生产模式**
   ```yaml
   app:
     environment:
       GIN_MODE: release  # Gin 生产模式
   ```

---

## 卸载

### 完全清理

```bash
# 1. 停止并删除所有容器
docker-compose down

# 2. 删除数据卷（⚠️ 数据会永久丢失）
docker-compose down -v

# 3. 删除镜像
docker rmi $(docker images | grep hotel_luggage | awk '{print $3}')

# 4. 清理未使用的资源
docker system prune -a
```

---

## 常见问题 FAQ

### Q1: 如何查看容器内部文件？

```bash
docker exec -it hotel_luggage_app sh
ls -la /app
```

### Q2: 如何在容器中执行 Go 命令？

```bash
# 进入应用容器
docker exec -it hotel_luggage_app sh

# 但注意：容器内只有编译好的二进制文件，没有源代码
```

### Q3: 本地开发时也用 Docker 吗？

可以混合使用：
- MySQL 和 Redis 用 Docker：`docker-compose up -d mysql redis`
- Go 应用本地运行：`go run ./cmd/main.go`
- 需要修改数据库连接为 `localhost:3306`

### Q4: 如何更新应用？

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动
docker-compose up -d --build
```

---

## 技术支持

如有问题，请查看：
1. 应用日志：`docker-compose logs app`
2. 容器状态：`docker-compose ps`
3. 网络连接：`docker network inspect hotel_network`

---

**最后更新**：2026年1月
