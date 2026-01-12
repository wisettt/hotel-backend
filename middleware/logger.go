package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        status := c.Writer.Status()
        path := c.Request.URL.Path
        method := c.Request.Method
        clientIP := c.ClientIP()
        // Simple log to stdout — replace with logger as needed
        println(method, path, clientIP, status, latency.String())
        _ = status
    }
}
