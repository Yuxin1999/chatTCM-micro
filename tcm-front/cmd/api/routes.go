package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *Config) Router() *gin.Engine {
	r := gin.Default()
	// 定义跨域资源共享
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 加载HTML模板
	r.HTMLRender = MyRender()
	r.Static("/static", "/")
	// 默认页面
	r.GET("/", app.RenderLogin)
	r.GET("/sign", app.RenderSignup)

	// 注册
	r.POST("/signup", app.Signup)
	// 登录
	r.POST("/login", app.Login)

	// 需要验证的路由
	r.GET("/doctor/:token", app.RenderDiagnose)
	r.GET("/chat", app.Chat)

	return r
}

/* 前端渲染 */

func renderFiles(file string) []string {
	var slice []string
	slice = append(slice, file)
	baseSlice := []string{"templates/public/base.layout.html",
		"templates/public/header.partial.html",
		"templates/public/footer.partial.html"}
	slice = append(slice, baseSlice...)
	return slice
}

func MyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", renderFiles("templates/login.html")...)
	r.AddFromFiles("signup", renderFiles("templates/signup.html")...)
	r.AddFromFiles("diagnose", renderFiles("templates/diagnose.html")...)
	return r
}

// RenderSignup 渲染注册页面
func (app *Config) RenderSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{"title": "ChatTCM"})
}

// RenderLogin 渲染登录页面
func (app *Config) RenderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{"title": "ChatTCM"})
}

// RenderDiagnose 渲染诊断页面
func (app *Config) RenderDiagnose(c *gin.Context) {
	token := c.Param("token")
	app.Valid(token)
	c.SetCookie("userToken", token, 3600, "/", "localhost", false, false)
	c.HTML(http.StatusOK, "diagnose", gin.H{"title": "ChatTCM"})
}
