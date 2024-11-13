// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"bytes"
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

const digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
const digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"

// escapeBytesBackslash escapes []byte with backslashes (\)
// This escapes the contents of a string (provided as []byte) by adding backslashes before special
// characters, and turning others into specific escape sequences, such as
// turning newlines into \n and null bytes into \0.
// https://github.com/mysql/mysql-server/blob/mysql-5.7.5/mysys/charset.c#L823-L932
func escapeBytesBackslash(buf *bytes.Buffer, v []byte) {
	for _, c := range v {
		switch c {
		case '\x00':
			buf.WriteByte('\\')
			buf.WriteByte('0')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\x1a':
			buf.WriteByte('\\')
			buf.WriteByte('Z')
		case '\'':
			buf.WriteByte('\\')
			buf.WriteByte('\'')
		case '"':
			buf.WriteByte('\\')
			buf.WriteByte('"')
		case '\\':
			buf.WriteByte('\\')
			buf.WriteByte('\\')
		default:
			buf.WriteByte(c)
		}
	}
}

// escapeBytesQuotes escapes apostrophes in []byte by doubling them up.
// This escapes the contents of a string by doubling up any apostrophes that
// it contains. This is used when the NO_BACKSLASH_ESCAPES SQL_MODE is in
// effect on the server.
// https://github.com/mysql/mysql-server/blob/mysql-5.7.5/mysys/charset.c#L963-L1038
func escapeBytesQuotes(buf *bytes.Buffer, v []byte) {
	for _, c := range v {
		if c == '\'' {
			buf.WriteByte('\'')
			buf.WriteByte('\'')
		} else {
			buf.WriteByte(c)
		}
	}
}

// escapeStringQuotes is similar to escapeBytesQuotes but for string.
func escapeStringQuotes(buf *bytes.Buffer, v string) {
	for i := 0; i < len(v); i++ {
		c := v[i]
		if c == '\'' {
			buf.WriteByte('\'')
			buf.WriteByte('\'')
		} else {
			buf.WriteByte(c)
		}
	}
}

func InterpolateParams(buf *bytes.Buffer, query string, args []interface{}) error {
	// Number of ? should be same to len(args)
	if n := strings.Count(query, "?"); n != len(args) {
		return driver.ErrSkip
	}

	argPos := 0

	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf.WriteString(query[i:])
			break
		}
		buf.WriteString(query[i : i+q])
		i += q

		arg, err := driver.DefaultParameterConverter.ConvertValue(args[argPos])
		if err != nil {
			return driver.ErrSkip
		}
		argPos++

		if arg == nil {
			buf.WriteString("NULL")
			continue
		}

		switch v := arg.(type) {
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case float64:
			buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
		case bool:
			if v {
				buf.WriteByte('1')
			} else {
				buf.WriteByte('0')
			}
		case time.Time:
			if v.IsZero() {
				buf.WriteString("'0000-00-00'")
			} else {
				v := v.In(time.Local)
				v = v.Add(time.Nanosecond * 500) // To round under microsecond
				year := v.Year()
				year100 := year / 100
				year1 := year % 100
				month := v.Month()
				day := v.Day()
				hour := v.Hour()
				minute := v.Minute()
				second := v.Second()
				micro := v.Nanosecond() / 1000

				buf.Write([]byte{
					'\'',
					digits10[year100], digits01[year100],
					digits10[year1], digits01[year1],
					'-',
					digits10[month], digits01[month],
					'-',
					digits10[day], digits01[day],
					' ',
					digits10[hour], digits01[hour],
					':',
					digits10[minute], digits01[minute],
					':',
					digits10[second], digits01[second],
				})

				if micro != 0 {
					micro10000 := micro / 10000
					micro100 := micro / 100 % 100
					micro1 := micro % 100
					buf.Write([]byte{
						'.',
						digits10[micro10000], digits01[micro10000],
						digits10[micro100], digits01[micro100],
						digits10[micro1], digits01[micro1],
					})
				}
				buf.WriteByte('\'')
			}
		case []byte:
			if v == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString("_binary'")
				escapeBytesQuotes(buf, v)
				buf.WriteByte('\'')
			}
		case string:
			buf.WriteByte('\'')
			escapeStringQuotes(buf, v)
			buf.WriteByte('\'')
		default:
			return driver.ErrSkip
		}
	}
	if argPos != len(args) {
		return driver.ErrSkip
	}
	return nil
}
