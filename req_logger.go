package main

import (
    "github.com/codegangsta/martini"
    "github.com/technoweenie/grohl"
    "time"
    "log"
    "net/http"
)

func ReqLogger() martini.Handler {
    return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
        start := time.Now()
        rw := res.(martini.ResponseWriter)
        c.Next()

        grohl.Log(grohl.Data{
            "method":   req.Method,
            "path":     req.URL.Path,
            "status":   rw.Status(),
            "duration": time.Since(start).Seconds(),
        })
    }
}
