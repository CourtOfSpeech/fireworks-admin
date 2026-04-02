package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/response"
)

// JWTConfig 定义JWT中间件的配置结构
// Secret: JWT签名密钥
// ExpireTime: 令牌过期时间（小时）
type JWTConfig struct {
	Secret     string
	ExpireTime int
}

// defaultJWTConfig 项目级别的默认JWT配置
var defaultJWTConfig = JWTConfig{
	Secret:     "default-secret-key-please-change-in-production",
	ExpireTime: 24,
}

// SetDefaultJWTConfig 设置项目级别的默认JWT配置
// 用于在应用启动时统一配置JWT参数
func SetDefaultJWTConfig(config JWTConfig) {
	if config.Secret != "" {
		defaultJWTConfig.Secret = config.Secret
	}
	if config.ExpireTime > 0 {
		defaultJWTConfig.ExpireTime = config.ExpireTime
	}
}

// Skipper 定义跳过中间件的函数类型
type Skipper func(c *echo.Context) bool

// jwtMiddlewareConfig JWT中间件内部配置结构
type jwtMiddlewareConfig struct {
	signingKey     []byte
	signingMethod  jwt.SigningMethod
	tokenLookup    string
	authScheme     string
	skipper        Skipper
	errorHandler   func(error) error
	successHandler func(*echo.Context)
}

// NewJWTMiddleware 创建JWT认证中间件
// 使用默认配置创建JWT中间件，适用于需要JWT认证的路由组
// 参数:
//   - config: JWT配置，如果为nil则使用默认配置
// 返回:
//   - echo.MiddlewareFunc: Echo中间件函数
func NewJWTMiddleware(config *JWTConfig) echo.MiddlewareFunc {
	return NewJWTMiddlewareWithSkipper(config, nil)
}

// NewJWTMiddlewareWithSkipper 创建带跳过功能的JWT认证中间件
// 允许通过Skipper函数跳过某些路由的JWT认证
// 参数:
//   - config: JWT配置，如果为nil则使用默认配置
//   - skipper: 跳过中间件的函数，返回true时跳过认证
// 返回:
//   - echo.MiddlewareFunc: Echo中间件函数
func NewJWTMiddlewareWithSkipper(config *JWTConfig, skipper Skipper) echo.MiddlewareFunc {
	return NewJWTMiddlewareWithHandler(config, skipper, nil, nil)
}

// NewJWTMiddlewareWithHandler 创建完全自定义的JWT认证中间件
// 提供完整的自定义选项，包括错误处理和成功处理回调
// 参数:
//   - config: JWT配置，如果为nil则使用默认配置
//   - skipper: 跳过中间件的函数，返回true时跳过认证
//   - errorHandler: 自定义错误处理函数
//   - successHandler: 认证成功后的回调函数
// 返回:
//   - echo.MiddlewareFunc: Echo中间件函数
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

// GetUserIDFromToken 从JWT令牌中获取用户ID
// 在认证成功后的处理器中调用，从令牌的claims中提取user_id字段
// 参数:
//   - c: Echo上下文
// 返回:
//   - string: 用户ID
//   - error: 错误信息
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

// GetUsernameFromToken 从JWT令牌中获取用户名
// 在认证成功后的处理器中调用，从令牌的claims中提取username字段
// 参数:
//   - c: Echo上下文
// 返回:
//   - string: 用户名
//   - error: 错误信息
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

// GetClaimFromToken 从JWT令牌中获取指定的claim值
// 通用的claim提取函数，可以获取任意key对应的值
// 参数:
//   - c: Echo上下文
//   - claimKey: 要获取的claim键名
// 返回:
//   - interface{}: claim值
//   - error: 错误信息
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

// JWTErrorResponse 生成JWT错误响应
// 用于在处理器中返回统一的JWT认证错误响应
// 参数:
//   - c: Echo上下文
//   - err: 错误信息
// 返回:
//   - error: HTTP错误
func JWTErrorResponse(c *echo.Context, err error) error {
	return response.Error(c, http.StatusUnauthorized, fmt.Sprintf("认证失败: %v", err))
}
