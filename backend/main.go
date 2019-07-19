package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"yk-product/backend/web/controllers"
	"yk-product/common"
	"yk-product/repositories"
	"yk-product/services"
)

func main() {

	// 1 创建iris实例
	app := iris.New()

	// 2 设置错误模式,在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	// 3 注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout(
		"share/layout.html").Reload(true)

	app.RegisterView(tmplate)

	// 4 设置模板目录
	app.StaticWeb("/assets", "./backend/web/assets")

	// 出现异常跳转指定页面
	app.OnAnyErrorCode(func(context iris.Context) {
		context.ViewData("message", context.Values().GetStringDefault("message", "访问的页面出错"))
		context.ViewLayout("")
		_ = context.View("share/error.html")
	})

	// 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5 注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)

	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	// 6 启动服务
	_ = app.Run(iris.Addr("localhost:8088"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
