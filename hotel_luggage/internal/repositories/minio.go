package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"hotel_luggage/configs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient 全局 MinIO 客户端（可为空）
// 说明：
// - MinIO 是一个兼容 Amazon S3 API 的对象存储服务
// - 用于存储上传的行李照片，相比本地文件系统更适合分布式部署
// - 如果 MinIO 连接失败，系统会自动降级到本地文件存储（./uploads 目录）
// - 使用前需要判断 MinIOClient != nil
//
// 设计理念：MinIO 是可选的存储优化组件，不影响核心功能
var MinIOClient *minio.Client

// MinIOBucketName MinIO 存储桶名称
// Bucket（存储桶）相当于文件系统中的"根目录"，所有文件都存储在 bucket 中
var MinIOBucketName string

// InitMinIO 初始化 MinIO 对象存储客户端（失败则自动降级到本地存储）
// 功能：
// 1. 从环境变量读取 MinIO 配置（服务器地址、凭证等）
// 2. 创建 MinIO 客户端并测试连接
// 3. 检查 bucket 是否存在，不存在则创建
// 4. 设置 bucket 为公开读取（可选，方便直接访问图片）
// 5. 连接失败：打印警告日志，MinIOClient 设为 nil（降级到本地存储）
// 6. 连接成功：设置全局 MinIOClient 和 MinIOBucketName
//
// 环境变量配置：
//   MINIO_ENDPOINT        - MinIO 服务器地址（如：localhost:9000 或 minio.example.com）
//   MINIO_ACCESS_KEY      - 访问密钥（Access Key ID）
//   MINIO_SECRET_KEY      - 私密密钥（Secret Access Key）
//   MINIO_USE_SSL         - 是否使用 HTTPS（true/false）
//   MINIO_BUCKET_NAME     - 存储桶名称（如：hotel-luggage）
//
// 设置示例（Windows）：
//   set MINIO_ENDPOINT=minio.2huo.tech
//   set MINIO_ACCESS_KEY=minioadmin
//   set MINIO_SECRET_KEY=minioadmin
//   set MINIO_USE_SSL=true
//   set MINIO_BUCKET_NAME=hotel-luggage
//
// 降级策略：
//   MinIO 连接失败不会导致程序退出，而是打印警告日志并继续运行
//   Upload handler 会自动降级到本地文件存储（./uploads 目录）
//
// 容错设计：
//   - 权限不足时（如无 ListBucket 权限），仍尝试使用 bucket
//   - bucket 创建失败时（可能已存在），不中断初始化
//   - 设置 bucket 策略失败时（权限不足），不影响上传功能
//   - 所有操作都有 5 秒超时，避免长时间等待
func InitMinIO() {
	// 1. 加载 MinIO 配置（从环境变量读取）
	config := configs.LoadMinIOConfig()

	// 2. 创建 MinIO 客户端
	// - Endpoint: MinIO 服务器地址
	// - Creds: 使用静态凭证（Access Key + Secret Key）
	// - Secure: 是否使用 HTTPS（true: https://, false: http://）
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		// 客户端创建失败：打印警告，降级到本地存储
		log.Printf("⚠️  MinIO初始化失败: %v (将使用本地文件存储)", err)
		return
	}

	// 3. 测试连接 - 检查 bucket 是否存在（5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, config.BucketName)
	
	if err != nil {
		// 容错处理：如果是权限错误（Access Denied），可能是没有 ListBucket 权限
		// 但 bucket 可能存在，我们继续尝试使用它（后续 PutObject 可能有权限）
		log.Printf("⚠️  无法检查bucket状态: %v (尝试直接使用bucket)", err)
		// 不要 return，继续执行后续步骤
	} else if !exists {
		// bucket 不存在，尝试创建
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			// 创建失败：可能是 bucket 已存在但我们没有创建权限
			log.Printf("⚠️  创建MinIO bucket失败: %v (bucket可能已存在，尝试继续)", err)
			// 不要 return，可能 bucket 已存在，后续上传可能成功
		} else {
			log.Printf("✅ MinIO bucket '%s' 创建成功", config.BucketName)
		}
	}

	// 4. 尝试设置 bucket 为公开读取（可选，方便直接访问图片）
	// 策略说明：允许任何人（Principal: *）执行 GetObject 操作（下载文件）
	// 注意：如果没有 SetBucketPolicy 权限，这个操作会失败，但不影响上传功能
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, config.BucketName)

	err = client.SetBucketPolicy(ctx, config.BucketName, policy)
	if err != nil {
		// 策略设置失败：可能是权限不足，但不影响核心上传功能
		log.Printf("ℹ️  设置bucket策略失败(可忽略): %v", err)
	}

	// 5. 初始化成功：设置全局变量
	MinIOClient = client
	MinIOBucketName = config.BucketName
	log.Printf("✅ MinIO初始化成功: %s/%s", config.Endpoint, config.BucketName)
}
