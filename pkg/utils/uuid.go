package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	maxRetries = 3
	retryDelay = 10 * time.Millisecond
)

// NewV7 生成 UUID v7，带重试机制。
// 该函数最多重试 3 次，每次重试间隔 10 毫秒，
// 如果 3 次重试后仍失败，则返回错误。
func NewV7() (uuid.UUID, error) {
	var lastErr error
	for range maxRetries {
		id, err := uuid.NewV7()
		if err == nil {
			return id, nil
		}
		lastErr = err
		time.Sleep(retryDelay)
	}
	return uuid.Nil, fmt.Errorf("uuid.NewV7() failed after %d retries: %w", maxRetries, lastErr)
}

// NewV7Safe 生成 UUID v7，带重试机制，失败时返回 Nil。
// 该函数最多重试 3 次，每次重试间隔 10 毫秒，
// 如果 3 次重试后仍失败，则返回 uuid.Nil。
func NewV7Safe() uuid.UUID {
	for range maxRetries {
		id, err := NewV7()
		if err == nil {
			return id
		}
		time.Sleep(retryDelay)
	}
	return uuid.Nil
}

// NewV4 生成 UUID v4。
// 该函数直接调用 uuid.New() 返回一个随机生成的 UUID v4。
func NewV4() uuid.UUID {
	return uuid.New()
}

// ToString 将 UUID 转换为无连字符的字符串。
// 该函数移除 UUID 字符串中的连字符，
// 例如 019449a8-7c3b-7d2e-8f1a-5b3c2d1e0f0a 转换为 019449a87c3b7d2e8f1a5b3c2d1e0f0a。
func ToString(id uuid.UUID) string {
	return strings.ReplaceAll(id.String(), "-", "")
}

// FromString 从字符串解析 UUID。
// 该函数支持带连字符（36位）和不带连字符（32位）两种格式的字符串，
// 如果输入是 32 位无连字符格式，会自动转换为标准 36 位格式后再解析。
func FromString(s string) (uuid.UUID, error) {
	if len(s) == 32 {
		s = s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
	}
	return uuid.Parse(s)
}

// Parse 解析 UUID 字符串，是 FromString 的别名。
// 该函数与 FromString 功能相同，提供更简洁的调用方式。
func Parse(s string) (uuid.UUID, error) {
	return FromString(s)
}
