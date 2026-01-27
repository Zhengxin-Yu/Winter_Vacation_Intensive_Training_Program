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

var MinIOClient *minio.Client
var MinIOBucketName string

// InitMinIO 初始化MinIO客户端
func InitMinIO() {
	config := configs.LoadMinIOConfig()

	// 创建MinIO客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Printf("⚠️  MinIO初始化失败: %v (将使用本地文件存储)", err)
		return
	}

	// 测试连接 - 尝试检查 bucket 是否存在（设置5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, config.BucketName)
	
	if err != nil {
		// 如果是权限错误(Access Denied)，可能是没有ListBucket权限
		// 但bucket可能存在，我们继续尝试使用它
		log.Printf("⚠️  无法检查bucket状态: %v (尝试直接使用bucket)", err)
		// 不要return，继续执行
	} else if !exists {
		// bucket不存在，尝试创建
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("⚠️  创建MinIO bucket失败: %v (bucket可能已存在，尝试继续)", err)
			// 不要return，可能是bucket已存在但我们没有创建权限
		} else {
			log.Printf("✅ MinIO bucket '%s' 创建成功", config.BucketName)
		}
	}

	// 尝试设置bucket为公开读取（可选）
	// 注意：如果没有权限，这个操作可能会失败，但不影响上传功能
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
		log.Printf("ℹ️  设置bucket策略失败(可忽略): %v", err)
	}

	MinIOClient = client
	MinIOBucketName = config.BucketName
	log.Printf("✅ MinIO初始化成功: %s/%s", config.Endpoint, config.BucketName)
}
