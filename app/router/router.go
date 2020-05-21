package router

import (
	"bacancy/go-boiler-plate/app/config"
	"bacancy/go-boiler-plate/app/controllers/user"
	"bacancy/go-boiler-plate/app/middleware"
	"net/http"

	limit "github.com/aviddiviner/gin-limit"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func ConfigureRouter() {
	if config.GetConfig().ENV == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func CreateRouter() {
	router = gin.New()

	router.Use(gin.Logger())
	router.Use(nice.Recovery(recoveryHandler))
	router.Use(limit.MaxAllowed(10))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET,PUT,POST,DELETE"},
		AllowHeaders:    []string{"accept,x-access-token,content-type,authorization"},
	}))

	/* Routes */
	public := router.Group("/")
	{
		// Public Routes
		public.POST("/signup", user.Signup)
	}

	protected := router.Group("/", middleware.ValidateToken())
	{
		// Private Routes
		protected.PUT("/user/profile", user.EditProfile)
		protected.GET("/user/profile", user.GetUserProfile)
	}

	// For Admin based Routes
	// admin := router.Group("/admin", middleware.ValidateAdminToken())
	// {
	// }
}

func RunRouter() {
	router.Run(":" + config.GetConfig().PORT)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	detail := ""
	if config.GetConfig().ENV == "develop" {
		detail = err.(error).Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": detail})
}
