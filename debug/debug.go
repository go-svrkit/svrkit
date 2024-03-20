package debug

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/svrkit.v1/zlog"
)

const timestampLayout = "2006-01-02T15:04:05.000-0700" // IOS8601

func TraceStack(title string, w io.Writer) {
	var stack = GetCallerStack(1)
	var now = time.Now()
	fmt.Fprintf(w, "%s\nstack traceback[%s] (most recent calls):\n", title, now.Format(timestampLayout))
	fmt.Fprintf(w, "%v \n", stack)
}

func CatchPanic(title string) {
	if v := recover(); v != nil {
		var now = time.Now()
		var stack = GetCallerStack(1)
		fmt.Fprintf(os.Stderr, "%s\nstack traceback[%s] (most recent calls):\n", title, now.Format(timestampLayout))
		fmt.Fprintf(os.Stderr, "%v %v\n", stack, v)
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
