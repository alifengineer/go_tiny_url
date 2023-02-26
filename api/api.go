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
func SetUpRouter(h handlers.Handler, cfg config.Config) (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	docs.SwaggerInfo.Title = cfg.ServiceName
	docs.SwaggerInfo.Version = cfg.Version
	// docs.SwaggerInfo.Host = cfg.ServiceHost + cfg.HTTPPort
	docs.SwaggerInfo.Schemes = []string{cfg.HTTPScheme}

	r.Use(customCORSMiddleware())

	// CLIENT SERVICE
	r.POST("/client-platform", h.CreateClientPlatform)
	r.GET("/client-platform", h.GetClientPlatformList)
	r.GET("/client-platform/:client-platform-id", h.GetClientPlatformByID)
	r.GET("/client-platform-detailed/:client-platform-id", h.GetClientPlatformByIDDetailed)
	r.PUT("/client-platform", h.UpdateClientPlatform)
	r.DELETE("/client-platform/:client-platform-id", h.DeleteClientPlatform)

	r.POST("/client-type", h.CreateClientType)
	r.GET("/client-type", h.GetClientTypeList)
	r.GET("/client-type/:client-type-id", h.GetClientTypeByID)
	r.PUT("/client-type", h.UpdateClientType)
	r.DELETE("/client-type/:client-type-id", h.DeleteClientType)

	r.POST("/client", h.AddClient)
	r.GET("/client/:project-id", h.GetClientMatrix)
	r.PUT("/client", h.UpdateClient)
	r.DELETE("/client", h.RemoveClient)

	r.POST("/relation", h.AddRelation)
	r.PUT("/relation", h.UpdateRelation)
	r.DELETE("/relation/:relation-id", h.RemoveRelation)

	r.POST("/user-info-field", h.AddUserInfoField)
	r.PUT("/user-info-field", h.UpdateUserInfoField)
	r.DELETE("/user-info-field/:user-info-field-id", h.RemoveUserInfoField)

	// PERMISSION SERVICE
	r.GET("/role/:role-id", h.GetRoleByID)
	r.GET("/role", h.GetRolesList)
	r.POST("/role", h.AddRole)
	r.PUT("/role", h.UpdateRole)
	r.DELETE("/role/:role-id", h.RemoveRole)

	r.POST("/permission", h.CreatePermission)
	r.GET("/permission", h.GetPermissionList)
	r.GET("/permission/:permission-id", h.GetPermissionByID)
	r.PUT("/permission", h.UpdatePermission)
	r.DELETE("/permission/:permission-id", h.DeletePermission)

	r.POST("/upsert-scope", h.UpsertScope)

	r.POST("/permission-scope", h.AddPermissionScope)
	r.DELETE("/permission-scope", h.RemovePermissionScope)

	r.POST("/role-permission", h.AddRolePermission)
	r.POST("/role-permission/many", h.AddRolePermissions)
	r.DELETE("/role-permission", h.RemoveRolePermission)

	r.POST("/user", h.CreateUser)
	r.GET("/user", h.GetUserList)
	r.GET("/user/:user-id", h.GetUserByID)
	r.PUT("/user", h.UpdateUser)
	r.DELETE("/user/:user-id", h.DeleteUser)
	r.PUT("/user/reset-password", h.ResetPassword)
	r.POST("/user/send-message", h.SendMessageToUserEmail)

	r.POST("/user-relation", h.AddUserRelation)
	r.DELETE("/user-relation", h.RemoveUserRelation)

	r.POST("/upsert-user-info/:user-id", h.UpsertUserInfo)

	r.POST("/login", h.Login)
	r.DELETE("/logout", h.Logout)
	r.PUT("/refresh", h.RefreshToken)
	r.POST("/has-acess", h.HasAccess)

	v1 := r.Group("/v1")
	{
		v1.POST("/short-url", h.CreateShortUrl)
		v1.GET("/short-url/:hash", h.GetShortUrl)
	}

	sigma := r.Group("/sigma")
	{
		sigma.GET("/:hash", h.HandleLonger)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
