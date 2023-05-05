package commands

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mix-go/web-skeleton/di"
	"github.com/mix-go/web-skeleton/routes"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xcli/flag"
	"github.com/mix-go/xcli/process"
	"github.com/mix-go/xutil/xenv"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type WebCommand struct {
}

func (t *WebCommand) Main() {
	if flag.Match("d", "daemon").Bool() {
		process.Daemon()
	}

	logger := di.Logrus()
	server := di.Server()
	addr := xenv.Getenv("GIN_ADDR").String(":8080")
	mode := xenv.Getenv("GIN_MODE").String(gin.ReleaseMode)

	// server
	gin.SetMode(mode)
	router := gin.New()
	// logger
	if mode != gin.ReleaseMode {
		handlerFunc := gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				return fmt.Sprintf("%s|%s|%d|%s",
					params.Method,
					params.Path,
					params.StatusCode,
					params.ClientIP,
				)
			},
			Output: logger.Writer(),
		})
		router.Use(handlerFunc)
	}
	routes.Load(router)
	server.Addr = flag.Match("a", "addr").String(addr)
	server.Handler = router

	// signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		logger.Info("Server shutdown")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("Server shutdown error: %s", err)
		}
	}()

	// templates
	router.LoadHTMLGlob(fmt.Sprintf("%s/../templates/*", xcli.App().BasePath))

	// static file
	router.Static("/static", fmt.Sprintf("%s/../public/static", xcli.App().BasePath))
	router.StaticFile("/favicon.ico", fmt.Sprintf("%s/../public/favicon.ico", xcli.App().BasePath))

	// run
	welcome()
	logger.Infof("Server start at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "http: Server closed") {
		panic(err)
	}
}
