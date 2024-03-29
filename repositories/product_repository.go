package repositories

import (
	"database/sql"
	"strconv"
	"yk-product/common"
	"yk-product/datamodels"
)

// 1 开发对应的接口
// 2 实现接口

type IProduct interface {
	// 连接数据
	Conn() error
	Insert(product *datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(product *datamodels.Product) error
	SelectByKey(int64 int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// 类似构造函数
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table, db}
}

// 数据库连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}

		p.mysqlConn = mysql
	}

	if p.table == "" {
		p.table = "product"
	}

	return
}

// 插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {

	// 判断连接
	if err = p.Conn(); err != nil {
		return
	}

	// 准备sql
	sql := "INSERT product set productName=?,productNum=?,productImage=?,productUrl=?"

	stmt, err := p.mysqlConn.Prepare(sql)

	if err != nil {
		return 0, err
	}
	// 传入参数
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)

	if err != nil {
		return 0, err
	}

	productId, err = result.LastInsertId()

	return
}

// 商品删除
func (p *ProductManager) Delete(productId int64) bool {

	// 判断连接
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "delete from product where ID=?"

	stmt, err := p.mysqlConn.Prepare(sql)

	if err != nil {
		return false
	}
	// 传入参数
	_, err = stmt.Exec(productId)

	if err != nil {
		return false
	}

	return true
}

// 商品的更新
func (p *ProductManager) Update(product *datamodels.Product) error {

	// 判断连接
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "update product set  productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)

	stmt, err := p.mysqlConn.Prepare(sql)

	if err != nil {
		return err
	}
	// 传入参数
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)

	if err != nil {
		return err
	}

	return nil
}

//根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return

}

// 获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errorProduct error) {

	// 判断连接
	if err := p.Conn(); err != nil {
		return nil, err
	}

	sql := "select * from " + p.table

	row, err := p.mysqlConn.Query(sql)
	defer row.Close()

	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(row)

	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {

		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}

	return

}
