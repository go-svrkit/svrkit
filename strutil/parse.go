// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"cmp"
	"fmt"
	"strings"

	"gopkg.in/svrkit.v1/logger"
)

const (
	KVPairQuote    = '\''
	SepVerticalBar = "|"
	SepColon       = ":"
	SepComma       = ","
	SepEqualSign   = "="
)

// ParseToMap 解析字符串为K-V map，
// 示例： ParseKVPairs("x=1|y=2", SepEqualSign, SepVerticalBar) -> {"a":"x,y", "c":"z"}
func ParseToMap[K, V cmp.Ordered](text string, sep1, sep2 string) (map[K]V, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}
	var arr = strings.Split(text, sep2)
	var dict = make(map[K]V, len(arr))
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		var parts = strings.Split(str, sep1)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key-value format: %s", str)
		}
		key, err := ParseTo[K](parts[0])
		if err != nil {
			return nil, err
		}
		value, err := ParseTo[V](parts[1])
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}
	return dict, nil
}

func MustParseToMap[K, V cmp.Ordered](text string, sep1, sep2 string) map[K]V {
	dict, err := ParseToMap[K, V](text, sep1, sep2)
	if err != nil {
		logger.Panicf("MustParseToMap: %v", err)
	}
	return dict
}

// ParseToMapN 解析字符串为K-V map，
// 示例： ParseToMapN("x=1|y=2", SepEqualSign, SepVerticalBar) -> {"a":"x,y", "c":"z"}
func ParseToMapN[K, V cmp.Ordered](text string) (map[K]V, error) {
	return ParseToMap[K, V](text, SepEqualSign, SepVerticalBar)
}

func MustParseToMapN[K, V cmp.Ordered](text string) map[K]V {
	dict, err := ParseToMap[K, V](text, SepEqualSign, SepVerticalBar)
	if err != nil {
		logger.Panicf("MustParseToMapN: %v", err)
	}
	return dict
}

func unquote(s string) string {
	n := len(s)
	if n >= 2 {
		if s[0] == KVPairQuote && s[n-1] == KVPairQuote {
			return s[1 : n-1]
		}
	}
	return s
}

// ParseKVPairs 解析字符串为K-V map，处理逗号分隔的特殊情况
// 示例： ParseKVPairs("100='x,y',101=z") -> {100:"x,y", 101:"z"}
func ParseKVPairs[K, V cmp.Ordered](text string) (map[K]V, error) {
	const sep1, sep2 = ',', '='
	var dict = make(map[K]V)
	var keyStr string
	var inQuote = false
	var start = 0
	for i := 0; i < len(text); i++ {
		var ch = text[i]
		switch ch {
		case sep1:
			if !inQuote {
				value := strings.TrimSpace(text[start:i])
				if keyStr == "" {
					keyStr = value
					value = ""
				}
				key, err := ParseTo[K](keyStr)
				if err != nil {
					return nil, err
				}
				val, err := ParseTo[V](unquote(value))
				if err != nil {
					return nil, err
				}
				dict[key] = val
				keyStr = ""
				start = i + 1
			}
		case sep2:
			if !inQuote {
				keyStr = strings.TrimSpace(text[start:i])
				start = i + 1
			}
		case KVPairQuote:
			inQuote = !inQuote
		}
	}
	if start < len(text) || keyStr != "" {
		s := strings.TrimSpace(text[start:])
		if keyStr == "" {
			keyStr = s
			s = ""
		}
		key, err := ParseTo[K](keyStr)
		if err != nil {
			return nil, err
		}
		val, err := ParseTo[V](unquote(s))
		if err != nil {
			return nil, err
		}
		dict[key] = val
	}
	return dict, nil
}

func MustParseKVPairs[K, V cmp.Ordered](text string) map[K]V {
	dict, err := ParseKVPairs[K, V](text)
	if err != nil {
		logger.Panicf("MustParseKVPairs: %v", err)
	}
	return dict
}
