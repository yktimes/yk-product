package datamodels

// 标签

type Product struct {
	ID           int64  `json:"id" sql:"ID" yk:"id"`
	ProductName  string `json:"ProductName" sql:"productName" yk:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" yk:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" yk:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" yk:"ProductUrl"`
}
