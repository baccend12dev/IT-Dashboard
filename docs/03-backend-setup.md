# Backend Setup - Go

## 1. Install Golang

**Minimal version:** Go 1.21+

## 2. Initialize Project

```bash
mkdir backend
cd backend
go mod init it-dashboard
```

## 3. Install Dependencies

### Gin

```bash
go get github.com/gin-gonic/gin
```

### GORM

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

## 4. Basic Server

**main.go**

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "server running",
        })
    })

    r.Run(":8080")
}
```