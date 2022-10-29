package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Customer struct {
	Model
	Id                    int    `json:"id"`
	Name                  string `json:"name"`
	IsPrivate             *int   `json:"is_private"`
	Appid                 string `json:"appid"`
	Secret                string `json:"secret"`
	CreativeCallbackUrl   string `json:"creative_callback_url"`
	AdvertiserCallbackUrl string `json:"advertiser_callback_url"`
}
type CustomerModel struct {
}

// TableName sets the insert table name for this struct type
func (c *CustomerModel) TableName() string {
	return "customer"
}

func (c CustomerModel) GetCustomerTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&Customer{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c CustomerModel) ExistsById(id int) (bool, error) {
	var customer *Customer
	err := GetDB().First(&customer, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if customer.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c CustomerModel) ExistsByAppid(appid string) (bool, error) {
	var customer *Customer
	err := GetDB().Where("appid=?", appid).First(&customer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if customer.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c CustomerModel) ExistsByName(name string) (bool, error) {
	var customer *Customer
	err := GetDB().Where("name=?", name).First(&customer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if customer.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c CustomerModel) GetCustomerById(id int) (*Customer, error) {
	var customer *Customer
	err := GetDB().First(&customer, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return customer, nil
}

func (c CustomerModel) GetCustomerByAppid(appid string) (*Customer, error) {
	var customer *Customer
	err := GetDB().Where("appid=?", appid).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return customer, nil
}


func (c CustomerModel) GetAllCustomer() ([]*Customer, error) {
	var customer []*Customer
	err := GetDB().Find(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return customer, nil
}

func (c CustomerModel) GetCustomers(page int, perPage int, maps interface{}) ([]*Customer, error) {
	var customers []*Customer
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&customers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return customers, nil
}

func (c CustomerModel) CreateCustomer(customer *Customer) (*Customer, error) {
	result := GetDB().Create(&customer)
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}

func (c CustomerModel) UpdateCustomer(customer *Customer) (*Customer, error) {
	result := GetDB().Updates(customer)
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}

func (c CustomerModel) DeleteCustomer(customer *Customer) (*Customer, error) {
	GetDB().Delete(&customer)
	return customer, nil
}
