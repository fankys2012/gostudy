package service

import "github.com/fankys2012/gostudy/customerManager/model"

type CustomerService struct {
	customer []model.Customer
	//自增ID
	customerNum int
}

func NewCustomerService() *CustomerService {
	customerService := &CustomerService{}
	customerService.customerNum = 1
	customer := model.NewCustomer(1, 30, "张三", "男", "15308006688", "3344@qq.com")
	customerService.customer = append(customerService.customer, customer)
	return customerService
}

func (this *CustomerService) List() []model.Customer {
	return this.customer
}

func (this *CustomerService) Add(customer model.Customer) bool {
	this.customerNum++
	customer.Id = this.customerNum
	this.customer = append(this.customer, customer)
	return true
}

func (this *CustomerService) GetFindById(id int) int {

	for index, customer := range this.customer {
		if customer.Id == id {
			return index
		}
	}
	return -1
}

func (this *CustomerService) Delete(id int) bool {
	index := this.GetFindById(id)
	if index == -1 {
		return false
	}
	//删除切片中的元素
	this.customer = append(this.customer[:index], this.customer[index+1:]...)
	return true
}

func (this *CustomerService) Edit(index int, customer model.Customer) bool {
	this.customer[index] = customer
	return true
}

func (this *CustomerService) GetCustomerByIndex(index int) model.Customer {
	return this.customer[index]
}
