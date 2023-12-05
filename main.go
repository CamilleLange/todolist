package main

import (
	"context"
	"fmt"
	"net/http"

	// #nosec
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Aloe-Corporation/logs"
	"github.com/CamilleLange/todolist/internal/configuration"
	"github.com/CamilleLange/todolist/internal/connectors"
	"github.com/CamilleLange/todolist/internal/ginrouters"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	PREFIX_ENV          = "TASK_API_"
	ENV_CONFIG          = PREFIX_ENV + "CONFIG"
	DEFAULT_PATH_CONFIG = "/config/"
)

var log = logs.Get()

func main() {
	go func() {
		log.Error("pprof", zap.Error(http.ListenAndServe("0.0.0.0:6060", nil)))
	}()

	if err := Init(); err != nil {
		panic(fmt.Errorf("fail to init API: %w", err))
	}

	// Get the router
	routerGin := ginrouters.Router
	addrGin := ginrouters.Config.Addr + ":" + strconv.Itoa(ginrouters.Config.Port)
	srv := &http.Server{
		ReadHeaderTimeout: time.Millisecond,
		Addr:              addrGin,
		Handler:           routerGin,
	}

	go RunGin(addrGin, routerGin)

	WaitSignalShutdown(srv)
}

func Init() error {
	// Load configuration
	log.Info("loading config...")
	pathFileConfig, present := os.LookupEnv(ENV_CONFIG)
	if !present {
		pathFileConfig = DEFAULT_PATH_CONFIG
	}

	err := configuration.LoadConf(pathFileConfig, PREFIX_ENV)
	if err != nil {
		return fmt.Errorf("fail to load config: %w", err)
	}
	log.Info("config is loaded")

	// Init modules
	log.Info("init all modules...")
	err = configuration.InitAllModules()
	if err != nil {
		return fmt.Errorf("fail to init modules: %w", err)
	}
	log.Info("all modules are ready")

	// Init all internal package
	log.Info("init all packages")
	err = configuration.InitAllPkg()
	if err != nil {
		return fmt.Errorf("fail to init packages: %w", err)
	}
	log.Info("all packages are ready")

	return nil
}

func RunGin(addr string, engine *gin.Engine) {
	log.Info("REST API listening on : "+addr,
		zap.String("package", "main"))

	log.Error(engine.Run(addr).Error(),
		zap.String("package", "main"))
}

func WaitSignalShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	if err := connectors.Close(); err != nil {
		log.Error("error during repositories.Close()", zap.Error(err))
	}

	// Time to wait before close forcing
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ginrouters.Config.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown: ", zap.Error(err))
	}

	log.Info("Server exiting")
}
