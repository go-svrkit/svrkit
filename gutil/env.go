// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// GetEnv 获取环境变量
func GetEnv(key string, def string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}
	return def
}

// GetEnvInt 获取环境变量int值
func GetEnvInt(key string, def int) int {
	if s := os.Getenv(key); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return def
}

func GetEnvInt64(key string, def int64) int64 {
	if s := os.Getenv(key); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return n
		}
	}
	return def
}

// GetEnvFloat 获取环境变量float值
func GetEnvFloat(key string, def float64) float64 {
	if s := os.Getenv(key); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return def
}

// GetEnvBool 获取环境变量bool值
func GetEnvBool(key string) bool {
	if s := os.Getenv(key); s != "" {
		if strings.EqualFold(s, "ON") {
			return true
		}
		b, _ := strconv.ParseBool(s)
		return b
	}
	return false
}

func LoadDotEnv(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}
	return godotenv.Overload(filenames...)
}
