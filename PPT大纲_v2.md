# 酒店行李寄存系统答辩PPT大纲（基于开发实践）

---

## 第1页：封面
- **标题**：酒店行李寄存管理系统
- **副标题**：基于 Go + Gin + MySQL 的后端开发实践
- **答辩人**：俞政鑫
- **日期**：2026年1月

---

## 第2页：目录
1. 项目开发流程
2. 项目架构设计
3. 各模块功能详解
4. 核心功能实现（代码讲解）
5. 数据库设计理念
6. 技术栈应用
7. API接口设计
8. JWT鉴权机制
9. 开发收获与成长
10. 问题与经验教训
11. 项目反思与不足

---

## 第3页：项目开发流程
### 完整开发流程
```
需求分析 → 数据库设计 → 接口设计 → 后端开发 → 前后端联调 → 测试部署
```

### 各阶段详解
1. **需求分析阶段**
   - 与前端讨论业务场景
   - 确定功能需求清单
   - 明确系统边界

2. **数据库设计阶段**
   - 设计ER图
   - 创建表结构
   - 建立索引和外键

3. **接口设计阶段**
   - 定义RESTful API规范
   - 统一请求/响应格式
   - 编写API文档

4. **后端开发阶段**
   - 搭建项目框架
   - 实现业务逻辑
   - 单元测试

5. **联调测试阶段**
   - 前后端对接
   - 接口调试
   - Bug修复

6. **部署上线**
   - 服务器配置
   - 环境变量设置
   - 持续优化

### 实际开发中的迭代
- ❌ **理想**：一次性完成所有设计
- ✅ **现实**：不断迭代优化
  - 多次修改数据库表结构
  - 多次调整接口名称和参数
  - 多次重构代码

---

## 第4页：项目架构设计
### 分层架构
```
hotel_luggage/
├── cmd/                    # 程序入口
│   ├── main.go            # 主程序
│   └── create_user/       # 用户创建工具（在前端未涉及创建用户接口，管理者可以用命令行工具创建用户）
├── configs/               # 配置管理（前置URL）
├── internal/              # 内部核心代码
│   ├── models/           # 数据模型
│   ├── handlers/         # HTTP处理层
│   ├── services/         # 业务逻辑层
│   ├── repositories/     # 数据访问层
│   └── middleware/       # 中间件
├── router/               # 路由配置
└── utils/                # 工具函数（例如code.go 即一个可自动生成6位数字取件码的函数）
```


## 第5页：各模块功能详解（1/3）
### 📁 cmd/ - 程序入口
```go
// cmd/main.go - 主程序
func main() {
    repositories.InitDB()      // 初始化数据库
    repositories.InitRedis()   // 初始化Redis
    r := router.SetupRouter()  // 设置路由
    r.Run(":8080")            // 启动服务
}

// cmd/create_user/main.go - 用户创建工具
// 为什么需要这个工具？
// - 系统不提供注册接口（安全考虑）
// - 由管理员通过命令行创建员工账号
// - 符合企业内部管理流程
```

### 📁 configs/ - 配置管理
```go
// 配置文件的作用：
// 1. 统一管理配置
// 2. 环境变量读取
// 3. 默认配置fallback
// 4. 便于部署到不同环境

type DBConfig struct {
    DSN string  // 数据库连接字符串
}
```

### 📁 models/ - 数据模型
```go
// 为什么需要独立的Model层？
// 1. 对应数据库表结构
// 2. GORM的映射规则
// 3. JSON序列化/反序列化
// 4. 业务实体的定义

type LuggageItem struct {
    ID            int64
    GuestName     string
    RetrievalCode string
    Status        string
    // ...
}
```

---

## 第6页：各模块功能详解（2/3）
### 📁 handlers/ - HTTP处理层
**职责**：
- 接收HTTP请求
- 参数校验（Gin binding）
- 调用Service层
- 封装HTTP响应

**示例**：
```go
func CreateLuggage(c *gin.Context) {
    // 1. 解析请求体
    var req CreateLuggageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // 2. 调用Service层
    result, err := services.CreateLuggage(req)
    
    // 3. 返回响应
    c.JSON(200, gin.H{"data": result})
}
```

