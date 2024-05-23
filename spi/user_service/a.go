package main

// 测试用
//go:generate go build --buildmode=plugin -o ./spi/user_service/a.so ./spi/user_service/a.go

type UserService struct{}

func (u UserService) Get() string {
	return "Get"
}

var UserSvc UserService
