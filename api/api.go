package api

import (
	"go_auth_api_gateway/api/docs"
	"go_auth_api_gateway/api/handlers"
	"go_auth_api_gateway/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetUpRouter godoc
// @description This is a api gateway
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func SetUpRouter(h handlers.Handler, cfg config.Config) (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	docs.SwaggerInfo.Title = cfg.ServiceName
	docs.SwaggerInfo.Version = cfg.Version
	// docs.SwaggerInfo.Host = cfg.ServiceHost + cfg.HTTPPort
	docs.SwaggerInfo.Schemes = []string{cfg.HTTPScheme}

	r.Use(customCORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register-user", h.RegisterUser)
	r.POST("/login-user", h.LoginUser)

	r.Use(h.AuthMiddleware) // this middleware will be applied to all routes
	{
		r.PUT("/user", h.UpdateUser)
		r.DELETE("/user", h.DeleteUser)
		r.PUT("/user/reset-password", h.ResetPassword)
		r.POST("/short-url", h.CreateShortUrl)
		r.PUT("/short-url", h.UpdateShortUrl)
		r.GET("/short-url/:hash", h.GetShortUrlData)
		r.GET("sigma/:hash", h.HandleLonger)
		r.PUT("/url-qrcode", h.UrlToQrcode)
	}

	// r.GET("/user", h.GetUserList)
	// r.GET("/user/:user-id", h.GetUserByID)

	return
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
