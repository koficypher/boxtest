package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// CustomLogger implements the Wails Logger interface
type CustomLogger struct {
	logger *log.Logger
}

// NewCustomLogger initializes the custom logger
func NewCustomLogger(output io.Writer) *CustomLogger {
	return &CustomLogger{
		logger: log.New(output, "WAILS: ", log.LstdFlags|log.Lmicroseconds),
	}
}

// Print implements the Logger interface
func (cl *CustomLogger) Print(message string) {
	cl.logger.Print(message)
}

// Trace implements the Logger interface
func (cl *CustomLogger) Trace(message string) {
	cl.logger.Print("TRACE: " + message)
}

// Debug implements the Logger interface
func (cl *CustomLogger) Debug(message string) {
	cl.logger.Print("DEBUG: " + message)
}

// Info implements the Logger interface
func (cl *CustomLogger) Info(message string) {
	cl.logger.Print("INFO: " + message)
}

// Warning implements the Logger interface
func (cl *CustomLogger) Warning(message string) {
	cl.logger.Print("WARNING: " + message)
}

// Error implements the Logger interface
func (cl *CustomLogger) Error(message string) {
	cl.logger.Print("ERROR: " + message)
}

// Fatal implements the Logger interface
func (cl *CustomLogger) Fatal(message string) {
	cl.logger.Fatal("FATAL: " + message)
}

// Add this method to CustomLogger
func (cl *CustomLogger) Writer() io.Writer {
	return cl.logger.Writer()
}

// func startLoggerOld() (*CustomLogger, error) {
// 	// fmt.Printf("APP ENV value %s", os.Getenv("APP_ENV"))
// 	isDevelopment := os.Getenv("APP_ENV") == "development"

// 	var logFilePath string
// 	if isDevelopment {
// 		logFilePath = filepath.Join(".", "app.log")
// 	} else {
// 		logFilePath = filepath.Join(xdg.DataHome, "BoxTest", "logs", "app.log")
// 	}

// 	// Ensure the log directory exists
// 	if err := os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm); err != nil {
// 		return nil, err
// 	}

// 	// Open the log file
// 	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Multi-writer to log to both stdout and the log file
// 	multiWriter := io.MultiWriter(os.Stdout, logFile)
// 	return NewCustomLogger(multiWriter), nil
// }

func startLogger() (*CustomLogger, error) {
	isDevelopment := os.Getenv("APP_ENV") == "development"
	var logDir string

	// First create the logger with minimal configuration
	tempLogger := log.New(os.Stdout, "INIT: ", log.LstdFlags|log.Lmicroseconds)

	if isDevelopment {
		logDir = filepath.Join(".", "logs")
	} else {
		logDir = filepath.Join(xdg.DataHome, "BoxTest", "logs")
	}

	// Create log directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Fallback with warning
		logDir = os.TempDir()
		tempLogger.Printf("Failed to create log directory: %v", err)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create fallback log directory: %w", err)
		}
		tempLogger.Printf("Using fallback log directory: %s", logDir)
	}

	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	// Now create the actual logger
	var output io.Writer = logFile
	if isDevelopment {
		output = io.MultiWriter(os.Stdout, logFile)
	}

	finalLogger := NewCustomLogger(output)

	// Use the final logger for these messages
	finalLogger.Info(fmt.Sprintf("Logging initialized in %s mode", os.Getenv("APP_ENV")))
	finalLogger.Info(fmt.Sprintf("Log file: %s", logFilePath))

	return finalLogger, nil
}

// func ensureWritableDir(path string) error {
// 	// Create the directory if it doesn't exist
// 	if err := os.MkdirAll(path, os.ModePerm); err != nil {
// 		return fmt.Errorf("failed to create directory: %w", err)
// 	}

// 	// Check if the directory is writable
// 	testFile := filepath.Join(path, ".testfile")
// 	err := os.WriteFile(testFile, []byte(""), 0644)
// 	if err != nil {
// 		return fmt.Errorf("directory is not writable: %w", err)
// 	}
// 	os.Remove(testFile)

// 	return nil
// }
