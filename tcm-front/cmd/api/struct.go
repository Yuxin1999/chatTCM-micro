package main

// SignRequest 注册信息：发送到RabbitMQ
type SignRequest struct {
	UserName string
	Email    string
	password string
}

// AuthRequest 登录验证需要的信息
type AuthRequest struct {
	Name     string
	Password string
}

// AuthResponse 验证成功获得的信息
type AuthResponse struct {
	AuthStatus bool
	Token      string
}

// SimilarResponse 相似词RPC调用响应
type SimilarResponse struct {
	similar string
}

type verifyResponse struct {
	status bool
	msg    string
}
