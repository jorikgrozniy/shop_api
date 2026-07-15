package main

import "shop_api/internal/app"

// @title Shop API
// @version 1.0
// @description REST API for shop service
// @@host shop.local
// @BasePath /api/v1
// @schemes https
func main() {
	app.New().Run()
}
