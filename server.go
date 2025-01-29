package main

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()
	updateData(ctx)

	// 启动一个定时任务，每小时更新一次数据
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			updateData(ctx)
		}
	}()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "ok. data updated at " + lastUpdateTime,
		})
	})
	r.GET("/get_holidays", func(c *gin.Context) {

		mu.RLock()

		var list []Event
		typeParam := c.Query("type")
		switch typeParam {
		case "1", "":
			// 展示法定节假日以及调休信息
			list = getHolidays(data)
		case "2":
			// 展示24节气
			list = getJieQiList(data)
		}

		mu.RUnlock()

		// 根据年份进行筛选
		year := c.Query("year")
		if len(year) > 0 {
			if len(year) != 4 {
				c.JSON(400, gin.H{
					"code":    400,
					"message": "invalid year, year must be 4 digits",
				})
				return
			}
			// 动态判断是否最近三年
			yearInt, err := strconv.Atoi(year)
			if err != nil {
				c.JSON(400, gin.H{
					"code":    400,
					"message": "invalid year, year must be a number",
				})
				return
			}
			nowYear := time.Now().Year()
			if yearInt < nowYear-3 || yearInt > nowYear {
				c.JSON(400, gin.H{
					"code":    400,
					"message": "year out of range, only support past 3 years",
				})
				return
			}
			list = filterByYear(list, year)
		}

		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    list,
		})
	})
	r.Run(":8080")

}