### 📁 services/ - 业务逻辑层
**职责**：
- 实现核心业务逻辑
- 数据校验
- 事务处理
- 调用Repository层

**示例**：
```go
func CreateLuggage(req Request) (Item, error) {
    // 1. 业务校验
    if !isValid(req) {
        return nil, errors.New("invalid")
    }
    
    // 2. 生成取件码
    code := utils.GenerateCode()
    
    // 3. 调用Repository保存
    return repositories.SaveLuggage(item)
}
```

---

## 第7页：各模块功能详解（3/3）
### 📁 repositories/ - 数据访问层
**职责**：
- 封装数据库操作
- GORM查询
- Redis缓存
- 事务管理

**示例**：
```go
func SaveLuggage(item *models.LuggageItem) error {
    return DB.Create(item).Error
}

func FindLuggageByCode(code string) ([]models.LuggageItem, error) {
    var items []models.LuggageItem
    err := DB.Where("retrieval_code = ? AND status = ?", 
                    code, "stored").Find(&items).Error
    return items, err
}
```

### 📁 middleware/ - 中间件
**职责**：
- JWT认证
- CORS处理
- 日志记录
- 权限控制

### 📁 router/ - 路由配置
**职责**：
- 定义所有API路径
- 绑定Handler
- 应用中间件
- 路由分组

### 📁 utils/ - 工具函数
**职责**：
- JWT生成/解析
- 取件码生成
- 密码加密
- 通用工具

---

## 第8页：核心功能实现 - 行李寄存（1/4）
### 业务流程
```
前端请求 → Handler解析 → Service校验 → 生成取件码 
→ Repository保存 → 返回结果
```

### 代码实现 - Handler层
```go
// handlers/luggage_handler.go
func CreateLuggage(c *gin.Context) {
    // 1. 获取当前登录用户
    username, _ := c.Get("username")
    
    // 2. 解析请求体
    var req CreateLuggageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    // 3. 区分单件/多件模式
    if len(req.Items) > 0 {
        // 多件模式：共用取件码
        sharedCode := generateUniqueCode()
        for _, item := range req.Items {
            item.RetrievalCode = sharedCode
            // 调用Service保存
        }
    } else {
        // 单件模式：独立取件码
        // 调用Service保存
    }
}
```

**关键点**：
1. 从JWT中提取用户信息
2. 参数校验（Gin binding）
3. 区分单件/多件逻辑

---

## 第9页：核心功能实现 - 行李寄存（2/4）
### Service层 - 业务逻辑
```go
// services/luggage_service.go
func CreateLuggage(req CreateLuggageRequest) (*models.LuggageItem, error) {
    // 1. 校验寄存室
    room, err := repositories.GetStoreroomByID(req.StoreroomID)
    if err != nil {
        return nil, errors.New("storeroom not found")
    }
    if !room.IsActive {
        return nil, errors.New("storeroom is inactive")
    }
    
    // 2. 校验容量
    stored, _ := repositories.CountStoredByStoreroom(req.StoreroomID)
    if stored >= room.Capacity {
        return nil, errors.New("storeroom is full")
    }
    
    // 3. 生成取件码（如果没有传入）
    if req.RetrievalCode == "" {
        req.RetrievalCode = generateUniqueCode()
    }
    
    // 4. 保存到数据库
    item := &models.LuggageItem{
        GuestName:     req.GuestName,
        RetrievalCode: req.RetrievalCode,
        StoreroomID:   req.StoreroomID,
        Status:        "stored",
        // ...
    }
    return repositories.SaveLuggage(item)
}
```

---

## 第10页：核心功能实现 - 行李寄存（3/4）
### 取件码生成逻辑
```go
// utils/code.go
func GenerateCode() string {
    // 生成6位随机数字
    rand.Seed(time.Now().UnixNano())
    code := fmt.Sprintf("%06d", rand.Intn(1000000))
    return code
}

// 为什么是6位数字？
// 1. 容易记忆和输入
// 2. 100万种组合，对单个酒店足够
// 3. 取件后归档，取件码可重复使用
```

