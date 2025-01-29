package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoggerMiddleware 中间件：记录请求日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%d] %s %s (%v)", statusCode, c.Request.Method, path, latency)
	}
}

// ErrorHandler 中间件：错误处理
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "内部服务器错误",
			})
			return
		}
	}
}

func main() {
	// 设置 release 模式
	gin.SetMode(gin.ReleaseMode)

	ctx := context.Background()
	if err := updateData(ctx); err != nil {
		log.Fatalf("初始化数据失败: %v", err)
	}

	// 启动一个定时任务，每小时更新一次数据
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			if err := updateData(ctx); err != nil {
				log.Printf("更新数据失败: %v", err)
			}
		}
	}()

	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(ErrorHandler())

	// 路由组
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/holidays", getHolidaysList)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok. data updated at " + lastUpdateTime,
	})
}

func getHolidaysList(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	var list []Event
	typeParam := c.Query("type")
	switch typeParam {
	case "1", "":
		list = getHolidays(data)
	case "2":
		list = getJieQiList(data)
	default:
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid type parameter",
		})
		return
	}

	// 根据年份进行筛选
	if year := c.Query("year"); year != "" {
		if err := validateYear(year); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: err.Error(),
			})
			return
		}
		list = filterByYear(list, year)
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    list,
	})
}

func validateYear(year string) error {
	if len(year) != 4 {
		return fmt.Errorf("invalid year, year must be 4 digits")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("invalid year, year must be a number")
	}

	nowYear := time.Now().Year()
	if yearInt < nowYear-3 || yearInt > nowYear {
		return fmt.Errorf("year out of range, only support past 3 years")
	}

	return nil
}
