package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Valid 局部中间件:用于验证token的有效期
func (app *Config) Valid(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token == "" {
			c.Redirect(http.StatusFound, "/")
		}
		// 验证登录
		var authed AuthResponse
		err := app.validAuth(token, &authed)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
		}
		c.Next()
	}
}

//
