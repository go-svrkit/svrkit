// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package debug

import (
	"fmt"
	"io"
	"net/http"
	"runtime/metrics"
	"strings"
	"time"

	"gopkg.in/svrkit.v1/gutil"
	"gopkg.in/svrkit.v1/zlog"
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
		zlog.Error(sb.String())
	}
}

func StartProfiler(addr string) {
	go func() {
		zlog.Infof("listen pprof at %s", addr)
		var httpServer = &http.Server{
			Addr: addr,
		}
		if err := httpServer.ListenAndServe(); err != nil {
			zlog.Infof("%v", err)
		}
	}()
}

// ReadMetrics 读取指定的metrics
// see https://pkg.go.dev/runtime/metrics
func ReadMetrics(category string) map[string]any {
	if category == "" {
		return nil
	}
	var allDesc = metrics.All()
	var names = make([]string, 0, 8)
	for i := 0; i < len(allDesc); i++ {
		if strings.HasPrefix(allDesc[i].Name, "/"+category) {
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
				result[name] = gutil.PrettyBytes(int64(val))
			} else {
				result[name] = val
			}
		case metrics.KindFloat64:
			result[name] = fmt.Sprintf("%v", sample.Value.Float64()) // JSON INF problem

		case metrics.KindFloat64Histogram:
			var histogram = sample.Value.Float64Histogram()
			if histogram != nil {
				result[name] = map[string]any{
					"counts":  histogram.Counts,
					"buckets": fmt.Sprintf("%v", histogram.Buckets), // JSON INF problem
				}
			}
		default:
		}
	}
	return result
}