### 唯一性检查
```go
// services/luggage_service.go
func generateUniqueCode() string {
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        code := utils.GenerateCode()
        
        // 检查是否已存在
        existing, _ := repositories.FindLuggageByCode(code)
        if len(existing) == 0 {
            return code  // 唯一，可以使用
        }
    }
    // 重试5次仍失败，返回错误
    return ""
}
```

**思考**：为什么重试5次？
- 6位数字重复概率很低
- 平衡性能和可靠性
- 如果5次都冲突，说明系统负载很高

---

## 第11页：核心功能实现 - 行李寄存（4/4）
### Repository层 - 数据库操作
```go
// repositories/luggage_repository.go
func SaveLuggage(item *models.LuggageItem) error {
    // GORM的Create方法
    // 1. 自动填充时间戳（created_at, updated_at）
    // 2. 返回自增ID
    // 3. 参数化查询，防止SQL注入
    return DB.Create(item).Error
}

func FindLuggageByCode(code string) ([]models.LuggageItem, error) {
    var items []models.LuggageItem
    
    // 1. Where条件：取件码匹配 AND 状态为stored
    // 2. Find返回多条记录（支持多件共用取件码）
    err := DB.Where("retrieval_code = ? AND status = ?", 
                    code, "stored").
            Find(&items).Error
    
    return items, err
}
```

### 为什么用GORM？
1. **简化SQL操作**：不用手写SQL
2. **类型安全**：编译时检查
3. **自动映射**：struct ↔ 数据库表
4. **防SQL注入**：自动参数化
5. **链式调用**：代码更优雅

---

## 第12页：数据库设计理念（1/2）
### 核心表结构
| 表名 | 说明 | 设计考虑 |
|------|------|----------|
| **users** | 用户表 | 密码hash存储，角色字段 |
| **hotels** | 酒店表 | 多租户基础 |
| **luggage_storerooms** | 寄存室表 | 容量管理，软删除 |
| **luggage_items** | 行李寄存表 | 核心业务表 |
| **luggage_history** | 取件历史 | 数据归档 |
| **行李寄存信息修改表** | 修改日志 | 操作追溯 |

### 设计原则

#### 1. **多租户隔离**
```sql
-- 所有表都有hotel_id字段
CREATE TABLE luggage_items (
    id BIGINT PRIMARY KEY,
    hotel_id BIGINT NOT NULL,  -- 关键字段
    guest_name VARCHAR(100),
    -- ...
    KEY idx_hotel_id (hotel_id)
);
```
- 确保不同酒店数据隔离
- 所有查询自动带hotel_id过滤
- 防止数据泄露

#### 2. **索引优化**
```sql
-- 高频查询字段建索引
KEY idx_retrieval_code (retrieval_code),
KEY idx_storeroom_id (storeroom_id),
KEY idx_status (status)
```
- 加速查询
- 减少全表扫描

---

## 第13页：数据库设计理念（2/2）
### 3. **软删除设计**
```sql
-- 寄存室表
CREATE TABLE luggage_storerooms (
    id BIGINT PRIMARY KEY,
    is_active TINYINT(1) DEFAULT 1,  -- 1启用，0停用
    -- ...
);
```
**为什么用软删除？**
- 历史数据关联完整
- 可以恢复误删除
- 数据追溯和审计

### 4. **数据归档**
```
寄存时：luggage_items（存储中）
取件后：luggage_items 删除 → luggage_history 归档
```
**优势**：
- luggage_items表保持轻量
- history表保留完整历史
- 查询性能更好

### 5. **JSON字段存储**
```sql
-- photo_urls存储为TEXT类型
photo_urls TEXT NULL
```
```go
// Go中定义为[]string
type LuggageItem struct {
    PhotoURLs []string `gorm:"column:photo_urls;type:text"`
}

// GORM Hooks实现自动转换
func (item *LuggageItem) BeforeSave(tx *gorm.DB) error {
    // []string → JSON string
}
func (item *LuggageItem) AfterFind(tx *gorm.DB) error {
    // JSON string → []string
}
```

