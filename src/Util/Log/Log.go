package Log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// Initalize the environment variables from .env file
func Init() {
	mode := os.Getenv("LOG_MODE")
	if mode == "" {
		mode = "production"
	}

	switch mode {
	case "production":
		initLog(PRODUCTION_MODE)
	case "development":
		initLog(DEVELOPMENT_MODE)
	default:
		initLog(PRODUCTION_MODE)
		Warn.Printf("Invalid logging mode: %s, set log level for production", mode)
	}
}

var (
	PRODUCTION_MODE  = 3
	DEVELOPMENT_MODE = 4
)

var (
	Verbose *log.Logger
	Warn    *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func initLog(level int) error {
	Verbose = log.New(io.Discard, "[DEBG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.Discard, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(io.Discard, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.Discard, "[ERR!] ", log.Ldate|log.Ltime|log.Lshortfile)

	date := time.Now().Format("2006-01-02_150405")
	pwd, err := os.Getwd()
	if err != nil {
		return err // os.Getwd() 예외처리
	}

	// 경로 포맷팅
	// FolderPath := fmt.Sprint(pwd, "/logs")
	// FilePath := fmt.Sprint(pwd, "/logs/", date, ".log")
	FolderPath := filepath.Join(pwd, "logs")
	FilePath := filepath.Join(pwd, "logs", fmt.Sprint(date, ".log"))

	os.MkdirAll(FolderPath, os.ModePerm)

	if !IsFileExists(FilePath) {
		_, err = os.Create(FilePath)
		if err != nil {
			return err // os.Create() 예외처리
		}
	}

	LogFile, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err // os.OpenFile() 예외처리
	}

	Writer := io.MultiWriter(LogFile, os.Stdout)

	if level < 1 {
		panic("Invaild logging mode. (1, 2, 3, 4)")
	}

	if level >= 1 {
		Error = log.New(os.Stderr, color.HiRedString("[ERR!] "), log.Ldate|log.Ltime)
		Error.SetOutput(Writer)
	}

	if level >= 2 {
		Warn = log.New(os.Stdout, color.YellowString("[WARN] "), log.Ldate|log.Ltime)
		Warn.SetOutput(Writer)
	}

	if level >= 3 {
		Info = log.New(os.Stdout, color.CyanString("[INFO] "), log.Ldate|log.Ltime)
		Info.SetOutput(Writer)
	}

	if level >= 4 {
		Verbose = log.New(os.Stdout, "[DEBG] ", log.Ldate|log.Ltime|log.Lshortfile)
		Info.Println("Verbose logging is enabled. Only use this mode for debugging.")
		Verbose.SetOutput(Writer)
	}

	if level > 4 {
		Warn.Println("Logging mode is greater than 4. Logging mode is set to the maximum level.")
	}

	return nil
}

func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		fmt.Println("Log System Initialization Error!")
		panic(err)
	}
}
