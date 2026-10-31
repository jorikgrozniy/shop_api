package main

import "shop_api/internal/app"

// @title Shop API
// @version 1.0
// @description REST API for shop service
// @BasePath /api/v1
func main() {
	app.New().Run()
}