---

## 第14页：技术栈应用（1/3）
### 1. Go语言
**为什么选择Go？**
- 高性能，编译型语言
- 天然支持并发（goroutine）
- 简洁的语法
- 丰富的标准库

**在项目中的应用**：
```go
// 并发处理（虽然当前项目未用到，但Go天生支持）
go func() {
    // 异步任务
}()

// 错误处理
if err != nil {
    return nil, err
}

// 结构体和方法
type Service struct {}
func (s *Service) DoSomething() error {}
```

### 2. Gin框架
**为什么选择Gin？**
- 轻量级，性能高
- 路由简洁
- 中间件机制完善
- 参数绑定方便

**在项目中的应用**：
```go
// 1. 路由定义
r := gin.Default()
r.POST("/api/login", handlers.Login)

// 2. 中间件应用
auth := r.Group("/api")
auth.Use(middleware.JWTAuth())  // JWT认证

// 3. 参数绑定
var req Request
c.ShouldBindJSON(&req)  // 自动解析JSON

// 4. 响应封装
c.JSON(200, gin.H{"data": result})
```

---

## 第15页：技术栈应用（2/3）
### 3. GORM
**为什么选择GORM？**
- Go最流行的ORM
- 功能完善
- 链式调用优雅
- 自动迁移

**在项目中的应用**：
```go
// 1. 模型定义
type LuggageItem struct {
    ID        int64     `gorm:"primaryKey"`
    GuestName string    `gorm:"column:guest_name;size:100"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

// 2. CRUD操作
DB.Create(&item)                          // 创建
DB.Where("id = ?", id).First(&item)      // 查询
DB.Model(&item).Update("status", "done") // 更新
DB.Delete(&item)                         // 删除

// 3. 复杂查询
DB.Where("hotel_id = ?", hotelID).
   Where("status = ?", "stored").
   Order("created_at DESC").
   Find(&items)

// 4. 关联查询
DB.Preload("Storeroom").Find(&items)  // 预加载关联
```

### 4. MySQL
**为什么选择MySQL？**
- 成熟稳定
- 社区活跃
- 完善的事务支持
- 丰富的生态

**在项目中的特性应用**：
- 事务（ACID）
- 外键约束
- 索引优化
- JSON字段

---

## 第16页：技术栈应用（3/3）
### 5. Redis
**注意**：项目实际上已经实现了Redis缓存！

**在项目中的应用**：
```go
// repositories/redis.go
var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    // 测试连接，失败则降级
}

// 缓存查询结果
func GetLuggageByCodeCache(code string) ([]models.LuggageItem, bool, error) {
    val, err := RedisClient.Get(ctx, "luggage:"+code).Result()
    if err == redis.Nil {
        return nil, false, nil  // 缓存未命中
    }
    // 反序列化返回
}

// 设置缓存（1分钟TTL）
func SetLuggageByCodeCache(code string, items []models.LuggageItem) {
    json, _ := json.Marshal(items)
    RedisClient.Set(ctx, "luggage:"+code, json, time.Minute)
}
```

**使用场景**：
- 按取件码查询行李（高频操作）
- 取件/修改时清除缓存
- 降级方案：Redis不可用时直接查MySQL

### 6. bcrypt
**密码加密**：
```go
// 创建用户时
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password), bcrypt.DefaultCost)

// 登录验证时
err := bcrypt.CompareHashAndPassword(
    []byte(user.PasswordHash), []byte(password))
