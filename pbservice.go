package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
)

var (
	pbApp    *pocketbase.PocketBase
	pbCancel context.CancelFunc
)

func runPB() {
	// Configure data directory
	dataDir := filepath.Join(xdg.DataHome, "BoxTest", "pb_data")
	if os.Getenv("APP_ENV") == "development" {
		dataDir = "./pb_data"
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	pbCancel = cancel

	go func() {
		//pbApp := pocketbase.New()

		// Initialize PocketBase
		pbApp = pocketbase.NewWithConfig(pocketbase.Config{
			DefaultDataDir: dataDir,
		})
		pbApp.Bootstrap()
		if err := apis.Serve(pbApp, apis.ServeConfig{
			HttpAddr:        "127.0.0.1:8090",
			ShowStartBanner: false,
		}); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal(fmt.Sprintf("PocketBase failed: %v", err))
		}
	}()

	// Graceful shutdown listener
	go func() {
		<-ctx.Done()
		shutdownPB()
	}()
}

// Add this function to resolve pbCancel warning
func StopPB() {
	if pbCancel != nil {
		pbCancel() // Now properly used
	}
}

func shutdownPB() {
	if pbApp == nil {
		return
	}

	// Graceful server shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pbApp.ResetBootstrapState(); err != nil {
		appLogger.Error(fmt.Sprintf("Cleanup error: %v", err))
	}
}
