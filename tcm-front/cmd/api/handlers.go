package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/rpc"
)

// Login 调用登录微服务
func (app *Config) Login(c *gin.Context) {
	// 解析表单，获取用户信息
	toAuth := AuthRequest{}
	toAuth.Name = c.PostForm("username")
	toAuth.Password = c.PostForm("password")
	// 验证登录
	var authed AuthResponse
	err := app.logAuth(toAuth, &authed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "登录失败",
			"data": err.Error(),
		})
		return
	}
	c.SetCookie("userToken", authed.Token, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, fmt.Sprintf("/doctor/%s", authed.Token))
}

// 登陆验证
func (app *Config) logAuth(toAuth AuthRequest, authed *AuthResponse) error {
	// 连接rpc客户端
	client, err := rpc.Dial("tcp", "tcm-authenticate:5001")
	if err != nil {
		return err
	}
	defer client.Close()
	// 调用验证微服务的登录验证
	err = client.Call("RPCServer.LogAuthenticate", toAuth, &authed)
	if err != nil {
		return err
	}
	return nil
}

// 调用验证
func (app *Config) validAuth(token string, authed *AuthResponse) error {
	// 连接rpc客户端
	client, err := rpc.Dial("tcp", "tcm-login:5001")
	if err != nil {
		return err
	}
	defer client.Close()
	// 调用验证微服务有效期验证函数
	err = client.Call("RPCServer.ValidAuthenticate", token, &authed)
	if err != nil {
		return err
	}
	return nil
}

// Signup 将注册信息发送到队列中
func (app *Config) Signup(c *gin.Context) {
	// 1.获取注册信息
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	info := SignRequest{
		username,
		email,
		password,
	}
	// 检查所有字段是否都不为空
	if username == "" || password == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请填写完整的信息",
		})
		return
	}

	// 2.发送信息到队列
	err := app.pushToQueue(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "注册失败，请稍候再试",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功，请输入验证码",
	})
	return
}

func (app *Config) Verify(c *gin.Context) {
	token := c.PostForm("verification-code")
	rpcClient, err := rpc.Dial("tcp", "tcm-diagnose:5001")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "验证失败，请重试",
		})
		return
	}
	var response *verifyResponse
	err = rpcClient.Call("Config.SignVerify", token, &response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "验证失败，请重试",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证成功",
	})
	c.Redirect(http.StatusFound, "/")
}

func (app *Config) Chat(c *gin.Context) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "websocket更新失败",
			"error": err.Error(),
		})
		return
	}
	// Handle the WebSocket connection
	app.chatWebsocket(conn)
}