```

---

## 第17页：API接口设计理念（1/2）
### RESTful设计原则
```
资源       HTTP方法    路径                      说明
行李寄存    POST       /api/luggage             创建寄存单
行李查询    GET        /api/luggage/by_code     按取件码查询
行李取件    POST       /api/luggage/:code/checkout  取件
行李修改    PUT        /api/luggage/:id         修改信息
寄存室列表  GET        /api/luggage/storerooms  查询列表
寄存室创建  POST       /api/luggage/storerooms  创建
```

### 为什么这样设计？

#### 1. **统一的URL规范**
```
/api/资源名称
/api/资源名称/:id
/api/资源名称/:id/动作
```
- 清晰的资源层次
- 符合RESTful规范
- 易于理解和维护

#### 2. **统一的请求格式**
```json
// 所有POST/PUT请求都用JSON
{
  "guest_name": "张三",
  "contact_phone": "13800000000",
  "storeroom_id": 1
}
```
- Content-Type: application/json
- 结构化数据
- 易于前端处理

---

## 第18页：API接口设计理念（2/2）
### 3. **统一的响应格式**
```json
// 成功响应
{
  "message": "success",
  "data": { /* 业务数据 */ }
}

// 错误响应
{
  "message": "error description",
  "error": "detailed error info"
}
```

**为什么要统一格式？**
- 前端可以统一处理
- 降低沟通成本
- 减少Bug

### 4. **路由分组**
```go
// router/router.go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // 公开接口
    r.POST("/api/login", handlers.Login)
    
    // 需要认证的接口
    auth := r.Group("/api")
    auth.Use(middleware.JWTAuth())  // JWT中间件
    {
        auth.POST("/luggage", handlers.CreateLuggage)
        auth.GET("/luggage/by_code", handlers.QueryByCode)
        // ...
    }
    
    return r
}
```

**优势**：
- 中间件复用
- 代码组织清晰
- 权限管理方便

---

## 第19页：JWT鉴权机制（1/2）
### JWT工作流程
```
1. 用户登录 → 后端验证 → 生成JWT token
2. 前端保存token → 后续请求携带token
3. 后端验证token → 提取用户信息 → 处理请求
```

### 代码实现

#### 1. **生成Token**
```go
// utils/jwt.go
type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    HotelID  int64  `json:"hotel_id"`
    jwt.StandardClaims
}

func GenerateToken(username, role string, hotelID int64) (string, error) {
    claims := Claims{
        Username: username,
        Role:     role,
        HotelID:  hotelID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),  // 24小时有效
            Issuer:    "hotel_luggage",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}
```

#### 2. **登录接口**
```go
// handlers/auth_handler.go
func Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    // 验证用户名密码
    user, err := repositories.GetUserByUsername(req.Username)
    err = bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash), []byte(req.Password))
    
    // 生成token
    token, _ := utils.GenerateToken(user.Username, user.Role, user.HotelID)
    
    c.JSON(200, gin.H{
        "token": token,
        "user":  user,
    })
}
```

---

## 第20页：JWT鉴权机制（2/2）
### 3. **验证Token中间件**
```go
// middleware/jwt.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Header获取token
        auth := c.GetHeader("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            c.JSON(401, gin.H{"message": "unauthorized"})
            c.Abort()
            return
        }
        
        // 2. 解析token
        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        claims, err := utils.ParseToken(tokenStr)
        if err != nil {
            c.JSON(401, gin.H{"message": "invalid token"})
            c.Abort()
            return
        }
        
        // 3. 将用户信息存入Context
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Set("hotel_id", claims.HotelID)
        
        c.Next()  // 继续处理请求
    }
}
```

### 4. **在Handler中使用**
```go
func CreateLuggage(c *gin.Context) {
    // 从Context获取用户信息
    username, _ := c.Get("username")
    hotelID, _ := c.Get("hotel_id")
    
    // 自动按hotel_id过滤数据
    // ...
}
```

### JWT的优势
- ✅ **无状态**：不需要服务端存储session
- ✅ **跨域友好**：适合前后端分离
- ✅ **可扩展**：可以携带自定义信息
- ✅ **性能好**：不需要查询session

---

## 第21页：开发收获与成长（1/4）
### 1. 利用AI辅助编程

#### 使用场景
- **代码生成**：快速生成模板代码
- **Bug调试**：分析错误原因
- **代码优化**：改进代码质量
- **学习新知识**：理解技术原理

#### 实际案例
```
我的问题：如何在Go中实现bcrypt密码加密？

