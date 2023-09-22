package Log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Verbose *log.Logger
	Warn    *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func Init(mode int) error {
	Verbose = log.New(io.Discard, "[VERBOSE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.Discard, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(io.Discard, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.Discard, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

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

	if mode < 1 {
		panic("Invaild logging mode. (1, 2, 3, 4)")
	}

	if mode >= 1 {
		Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
		Error.SetOutput(Writer)
	}

	if mode >= 2 {
		Warn = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
		Warn.SetOutput(Writer)
	}

	if mode >= 3 {
		Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		Info.SetOutput(Writer)
	}

	if mode >= 4 {
		Verbose = log.New(os.Stdout, "[VERBOSE] ", log.Ldate|log.Ltime|log.Lshortfile)
		Verbose.SetOutput(Writer)
	}

	if mode > 4 {
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
