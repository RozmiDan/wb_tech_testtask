// @title           WB Orders Demo API
// @version         1.0
// @description     Демонстрационный сервис: получение информации о заказе.
// @BasePath        /
// @schemes         http
package main

import (
	"github.com/RozmiDan/wb_tech_testtask/internal/app"
	"github.com/RozmiDan/wb_tech_testtask/internal/config"
)

func main() {
	cnfg := config.MustLoad()

	app.Run(cnfg)
}
