package debug

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
