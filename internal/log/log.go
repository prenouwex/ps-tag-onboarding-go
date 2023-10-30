package log

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wexinc/ps-tag-onboarding-go/internal/constants"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	Warn           *log.Logger
	Info           *log.Logger
	Error          *log.Logger
	fileWriter     io.Writer
	logMultiWriter io.Writer
)

func init() {
	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	filePath := filepath.Join(rootDir, constants.LogFile)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logMultiWriter = io.MultiWriter(os.Stdout, file)
	fileWriter = io.Writer(file)

	//log.SetOutput(logMultiWriter)
	Info = log.New(logMultiWriter, "[INFO]\t", log.LstdFlags|log.Lshortfile)
	Warn = log.New(logMultiWriter, "[WARN]\t", log.LstdFlags|log.Lshortfile)
	Error = log.New(logMultiWriter, "[ERROR]\t", log.LstdFlags|log.Lshortfile)
}

// RequestLogger middleware for printing info on terminal
func RequestLogger(next http.Handler) http.Handler {
	return getLogHttpHandler(next, os.Stdout, constants.NoColorLogs)
}

// RequestFileLogger middleware for log info to file
func RequestFileLogger(next http.Handler) http.Handler {
	return getLogHttpHandler(next, fileWriter, true)
}

// getLogHttpHandler method to get middleware
func getLogHttpHandler(next http.Handler, lw io.Writer, noColor bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := &middleware.DefaultLogFormatter{
			Logger:  log.New(lw, "[INFO]\t", log.LstdFlags),
			NoColor: noColor,
		}
		entry := f.NewLogEntry(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
		}()

		next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
	})
}
