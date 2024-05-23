package main

// 测试用


type UserService struct{}

// GetName returns the name of the service
func (u UserService) Get() string {
	return "A"
}

// 导出对象
var UserSvc UserService


