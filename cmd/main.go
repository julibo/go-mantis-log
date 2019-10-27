package main

import (
	"context"
	"fmt"
	log "go-mantis-log"
	"go-mantis-log/conf"
	"go-mantis-log/plugins/zaplog"
	"go-mantis-log/tracer"

	"time"
)

func main() {
	log.SetLogger(zaplog.New(
		conf.WithProjectName("zap"),
		conf.WithLogPath("tmp/log"),
		conf.WithLogLevel("info"),
		conf.WithIsStdOut("yes"),
	))

	ctx := context.WithValue(context.Background(), tracer.LogTraceKey, "46b1506e7332f7c1:7f75737aa70629cc:3bb947500f42ad71:1")
	// log.Infof("hello %s", "world", ctx)
	log.Infof("this is zap test %s", "test", ctx)

	ticker := time.NewTicker(time.Second)
	for i := 0; i < 300; i++ {
		fmt.Println(<-ticker.C)
	}
}
