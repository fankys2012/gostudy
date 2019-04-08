package main

import "fmt"

var (
	userId  int
	userPwd string
)

func main() {
	//接受用户选择
	var key int
	//判断是否显示菜单
	var loop = true

	for loop {
		fmt.Println("-----------欢迎登陆-----------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出聊天室")
		fmt.Println("\t\t\t 请选择（1-3）")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			loop = false
		case 2:
			fmt.Println("注册用")
			loop = false
		case 3:
			fmt.Println("退出聊天室")
			loop = false
		default:
			fmt.Println("输入有误，请重新输入...")
		}
	}

	if key == 1 {
		fmt.Println("请输入用户ID")
		fmt.Scanf("%d\n", &userId)
		fmt.Println("请输入密码")
		fmt.Scanf("%s\n", &userPwd)

		err := login(userId, userPwd)
		if err == nil {
			fmt.Println("登陆成功")
		} else {
			fmt.Println("登陆失败")
		}
	}

}
