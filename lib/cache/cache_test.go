package cache

import (
	"fmt"
	"testing"
)

//================原始数据操作========================

//demo model
type Prov struct {
	ID   int64
	Name string
}

// from mysql or other...
func (m *Prov) GetAllDataFromDB() (data []*Prov, err error) {
	data = make([]*Prov, 0)
	data = append(data, &Prov{
		ID:   51,
		Name: "四川",
	})
	data = append(data, &Prov{
		ID:   50,
		Name: "重庆",
	})
	data = append(data, &Prov{
		ID:   11,
		Name: "北京",
	})
	return
}

//===========================cache 接口实现=====================================

// ProvCache 省份信息缓存 一个ICache的实现
type ProvCache struct{}

func (m *ProvCache) KeyName() string {
	return "prov"
}

func (m *ProvCache) PrimaryKey(obj interface{}) string {
	return fmt.Sprintf("%v", obj.(*Prov).ID)
}

// GetAllData
func (m *ProvCache) GetAllData() (data interface{}, err error) {
	data, err = new(Prov).GetAllDataFromDB()
	return
}

//===================使用者===============

func TestCacheUtil_GetAll(t *testing.T) {
	nc := NewMemCache()
	data, err := nc.GetAll(new(ProvCache))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("===test all===")

	//返回值断言
	value, ok := data.([]*Prov)
	if ok {
		fmt.Printf("value: %#v \n", value)
		for _, item := range value {
			fmt.Println(item)
		}
	}

}

func TestCacheUtil_GetItem(t *testing.T) {
	nc := NewMemCache()
	data, err := nc.GetItem(new(ProvCache), 51)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("===test item===")
	value, ok := data.(*Prov)
	if ok {
		fmt.Println(value)
	}

}
