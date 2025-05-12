package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// 安全相关常量
const (
	// Argon2参数
	Argon2Time    = 3
	Argon2Memory  = 64 * 1024
	Argon2Threads = 4
	Argon2KeyLen  = 32

	// JWT过期时间（秒）
	JWTExpiration = 3600
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GeneratePasswordHash 使用Argon2id生成密码哈希
func GeneratePasswordHash(password string, salt []byte) (string, error) {
	if salt == nil {
		// 生成随机盐值
		salt = make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return "", err
		}
	}

	// 使用Argon2id生成哈希
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLen,
	)

	// 编码为Base64字符串
	hashedPassword := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		Argon2Memory,
		Argon2Time,
		Argon2Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return hashedPassword, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) (bool, error) {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("无效的哈希格式")
	}

	var version int
	var memory uint32
	var time uint32
	var parallelism uint8

	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false, err
	}

	// 解码盐值和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// 使用相同参数计算哈希
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		parallelism,
		uint32(len(decodedHash)),
	)

	// 比较哈希值
	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1, nil
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID, email, privateKeyPEM string) (string, error) {
	// 解析私钥
	privateKey, err := ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// 创建声明
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-enhance-test",
			Subject:   userID,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT 验证JWT令牌
func VerifyJWT(tokenString, publicKeyPEM string) (*JWTClaims, error) {
	// 解析公钥
	publicKey, err := ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// EncryptData 使用AES-GCM加密数据
func EncryptData(plaintext []byte, key []byte) ([]byte, error) {
	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptData 使用AES-GCM解密数据
func DecryptData(ciphertext []byte, key []byte) ([]byte, error) {
	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 检查长度
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("密文太短")
	}

	// 提取随机数
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// ParseRSAPrivateKeyFromPEM 解析PEM格式的RSA私钥
func ParseRSAPrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("无法解码PEM块")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("无法解析私钥")
	}

	return key, nil
}

// ParseRSAPublicKeyFromPEM 解析PEM格式的RSA公钥
func ParseRSAPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("无法解码PEM块")
	}

	// 使用 ParsePKIXPublicKey 替代 ParsePKCS1PublicKey
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是RSA公钥")
	}

	return rsaPubKey, nil
}

func GetUserIDFromContext(c *gin.Context) string {
	id, _ := c.Get("userID")
	idstr, ok := id.(string)
	if !ok {
		log.Println("GetUserIDFromContext error")
		return ""
	}
	return idstr

}
