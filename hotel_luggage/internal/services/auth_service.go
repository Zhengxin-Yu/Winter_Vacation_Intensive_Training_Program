package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login 用户登录验证（用户名密码校验）
// 功能：
// 1. 验证用户名和密码是否为空
// 2. 从数据库查询用户信息
// 3. 使用 bcrypt 验证密码哈希
// 4. 验证成功返回用户信息，失败返回统一的错误信息（防止用户名枚举攻击）
//
// 参数：
//   - username: 用户名
//   - password: 明文密码
//
// 返回：
//   - models.User: 用户信息（包含 username, role, hotel_id 等）
//   - error: 验证失败时返回错误
//
// 安全设计：
// 1. 密码使用 bcrypt 加密存储（不存储明文）
// 2. 登录失败时返回统一错误信息 "invalid username or password"
//    - 防止攻击者通过错误信息判断用户名是否存在
//    - 无论是用户名不存在还是密码错误，都返回相同的错误
// 3. bcrypt 自带 salt（随机盐），防止彩虹表攻击
//
// 使用示例：
//   user, err := services.Login("staff_user", "123456")
//   if err != nil {
//       // 登录失败
//       return gin.H{"message": "login failed"}
//   }
//   // 生成 JWT token
//   token, _ := utils.GenerateToken(user.Username, user.Role)
//
// 后续优化建议：
// - 添加登录失败次数限制（防止暴力破解）
// - 添加账号锁定机制（连续失败 N 次后锁定）
// - 记录登录日志（时间、IP、是否成功）
func Login(username, password string) (models.User, error) {
	// 1. 参数验证：用户名和密码不能为空
	if username == "" || password == "" {
		return models.User{}, errors.New("username or password is empty")
	}

	// 2. 查询用户信息
	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		// 用户不存在：返回统一的错误信息（不暴露用户名是否存在）
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("invalid username or password")
		}
		// 数据库错误：返回原始错误
		return models.User{}, err
	}

	// 3. 验证密码：使用 bcrypt 对比密码哈希
	// bcrypt.CompareHashAndPassword() 会：
	// - 从哈希中提取 salt
	// - 使用相同的 salt 对输入密码进行哈希
	// - 比较两个哈希值是否相同
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// 密码错误：返回统一的错误信息（不暴露密码错误）
		return models.User{}, errors.New("invalid username or password")
	}

	// 4. 验证成功：返回用户信息
	return user, nil
}
