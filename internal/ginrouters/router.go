package ginrouters

import (
	"net/http"
	"time"

	"github.com/Aloe-Corporation/cors"
	"github.com/Aloe-Corporation/logs"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	log = logs.Get()
	// Config of the routers package.
	Config Conf
	// Router of the API.
	Router *gin.Engine

	// Validate singleton, used by all routers to validate request body data.
	ValidateInstance *validator.Validate
)

type HomeMessage struct {
	Service   string `json:"service"`
	Copyright string `json:"copyright"`
}

// Conf for the routers package.
type Conf struct {
	GinMode         string `mapstructure:"gin_mode"`
	Addr            string `mapstructure:"addr"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

// Init create a gin.Engine and define multiplexer of the Engine.
func Init() {
	ValidateInstance = validator.New()

	log.Info("init ginrouters package...")
	gin.SetMode(Config.GinMode)

	// set up the routing.
	log.Info("instantiate gin engine...")
	Router = gin.New()
	log.Info("gin engine instantiate")

	// Middleware.
	log.Info("load middlewares...")
	Router.Use(ginzap.RecoveryWithZap(log, true))
	Router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	Router.Use(cors.Middleware(new(cors.Builder).New().WithOrigins("http://localhost:8080").Build()))
	log.Info("middlewares loaded")

	// Add your handler below this log.
	log.Info("load handlers...")

	Router.GET("/tasks", GetInstanceTaskRouter().GetAll)
	Router.Group("/task").
		POST("", GetInstanceTaskRouter().Post).
		GET("/:task_uuid", GetInstanceTaskRouter().Get).
		PUT("/:task_uuid", GetInstanceTaskRouter().Put).
		DELETE("/:task_uuid", GetInstanceTaskRouter().Delete)

	// Specific handler
	log.Info("load specific handlers...")
	Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, HomeMessage{
			Service:   "default",
			Copyright: "default",
		})
	})

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Not Found")
	})

	Router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, "Not Allowed")
	})
	log.Info("specific handlers loaded")

	log.Info("handlers loaded")

	log.Info("ginrouters package ready")
}
