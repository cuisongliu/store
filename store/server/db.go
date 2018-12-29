package serve

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

//GetEngine is
func GetEngine() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:Fanux#123@/store?charset=utf8")
	if err != nil {
		fmt.Println(err)
		panic("start engine failed")
	}
}

//Product is
type Product struct {
	ProductName string

	ProductURL   string // http://sealyun.oss-cn-beijing.aliyuncs.com/c89602f7cb2a/kube1.13.1.tar.gz
	ProductPrice float64

	ProductDivide float64
}

//UserProduct is
type UserProduct struct {
	ID           string
	Login        string
	ProductName  string
	ProductPrice float64

	Referrer    string
	PayReferrer float64

	Status     string // [see,payed,unknow]
	ClickCount int
}

//UserPayeeAccount use for alipay
type UserPayeeAccount struct {
	Login        string
	PayeeAccount string
	Amount       float64 //user earned money
	Passwd       string  //
}

//CreateTables is
func CreateTables() {
	err := engine.CreateTables(new(User))
	err = engine.CreateTables(new(Product))
	err = engine.CreateTables(new(UserProduct))
	err = engine.CreateTables(new(UserPayeeAccount))

	err = engine.Sync(new(User))
	err = engine.Sync(new(Product))
	err = engine.Sync(new(UserProduct))
	err = engine.Sync(new(UserPayeeAccount))

	if err != nil {
		fmt.Println("new table failed", err)
	}
}

//Get is
func (upa *UserPayeeAccount) Get(login string) (bool, error) {
	return engine.Where("login = ?", login).Get(upa)
}

//Save is
func (upa *UserPayeeAccount) Save() (int64, error) {
	return engine.Insert(upa)
}

//Update is
func (upa *UserPayeeAccount) Update() (int64, error) {
	return engine.Where("login = ?", upa.Login).Update(upa)
}

//Save is
func (user *User) Save() (int64, error) {
	return engine.Insert(user)
}

//Get is
func (user *User) Get(login string) (bool, error) {
	return engine.Where("login = ?", login).Get(user)
}

//Save is
func (up *UserProduct) Save() (int64, error) {
	return engine.Insert(up)
}

//Update is
func (up *UserProduct) Update() (int64, error) {
	return engine.Where("login = ?", up.Login).And("product_name = ?", up.ProductName).Update(up)
}

//Get is
func (up *UserProduct) Get(login, product string) (bool, error) {
	return engine.Where("login = ?", login).And("product_name = ?", product).Get(up)
}

//Save is
func (p *Product) Save() (int64, error) {
	return engine.Insert(p)
}

//Get is
func (p *Product) Get(name string) (bool, error) {
	return engine.Where("product_name= ?", name).Get(p)
}

//GetProductURL is
func GetProductURL(name string) string {
	p := &Product{ProductName: name}
	p.Get(name)
	return p.ProductURL
}

//GetProductPrice is
func GetProductPrice(name string) float64 {
	p := &Product{ProductName: name}
	p.Get(name)
	return p.ProductPrice
}

//GetProductDevide is
func GetProductDevide(name string) float64 {
	p := &Product{ProductName: name}
	has, err := p.Get(name)
	if !has || err != nil {
		fmt.Println("get devide failed")
		return 0
	}
	return p.ProductPrice * p.ProductDivide
}

func init() {
	GetEngine()
	CreateTables()

	p := &Product{
		ProductName:   "kubernetes1.13.1",
		ProductURL:    "http://sealyun.oss-cn-beijing.aliyuncs.com/c89602f7cb2a/kube1.13.1.tar.gz",
		ProductPrice:  5,
		ProductDivide: 0.6,
	}
	_, err := p.Save()
	if err != nil {
		fmt.Println("save product failed")
	}
}
