package main

import (
	"fmt"
	"strconv"

	"github.com/fankys2012/gostudy/customerManager/model"
	"github.com/fankys2012/gostudy/customerManager/service"
)

type CustomerView struct {
	key             string
	loop            bool
	customerService *service.CustomerService
}

func (this *CustomerView) mainView() {
	for {
		fmt.Println("------------客户信息管理软件-----------\n")
		fmt.Println("           1、添 加 客 户\n")
		fmt.Println("           2、修 改 客 户\n")
		fmt.Println("           3、删 除 客 户\n")
		fmt.Println("           4、客 户 列 表\n")
		fmt.Println("           5、退   出\n")
		fmt.Println("请选择1-5")
		fmt.Scanln(&this.key)
		switch this.key {
		case "1":
			this.Add()
		case "2":
			this.Edit()
		case "3":
			this.Delete()
		case "4":
			this.list()
		case "5":
			fmt.Println("确认退出？y/n")
			var sure string
			for {
				fmt.Scanln(&sure)
				switch sure {
				case "y", "Y":
					this.loop = true
					goto outerquict
				case "n", "N":
					this.loop = false
					goto outerquict
				default:
					fmt.Println("输入有误请重新输入")
				}
			}
		default:
			fmt.Println("输入有误请重新输入")
		}
	outerquict:
		if this.loop {
			break
		}
	}
	fmt.Println("您已退出客服管理软件")
}

func (this *CustomerView) list() {
	customers := this.customerService.List()
	fmt.Println("--------- 客户列表--------")
	fmt.Println("编号\t姓名\t性别\t年龄\t电话    \t邮件")
	for _, c := range customers {
		fmt.Println(c.GetInfo())
	}
}

func (this *CustomerView) Add() {
	var name, gender, phone, email, age string
	fmt.Println("请输入姓名：")
	fmt.Scanln(&name)
	fmt.Println("请输入性别：")
	fmt.Scanln(&gender)
	fmt.Println("请输入年龄：")
	fmt.Scanln(&age)
	fmt.Println("请输入电话：")
	fmt.Scanln(&phone)
	fmt.Println("请输入邮件:")
	fmt.Scanln(&email)
	intAge, _ := strconv.Atoi(age)
	customer := model.NewCustomer(0, intAge, name, gender, phone, email)
	this.customerService.Add(customer)
}

func (this *CustomerView) Delete() {
	var id int
	fmt.Println("请输入删除客户ID：")
	fmt.Scanln(&id)
	if id == -1 {
		return
	}
	fmt.Println("确认删除Y/y?")
	confire := ""
	fmt.Scanln(&confire)
	if confire == "Y" || confire == "y" {
		bool := this.customerService.Delete(id)
		if bool {
			fmt.Println("删除成功")
		} else {
			fmt.Println("删除失败，ID不存在")
		}
	}
}

func (this *CustomerView) Edit() {
	fmt.Println("请选择待修改客户编号(-1退出)：")
	var id int
	fmt.Scanln(&id)
	if id == -1 {
		return
	}
	index := this.customerService.GetFindById(id)
	if index == -1 {
		fmt.Println("输入客户编号不存在！")
		return
	}
	customer := this.customerService.GetCustomerByIndex(index)
	var name, gender, phone, email string
	var age int
	fmt.Println("请输入姓名(" + customer.Name + ")：")
	fmt.Scanln(&name)
	fmt.Println("请输入性别(" + customer.Gender + ")：")
	fmt.Scanln(&gender)
	fmt.Println("请输入年龄(" + string(customer.Age) + ")：")
	fmt.Scanln(&age)
	fmt.Println("请输入电话(" + customer.Phone + ")：")
	fmt.Scanln(&phone)
	fmt.Println("请输入邮件(" + customer.Email + "):")
	fmt.Scanln(&email)
	if name != "" {
		customer.Name = name
	}
	if gender != "" {
		customer.Gender = gender
	}
	if phone != "" {
		customer.Phone = phone
	}
	if age > 0 {
		customer.Age = age
	}
	if email != "" {
		customer.Email = email
	}
	this.customerService.Edit(index, customer)

}

func main() {
	customerView := CustomerView{
		key:  "",
		loop: false,
	}
	customerView.customerService = service.NewCustomerService()
	customerView.mainView()

}