AI帮助：
1. 提供完整代码示例
2. 解释bcrypt原理
3. 说明成本因子的选择
4. 给出最佳实践建议
```

#### 收获
- ✅ **提高开发效率**：减少查资料时间
- ✅ **快速学习**：理解不熟悉的技术
- ✅ **规范代码**：学习最佳实践
- ⚠️ **注意盲目依赖**：要理解代码，不能直接复制

---

## 第22页：开发收获与成长（2/4）
### 2. 前后端分离协作

#### 协作流程
```
需求讨论 → 接口文档编写 → 并行开发 → 联调测试
```

#### 编写API文档的重要性
**我的实践**：编写了详细的 `Frontend_Integration_Guide.md`

**文档内容**：
```markdown
## POST /api/luggage - 创建寄存单

### 请求参数
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| guest_name | string | 是 | 客人姓名 |
| contact_phone | string | 是 | 联系电话 |

### 请求示例
{
  "guest_name": "张三",
  "contact_phone": "13800000000"
}

### 响应示例
{
  "message": "success",
  "data": { ... }
}

### 错误码
- 400: 参数错误
- 401: 未认证
```

#### 收获
- ✅ **减少沟通成本**：文档写清楚，避免反复问
- ✅ **提高开发效率**：前后端并行开发
- ✅ **降低Bug率**：需求明确，减少误解

---

## 第23页：开发收获与成长（3/4）
### 3. 规范的开发流程

#### 统一接口命名规范
```
资源操作      路径                 HTTP方法
创建         /api/resource        POST
查询列表     /api/resource        GET
查询详情     /api/resource/:id    GET
更新         /api/resource/:id    PUT
删除         /api/resource/:id    DELETE
```

#### 统一请求体/响应体格式
```json
// 统一的响应结构
{
  "message": "操作结果描述",
  "data": {
    // 业务数据
  }
}

// 统一的错误响应
{
  "message": "错误描述",
  "error": "详细错误信息"
}
```

#### 为什么要统一？
1. **前端易于处理**
   ```javascript
   // 统一的响应处理
   axios.post('/api/luggage', data)
     .then(res => {
       if (res.data.message === 'success') {
         // 成功处理
       }
     })
   ```

2. **降低沟通成本**
   - 不用每个接口都问一遍格式
   - 新接口直接参考现有规范

3. **便于维护**
   - 统一的错误处理
   - 统一的日志记录

---

## 第24页：开发收获与成长（4/4）
### 4. 团队开发中的沟通

#### 有效沟通的实践

**1. 需求确认**
```
❌ 不好的沟通：
前端："我要一个查询接口"
后端："好的"（开始写代码）

✅ 好的沟通：
前端："我需要按取件码查询行李，返回所有字段"
后端："确认一下，需要返回图片URL吗？需要寄存室信息吗？"
前端："需要图片，不需要寄存室详情"
后端："好的，我写一下接口文档，你看看是否符合需求"
```

**2. 问题反馈**
```
❌ 不好的反馈：
"接口报错了"

