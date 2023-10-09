// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"strings"
)

const (
	KVPairQuote    = '\''
	SepVerticalBar = "|"
	SepColon       = ":"
	SepComma       = ","
	SepEqualSign   = "="
)

func unquote(s string) string {
	n := len(s)
	if n >= 2 {
		if s[0] == KVPairQuote && s[n-1] == KVPairQuote {
			return s[1 : n-1]
		}
	}
	return s
}

// 解析字符串为K-V map，
// 示例：s = "a='x,y',c=z"
//
//	ParseKVPairs(s,",","=") ==> {"a":"x,y", "c":"z"}
func ParseKVPairs(text string, sep, equal byte) map[string]string {
	var result = make(map[string]string)
	var key string
	var inQuote = false
	var start = 0
	for i := 0; i < len(text); i++ {
		var ch = text[i]
		switch ch {
		case sep:
			if !inQuote {
				value := strings.TrimSpace(text[start:i])
				if key == "" {
					key = value
					value = ""
				}
				result[key] = unquote(value)
				key = ""
				start = i + 1
			}
		case equal:
			if !inQuote {
				key = strings.TrimSpace(text[start:i])
				start = i + 1
			}
		case KVPairQuote:
			inQuote = !inQuote
		}
	}
	if start < len(text) || key != "" {
		s := strings.TrimSpace(text[start:])
		if key == "" {
			key = s
			s = ""
		}
		result[key] = unquote(s)
	}
	return result
}

//// 解析字符串为map, 格式：1:100|2:200|3:300
//func ParseKVMap(text string, sep1, sep2 string) map[int32]int32 {
//	text = strings.TrimSpace(text)
//	if text == "" {
//		return nil
//	}
//	if sep1 == "" {
//		sep1 = SepVerticalBar
//	}
//	if sep2 == "" {
//		sep2 = SepColon
//	}
//	var ret = make(map[int32]int32)
//	var pairs = strings.Split(text, sep1)
//	for _, pair := range pairs {
//		pair = strings.TrimSpace(pair)
//		if pair == "" {
//			continue
//		}
//		kv := strings.Split(pair, sep2)
//		if len(kv) == 2 {
//			var key = Parse3(kv[0])
//			var val = Sto32(kv[1])
//			ret[key] += val
//		} else {
//			return nil
//		}
//	}
//	return ret
//}
