// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"runtime/metrics"
	"strings"
	"time"

	"gopkg.in/svrkit.v1/qlog"
)

const timestampLayout = "2006-01-02T15:04:05.000-0700" // IOS8601

func TraceStack(skip int, title string, err interface{}, w io.Writer) {
	var stack = GetCallerStack(skip + 1)
	var now = time.Now()
	fmt.Fprintf(w, "%s\n%v\nstack traceback[%s] (most recent calls):\n", title, err, now.Format(timestampLayout))
	fmt.Fprintf(w, "%v\n", stack)
}

func CatchPanic(title string) {
	if v := recover(); v != nil {
		var stack = GetCallerStack(1)
		var now = time.Now()
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s\n%v\nstack traceback[%s] (most recent calls):", title, v, now.Format(timestampLayout))
		fmt.Fprintf(&sb, "%v\n", stack)
		qlog.Error(sb.String())
	}
}

func StartProfiler(addr string) {
	go func() {
		qlog.Infof("listen pprof at %s", addr)
		var httpServer = &http.Server{
			Addr: addr,
		}
		if err := httpServer.ListenAndServe(); err != nil {
			qlog.Infof("%v", err)
		}
	}()
}

// ReadGCPercent see https://pkg.go.dev/runtime/debug#SetGCPercent
func ReadGCPercent() uint64 {
	var sample = []metrics.Sample{{Name: "/gc/gogc:percent"}}
	metrics.Read(sample)
	return sample[0].Value.Uint64()
}

// ReadMemoryLimit see https://pkg.go.dev/runtime/debug#SetMemoryLimit
func ReadMemoryLimit() string {
	var sample = []metrics.Sample{{Name: "/gc/gomemlimit:bytes"}}
	metrics.Read(sample)
	var bytes = sample[0].Value.Uint64()
	if bytes == math.MaxInt64 {
		return "MaxInt64"
	}
	return PrettyBytes(int64(bytes))
}

// ReadMetrics 读取指定的metrics see https://pkg.go.dev/runtime/metrics
func ReadMetrics(category string) map[string]any {
	if category == "" {
		return nil
	}
	var allDesc = metrics.All()
	var names = make([]string, 0, 8)
	var prefix = "/" + category
	for i := 0; i < len(allDesc); i++ {
		if strings.HasPrefix(allDesc[i].Name, prefix) {
			names = append(names, allDesc[i].Name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	var samples = make([]metrics.Sample, len(names))
	for i, name := range names {
		samples[i].Name = name
	}
	metrics.Read(samples)
	var result = make(map[string]any)
	for i := 0; i < len(samples); i++ {
		var name = names[i]
		var sample = &samples[i]
		switch sample.Value.Kind() {
		case metrics.KindUint64:
			var val = sample.Value.Uint64()
			if strings.HasSuffix(name, "bytes") {
				result[name] = PrettyBytes(int64(val))
			} else {
				result[name] = val
			}
		case metrics.KindFloat64:
			result[name] = fmt.Sprintf("%v", sample.Value.Float64()) // JSON INF problem

		case metrics.KindFloat64Histogram:
			var histogram = sample.Value.Float64Histogram()
			if histogram != nil {
				var buckets = make([]string, 0, len(histogram.Buckets))
				for j := 0; j < len(histogram.Buckets); j++ {
					if j >= len(histogram.Counts) || histogram.Counts[j] == 0 {
						continue
					}
					var text = fmt.Sprintf("%v = %d", histogram.Buckets[j], histogram.Counts[j])
					buckets = append(buckets, text)
				}
				result[name] = buckets
			}
		default:
		}
	}
	return result
}
