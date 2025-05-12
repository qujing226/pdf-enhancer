package main

import (
	"context"
	"database/sql"
	"fmt"
	"log" // 新增导入

	// 新增导入
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"

	"github.com/qujing226/pdf-enhancer/backend/database"
	"github.com/qujing226/pdf-enhancer/backend/handlers"
	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/repository"
	"github.com/qujing226/pdf-enhancer/backend/services"
	"github.com/qujing226/pdf-enhancer/backend/utils"
)

// @title           AI增强报告系统 API
// @version         1.0
// @description     提供用户登录、报告查询和AI摘要生成功能的API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Panic("未找到.env文件，使用环境变量")
	}

	// 设置运行模式
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化数据库连接
	db, err := initDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 初始化MinIO客户端
	minioClient, err := initMinIO()
	if err != nil {
		log.Fatalf("MinIO连接失败: %v", err)
	}

	// 初始化DeepSeek客户端
	deepseekClient := initDeepSeekClient()

	// 初始化服务
	userService := initUserService(db)
	reportService := initReportService(db, minioClient, deepseekClient)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户认证
		authHandler := handlers.NewAuthHandler(userService)
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(authMiddleware())
		{
			// 报告相关API
			reportHandler := handlers.NewReportHandler(reportService)
			auth.GET("/reports", reportHandler.GetReports)
			auth.GET("/report/:report_id", reportHandler.GetReport)
			auth.POST("/report/:report_id/summary", reportHandler.GenerateSummary)
			auth.GET("/report/:report_id/pdf", reportHandler.GetReportPDF)
			auth.POST("/reports/upload", reportHandler.UploadReport)
		}
	}

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	port := getEnv("PORT", "8080")
	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 初始化数据库连接
func initDB() (*sql.DB, error) {
	// 从环境变量创建数据库配置
	config := database.NewConfig()

	// 连接到数据库
	return database.Connect(config)
}

// 初始化MinIO客户端
func initMinIO() (*minio.Client, error) {
	// 从环境变量获取MinIO配置
	endpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	accessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	useSSL, _ := strconv.ParseBool(getEnv("MINIO_USE_SSL", "false"))

	// 创建MinIO客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("MinIO客户端初始化失败: %w", err)
	}

	// 确保存储桶存在
	bucketName := getEnv("MINIO_BUCKET_NAME", "reports")
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("检查MinIO存储桶失败: %w", err)
	}

	if !exists {
		// 创建存储桶
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("创建MinIO存储桶失败: %w", err)
		}
		log.Printf("已创建MinIO存储桶: %s", bucketName)
	}

	log.Printf("成功连接到MinIO服务: %s", endpoint)
	return minioClient, nil
}

// 初始化DeepSeek客户端
func initDeepSeekClient() *services.DeepSeekClient {
	// 从环境变量获取DeepSeek配置
	config := services.DeepSeekConfig{
		APIKey:      getEnv("DEEPSEEK_API_KEY", ""),
		BaseURL:     getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
		ModelName:   getEnv("DEEPSEEK_MODEL_NAME", "deepseek-chat"),
		MaxTokens:   getIntEnv("DEEPSEEK_MAX_TOKENS", 1000),
		Temperature: getFloatEnv("DEEPSEEK_TEMPERATURE", 0.7),
	}

	// 创建DeepSeek客户端
	return services.NewDeepSeekClient(config)
}

// 初始化用户服务
func initUserService(db *sql.DB) *services.UserService {
	userRepo := repository.NewUserRepository(db)
	return services.NewUserService(userRepo)
}

// 初始化报告服务
func initReportService(db *sql.DB, minioClient *minio.Client, deepseekClient *services.DeepSeekClient) *services.ReportService {
	reportRepo := repository.NewReportRepository(db)
	return services.NewReportService(reportRepo, minioClient, deepseekClient, getEnv("MINIO_BUCKET_NAME", "reports"))
}

// JWT认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, models.NewAPIResponse(401, "未提供有效的认证令牌", nil))
			return
		}

		// 提取令牌
		tokenString := authHeader[7:]

		// 解析JWT令牌
		token, err := jwt.ParseWithClaims(tokenString, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 获取公钥
			publicKeyPEM := getEnv("JWT_PUBLIC_KEY", "")
			if publicKeyPEM == "" {
				return nil, fmt.Errorf("未配置JWT公钥")
			}

			// 解析公钥
			publicKey, err := utils.ParseRSAPublicKeyFromPEM(publicKeyPEM)
			if err != nil {
				return nil, fmt.Errorf("解析JWT公钥失败: %w", err)
			}

			return publicKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, models.NewAPIResponse(401, "无效的认证令牌", err.Error()))
			return
		}

		// 验证令牌
		if !token.Valid {
			c.AbortWithStatusJSON(401, models.NewAPIResponse(401, "认证令牌已过期或无效", nil))
			return
		}

		// 提取用户信息
		claims, ok := token.Claims.(*utils.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(401, models.NewAPIResponse(401, "无效的令牌声明", nil))
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// 获取环境变量整数值
func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("警告: 环境变量 %s 不是有效的整数，使用默认值 %d", key, defaultValue)
		return defaultValue
	}

	return value
}

// 获取环境变量浮点数值
func getFloatEnv(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Printf("警告: 环境变量 %s 不是有效的浮点数，使用默认值 %f", key, defaultValue)
		return defaultValue
	}

	return value
}
