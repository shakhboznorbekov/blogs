package router

import (
	_ "github.com/blogs/internal/controller/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Auth interface {
	HasPermission(roles ...string) gin.HandlerFunc
}

type Authorization interface {
	SignIn(*gin.Context)
}

type User interface {
	AdminGetUserList(*gin.Context)
	AdminGetUserDetail(*gin.Context)
	AdminCreateUser(*gin.Context)
	AdminUpdateUser(*gin.Context)
	AdminDeleteUser(*gin.Context)
}

type Blog interface {
	AdminGetBlogList(*gin.Context)
	AdminGetBlogDetail(*gin.Context)
	AdminCreateBlog(*gin.Context)
	AdminUpdateBlog(*gin.Context)
	AdminDeleteBlog(*gin.Context)
}

type News interface {
	AdminGetNewList(ctx *gin.Context)
	AdminGetNewDetail(*gin.Context)
	AdminCreateNew(ctx *gin.Context)
	AdminUpdateNew(ctx *gin.Context)
	AdminDeleteNew(ctx *gin.Context)
}

type Router struct {
	auth          Auth
	user          User
	authorization Authorization
	blog          Blog
	new           News
}

func New(auth Auth, user User, authorization Authorization, blog Blog, new News) *Router {
	return &Router{auth: auth,
		user:          user,
		authorization: authorization,
		blog:          blog,
		new:           new,
	}
}

// Init ...
// @title API
// @version 1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name Shaxboz
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func (r *Router) Init(port string) error {
	router := gin.Default()

	// gin engine
	router.Use(customCORSMiddleware())

	// auth
	router.POST("/api/v1/user/sign-in", r.authorization.SignIn)

	//user
	router.GET("/api/v1/admin/user/list", r.auth.HasPermission("Admin"), r.user.AdminGetUserList)
	router.GET("/api/v1/admin/user/:id", r.auth.HasPermission("Admin"), r.user.AdminGetUserDetail)
	router.POST("/api/v1/admin/user/create", r.auth.HasPermission("Admin"), r.user.AdminCreateUser)
	router.PUT("/api/v1/admin/user/:id", r.auth.HasPermission("Admin"), r.user.AdminUpdateUser)
	router.DELETE("/api/v1/admin/user/:id", r.auth.HasPermission("Admin"), r.user.AdminDeleteUser)

	//blog
	router.GET("/api/v1/admin/blog/list", r.auth.HasPermission("Admin"), r.blog.AdminGetBlogList)
	router.GET("/api/v1/admin/blog/:id", r.auth.HasPermission("Admin"), r.blog.AdminGetBlogDetail)
	router.POST("/api/v1/admin/blog/create", r.auth.HasPermission("Admin"), r.blog.AdminCreateBlog)
	router.PUT("/api/v1/admin/blog/:id", r.auth.HasPermission("Admin"), r.blog.AdminUpdateBlog)
	router.DELETE("/api/v1/admin/blog/:id", r.auth.HasPermission("Admin"), r.blog.AdminDeleteBlog)

	//news
	router.GET("/api/v1/admin/news/list", r.auth.HasPermission("Admin"), r.new.AdminGetNewList)
	router.GET("/api/v1/admin/news/:id", r.auth.HasPermission("Admin"), r.new.AdminGetNewDetail)
	router.POST("/api/v1/admin/news/create", r.auth.HasPermission("Admin"), r.new.AdminCreateNew)
	router.PUT("/api/v1/admin/news/:id", r.auth.HasPermission("Admin"), r.new.AdminUpdateNew)
	router.DELETE("/api/v1/admin/news/:id", r.auth.HasPermission("Admin"), r.new.AdminDeleteNew)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router.Run(port)
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)

			return
		}

		c.Next()
	}
}
