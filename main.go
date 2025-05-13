package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS
var appLogger *CustomLogger

func init() {
	// Initialize logging
	var err error
	appLogger, err = startLogger()
	if err != nil {
		fmt.Printf("Failed to set up logging: %v\n", err)
		return
	}

	// Set PocketBase output to logger
	// SetLoggerWriter(appLogger.Writer())
}

func main() {

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:              "Box Test",
		Width:              1024,
		Height:             768,
		Logger:             appLogger,
		LogLevelProduction: logger.ERROR,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		appLogger.Error(err.Error())
	}
}
