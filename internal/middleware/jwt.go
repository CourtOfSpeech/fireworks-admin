// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// JWTConfig 定义 JWT 中间件的配置结构。
type JWTConfig struct {
	// Secret 是 JWT 签名密钥。
	Secret string
	// ExpireTime 是令牌过期时间，单位为小时。
	ExpireTime int
}

// defaultJWTConfig 是项目级别的默认 JWT 配置。
var defaultJWTConfig = JWTConfig{
	Secret:     "default-secret-key-please-change-in-production",
	ExpireTime: 24,
}

// SetDefaultJWTConfig 设置项目级别的默认 JWT 配置。
// 用于在应用启动时统一配置 JWT 参数。
func SetDefaultJWTConfig(config JWTConfig) {
	if config.Secret != "" {
		defaultJWTConfig.Secret = config.Secret
	}
	if config.ExpireTime > 0 {
		defaultJWTConfig.ExpireTime = config.ExpireTime
	}
}

// Skipper 定义跳过中间件的函数类型。
type Skipper func(c *echo.Context) bool

// jwtMiddlewareConfig 是 JWT 中间件内部配置结构。
type jwtMiddlewareConfig struct {
	// signingKey 是签名密钥。
	signingKey []byte
	// signingMethod 是签名算法。
	signingMethod jwt.SigningMethod
	// tokenLookup 是令牌查找位置。
	tokenLookup string
	// authScheme 是认证方案。
	authScheme string
	// skipper 是跳过中间件的函数。
	skipper Skipper
	// errorHandler 是错误处理函数。
	errorHandler func(error) error
	// successHandler 是认证成功后的回调函数。
	successHandler func(*echo.Context)
}

// NewJWTMiddleware 创建 JWT 认证中间件。
// 使用默认配置创建 JWT 中间件，适用于需要 JWT 认证的路由组。
func NewJWTMiddleware(config *JWTConfig) echo.MiddlewareFunc {
	return NewJWTMiddlewareWithSkipper(config, nil)
}

// NewJWTMiddlewareWithSkipper 创建带跳过功能的 JWT 认证中间件。
// 允许通过 Skipper 函数跳过某些路由的 JWT 认证。
func NewJWTMiddlewareWithSkipper(config *JWTConfig, skipper Skipper) echo.MiddlewareFunc {
	return NewJWTMiddlewareWithHandler(config, skipper, nil, nil)
}

// NewJWTMiddlewareWithHandler 创建完全自定义的 JWT 认证中间件。
// 提供完整的自定义选项，包括错误处理和成功处理回调。
func NewJWTMiddlewareWithHandler(config *JWTConfig, skipper Skipper, errorHandler func(error) error, successHandler func(*echo.Context)) echo.MiddlewareFunc {
	cfg := defaultJWTConfig
	if config != nil {
		if config.Secret != "" {
			cfg.Secret = config.Secret
		}
		if config.ExpireTime > 0 {
			cfg.ExpireTime = config.ExpireTime
		}
	}

	if skipper == nil {
		skipper = func(c *echo.Context) bool {
			return false
		}
	}

	if errorHandler == nil {
		errorHandler = func(err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("JWT认证失败: %v", err))
		}
	}

	middlewareConfig := &jwtMiddlewareConfig{
		signingKey:     []byte(cfg.Secret),
		signingMethod:  jwt.SigningMethodHS256,
		tokenLookup:    "header:Authorization:Bearer ",
		authScheme:     "Bearer",
		skipper:        skipper,
		errorHandler:   errorHandler,
		successHandler: successHandler,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if middlewareConfig.skipper(c) {
				return next(c)
			}

			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return middlewareConfig.errorHandler(errors.New("缺少Authorization头"))
			}

			if !strings.HasPrefix(auth, middlewareConfig.authScheme+" ") {
				return middlewareConfig.errorHandler(errors.New("无效的认证方案"))
			}

			tokenString := strings.TrimPrefix(auth, middlewareConfig.authScheme+" ")
			if tokenString == "" {
				return middlewareConfig.errorHandler(errors.New("令牌为空"))
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != middlewareConfig.signingMethod.Alg() {
					return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
				}
				return middlewareConfig.signingKey, nil
			})

			if err != nil {
				return middlewareConfig.errorHandler(err)
			}

			if !token.Valid {
				return middlewareConfig.errorHandler(errors.New("令牌无效"))
			}

			c.Set("user", token)

			if middlewareConfig.successHandler != nil {
				middlewareConfig.successHandler(c)
			}

			return next(c)
		}
	}
}

// GetUserIDFromToken 从 JWT 令牌中获取用户 ID。
// 在认证成功后的处理器中调用，从令牌的 claims 中提取 user_id 字段。
func GetUserIDFromToken(c *echo.Context) (string, error) {
	user := c.Get("user")
	if user == nil {
		return "", errors.New("未找到JWT令牌信息")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return "", errors.New("令牌类型错误")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("令牌声明类型错误")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("令牌中未找到user_id")
	}

	return userID, nil
}

// GetUsernameFromToken 从 JWT 令牌中获取用户名。
// 在认证成功后的处理器中调用，从令牌的 claims 中提取 username 字段。
func GetUsernameFromToken(c *echo.Context) (string, error) {
	user := c.Get("user")
	if user == nil {
		return "", errors.New("未找到JWT令牌信息")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return "", errors.New("令牌类型错误")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("令牌声明类型错误")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("令牌中未找到username")
	}

	return username, nil
}

// GetClaimFromToken 从 JWT 令牌中获取指定的 claim 值。
// 通用的 claim 提取函数，可以获取任意 key 对应的值。
func GetClaimFromToken(c *echo.Context, claimKey string) (interface{}, error) {
	user := c.Get("user")
	if user == nil {
		return nil, errors.New("未找到JWT令牌信息")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, errors.New("令牌类型错误")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("令牌声明类型错误")
	}

	value, ok := claims[claimKey]
	if !ok {
		return nil, fmt.Errorf("令牌中未找到%s", claimKey)
	}

	return value, nil
}

// JWTErrorResponse 生成 JWT 错误响应。
// 用于在处理器中返回统一的 JWT 认证错误响应。
func JWTErrorResponse(c *echo.Context, err error) error {
	return bizerr.New(http.StatusUnauthorized, fmt.Sprintf("认证失败: %v", err), http.StatusUnauthorized)
}
