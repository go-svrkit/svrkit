package helper

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gopkg.in/svrkit.v1/logger"
)

func StartProfiler(addr string) {
	go func() {
		logger.Infof("listen pprof at %s", addr)
		var httpServer = &http.Server{
			Addr: addr,
		}
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Infof("%v", err)
		}
	}()
}

func LoadDotEnv() {
	var filename = ".env"
	if IsFileExist(filename) {
		if err := godotenv.Load(filename); err != nil {
			log.Printf("load %s failed: %v", filename, err)
		} else {
			log.Printf("load %s OK", filename)
		}
	}
}
