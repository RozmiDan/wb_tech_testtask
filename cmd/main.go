package main

import (
	"github.com/RozmiDan/wb_tech_testtask/internal/app"
	"github.com/RozmiDan/wb_tech_testtask/internal/config"
)

func main() {
	cnfg := config.MustLoad()

	app.Run(cnfg)
}