✅ 好的反馈：
"POST /api/luggage 接口报400错误，
请求参数：{...}
响应内容：{error: 'invalid request'}
是storeroom_id必填吗？"
```

**3. 变更通知**
```
✅ 主动通知：
"我修改了/api/luggage接口，增加了photo_urls字段，
返回数组类型，已更新文档，请同步一下"
```


## 第25页：开发中的问题与教训（1/3）
### 问题1：图片上传功能实现

#### 遇到的困难
```
❓ 问题：
1. 如何接收multipart/form-data？
2. 如何保存文件？
3. 如何生成访问URL？
4. 如何防止文件名冲突？
```

#### 解决过程
**1. 请教学长**
- 学习了Gin的文件上传API
- 了解了文件存储的最佳实践
- 学会了安全的文件名生成

**2. 实现代码**
```go
func UploadImage(c *gin.Context) {
    // 1. 接收文件
    file, _ := c.FormFile("file")
    
    // 2. 验证文件类型和大小
    if file.Size > 10*1024*1024 {  // 限制10MB
        c.JSON(400, gin.H{"error": "file too large"})
        return
    }
    
    // 3. 生成安全的文件名
    ext := filepath.Ext(file.Filename)
    filename := generateRandomFilename() + ext
    
    // 4. 保存文件
    savePath := "uploads/" + time.Now().Format("2006/01") + "/" + filename
    c.SaveUploadedFile(file, savePath)
    
    // 5. 返回访问URL
    url := "/uploads/" + savePath
    c.JSON(200, gin.H{"url": url})
}
```

#### 经验教训
- ✅ **遇到不会的，及时请教**
- ✅ **学习别人的最佳实践**
- ✅ **注意安全问题**（文件类型、大小、文件名）

---

## 第26页：开发中的问题与教训（2/3）
### 问题2：频繁修改数据库和接口

#### 问题描述
**开发初期的混乱**：
```
第1天：创建luggage_items表
第2天：发现缺少photo_url字段，ALTER TABLE
第3天：需要支持多图，改为photo_urls TEXT
第4天：增加order_id字段
第5天：发现order_id不必要，删除字段
第6天：接口名称从/storage改为/luggage
第7天：响应格式又改了...
```

**导致的问题**：
- ❌ 前端频繁修改代码
- ❌ 测试数据需要重新录入
- ❌ 文档需要反复更新
- ❌ 大量时间浪费在修改上

#### 根本原因
```
没有先设计好 → 边开发边改 → 反复返工
```

---

## 第27页：开发中的问题与教训（3/3）
### 正确的开发流程应该是：

```
1. 需求分析（充分讨论）
   ↓
2. 数据库设计（确定表结构）
   ↓
3. 接口设计（编写API文档）
   ↓
4. 评审（前后端一起review）
   ↓
5. 开发（按照设计实现）
   ↓
6. 测试联调
```

#### 经验教训

**🎯 教训1：先设计，后开发**
```
✅ 花2小时设计好，避免20小时返工
✅ 接口文档要先写，不要边开发边写
✅ 数据库设计要考虑扩展性
```

**🎯 教训2：充分沟通**
```
✅ 需求不清楚，不要开始写代码
✅ 设计方案要和前端确认
✅ 变更要及时通知
```

**🎯 教训3：版本管理**
```
✅ 使用Git管理代码
✅ 重要修改前先commit
✅ 数据库变更记录在migration文档中
```

---

## 第28页：项目反思与不足（1/2）
### 已实现但可以优化的

#### 1. Redis缓存
**现状**：
- ✅ 已实现按取件码查询的缓存
- ✅ 已实现降级方案

**可以改进**：
```go
// 当前：只缓存了按取件码查询
// 可以扩展：
// 1. 缓存寄存室信息（很少变化）
// 2. 缓存用户信息（每次请求都要查）
// 3. 使用缓存预热（系统启动时加载热点数据）
```

#### 2. 错误处理
**可以改进**：
```go
// 当前：简单返回错误信息
c.JSON(400, gin.H{"error": err.Error()})

// 可以改进：统一的错误码系统
const (
    ErrStoreroomNotFound = 10001
    ErrStoreroomFull    = 10002
    ErrInvalidCode      = 10003
)

c.JSON(400, gin.H{
    "code": ErrStoreroomFull,
    "message": "寄存室已满",
})
```

---

## 第29页：项目反思与不足（2/2）
### 尝试但未能完全实现的功能

#### 1. MinIO 对象存储集成
**设计初衷**：
- 将上传的图片存储到专业的对象存储服务（MinIO）
- 提升文件管理能力和访问性能
- 学习云存储服务的集成

**实现情况**：
```go
// ✅ 已完成的部分
// 1. MinIO 客户端初始化
func InitMinIO() {
    client, err := minio.New(config.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: true,  // 使用HTTPS
    })
    // 设置超时控制、降级方案等
}

