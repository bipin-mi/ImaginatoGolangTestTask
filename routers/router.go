package routers

import (
	"ImaginatoGolangTestTask/controllers"
	"ImaginatoGolangTestTask/services"
	"ImaginatoGolangTestTask/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"ImaginatoGolangTestTask/shared/log"
	"ImaginatoGolangTestTask/shared/utils/middleware"
)

func ClientSecretAuth() gin.HandlerFunc {

	//context.WithValue(context.Background(), c.CtxUserID, int64(0))
	return func(c *gin.Context) {

		//secretKey := c.Request.Header.Get("Client-Secret")
		//if secretKey != os.Getenv("CLIENT_SECRET") {
		//	c.JSON(403, gin.H{"message": "Your request is not authorized"})
		//	c.Abort()
		//	return
		//}
		c.Next()
	}
}

// IRoutes is
type IRoutes interface {
	Setup()
	Run()
	Close(ctx context.Context) error
}

// Routes is
type Routes struct {
	router     *gin.Engine
	server     *http.Server
	middleware middleware.IMiddleware
	adminCtrl  *controllers.AdminController
}

// NewRouter is
func NewRouter() IRoutes {
	validation := validator.NewValidatorService()
	adminServe := services.NewAdminService()
	adminMiddleware := middleware.NewMiddlewareService()
	adminCtrl := controllers.InitController(validation, adminServe)
	router := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("Port")),
		Handler: router,
	}

	return &Routes{
		router,
		server,
		adminMiddleware,
		adminCtrl,
	}
}
func (route *Routes) Run() {
	log.GetLog().Info("", "service listen on "+os.Getenv("Port"))

	err := route.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.GetLog().Fatal("", "listen: %s\n", err)
	}
}

func (route *Routes) Close(ctx context.Context) error {
	if route.server != nil {
		return route.server.Shutdown(ctx)
	}
	return nil
}

func (route *Routes) Setup() {
	route.setupCors()

	api := route.router.Group("/admin/")
	api.Use(ClientSecretAuth())
	api.POST("signup", route.adminCtrl.Create)
	api.POST("login", route.adminCtrl.Login)
	api.POST("forgot-password", route.adminCtrl.ForgotPassword)
	api.POST("reset-password", route.adminCtrl.ResetPassword)
	api.GET("verify-email/:reset_token", route.adminCtrl.VerifyEmail)

	//For API where Authorization required
	api.Use(route.middleware.AuthHandler())
	{
		api.DELETE("admin-delete/:id", route.adminCtrl.AdminDelete)
		api.GET("list", route.adminCtrl.List)
	}

	api.Use(ClientSecretAuth())
}

func (route *Routes) setupCors() {
	route.router.Use(cors.New(cors.Config{
		ExposeHeaders:   []string{"Data-Length"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		AllowAllOrigins: true,
		MaxAge:          12 * time.Hour,
	}))
}
