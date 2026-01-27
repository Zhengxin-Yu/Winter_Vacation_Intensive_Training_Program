package configs

import "os"

// MinIOConfig MinIO对象存储配置
type MinIOConfig struct {
	Endpoint        string // MinIO服务地址，例如：localhost:9000
	AccessKeyID     string // Access Key
	SecretAccessKey string // Secret Key
	UseSSL          bool   // 是否使用HTTPS
	BucketName      string // 存储桶名称
}

// LoadMinIOConfig 从环境变量读取MinIO配置
func LoadMinIOConfig() MinIOConfig {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000" // 默认本地MinIO
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin" // 默认用户名
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin" // 默认密码
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "hotel-luggage" // 默认桶名
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	return MinIOConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		UseSSL:          useSSL,
		BucketName:      bucketName,
	}
}