// 2. 上传逻辑实现（带降级）
func Upload(c *gin.Context) {
    // 优先尝试MinIO上传
    if repositories.MinIOClient != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        _, err = MinIOClient.PutObject(ctx, bucketName, objectName, file, size, options)
        if err != nil {
            // 自动降级到本地存储
            log.Printf("MinIO上传失败，降级到本地存储: %v", err)
        }
    }
    // 降级：保存到本地 uploads/ 目录
}
```

**未能完全实现的原因**：
- ❌ **网络连接问题**：校园网环境访问外部 MinIO 服务器延迟高（平均 431ms）
- ❌ **上传超时**：即使设置 30 秒超时，网络条件下仍然无法稳定上传
- ❌ **权限限制**：学长提供的账号可能对 bucket 操作权限有限（Access Denied）

**当前解决方案**：
```
✅ 实现了自动降级机制：
MinIO 上传失败 → 自动切换到本地文件存储
```

**技术收获**：
- ✅ 学会了 MinIO Go SDK 的使用
- ✅ 理解了对象存储的工作原理
- ✅ 掌握了超时控制和降级策略
- ✅ 学会了环境变量配置管理

**配置文件**：
```bash
# start_with_minio.bat
set MINIO_ENDPOINT=minio.2huo.tech
set MINIO_ACCESS_KEY=***
set MINIO_USE_SSL=true
set MINIO_BUCKET_NAME=traning-hotel
```

**反思**：
- 生产环境的网络条件和开发环境可能差异很大
- 外部服务依赖需要考虑网络稳定性和降级方案
- 对象存储更适合部署在同一内网或使用 CDN 加速

---

### 未实现的其他功能

#### 2. 限流防护
**问题**：
- 没有防止恶意请求
- 没有限制API调用频率
- 可能被刷接口

**解决方案**：
```go
// 可以使用中间件实现限流
// 例如：golang.org/x/time/rate
import "golang.org/x/time/rate"

func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(10, 20)  // 每秒10次，最大20次
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

#### 3. 日志系统
**可以改进**：
```go
// 当前：使用Gin默认日志
// 可以改进：使用结构化日志（zap）
import "go.uber.org/zap"

logger.Info("create luggage",
    zap.String("username", username),
    zap.Int64("luggage_id", id),
    zap.String("retrieval_code", code),
)
```

#### 4. 监控告警
- 接口响应时间监控
- 错误率监控
- 数据库慢查询监控

#### 5. 单元测试
- Handler层测试
- Service层测试
- Repository层测试

---

## 第30页：未来改进方向
### 功能扩展
- 📱 小程序端开发
- 📊 数据统计和可视化
- 🔔 消息通知（微信/短信）
- 📷 图片压缩和CDN

### 技术优化
- 🚀 性能优化（查询优化、缓存策略）
- 🔒 安全加固（限流、防刷、加密）
- 📝 完善日志和监控
- 🧪 增加单元测试覆盖

### 架构演进
- 📦 微服务拆分（如果业务复杂）
- 🐳 容器化部署（Docker）
- ☸️ K8s编排（如果需要）
- 📈 可观测性（Prometheus + Grafana）

---

## 第31页：总结
### 项目完成情况
- ✅ 完成了核心功能开发
- ✅ 实现了前后端分离
- ✅ 编写了详细文档
- ✅ 实现了基础优化（Redis缓存）

### 个人成长
- ✅ 掌握了Go后端开发
- ✅ 理解了分层架构
- ✅ 学会了前后端协作
- ✅ 积累了实战经验

### 不足与反思
- ⚠️ 前期设计不够充分
- ⚠️ 频繁返工浪费时间
- ⚠️ 部分功能未实现（限流等）

### 最大收获
> **学会了如何从0到1完成一个完整的后端项目**
> 
> 不仅是写代码，更重要的是：
> - 需求分析
> - 架构设计
> - 接口设计
> - 团队协作
> - 问题解决

---

## 第32页：Q&A
**欢迎提问**

### 可能的问题准备
1. JWT的安全性如何保证？
2. 为什么不用微服务架构？
3. 如何处理并发问题？
4. 数据库如何优化？
5. 如何保证数据一致性？

链接：[https://github.com/Zhengxin-Yu/Winter_Vacation_Intensive_Training_Program]
