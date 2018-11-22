package http

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Logger is the default logger that will be used in this package.
var Logger = log.New(os.Stdout, "http: ", log.LstdFlags)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		logResponse := &logResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(logResponse, req)

		Logger.Printf(`%s [%s] "%s %s %s" %d %s "%s"`,
			remoteAddr(req),
			time.Since(start),
			req.Method,
			req.RequestURI,
			req.Proto,
			logResponse.statusCode,
			http.StatusText(logResponse.statusCode),
			req.UserAgent(),
		)
	})
}

func remoteAddr(req *http.Request) string {
	addr, _, _ := net.SplitHostPort(req.RemoteAddr)
	return addr
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
