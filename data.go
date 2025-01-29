package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

const (
	// iCloud 假日日历地址
	calendarURL = "https://calendars.icloud.com/holidays/cn_zh.ics"
	// 请求超时时间
	requestTimeout = 10 * time.Second
)

// updateData 从 iCloud 更新节假日数据
func updateData(ctx context.Context) error {
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, calendarURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析 ICS 文件
	calendar, err := ics.ParseCalendar(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("解析日历数据失败: %w", err)
	}

	// 转换数据
	var events []Event
	for _, event := range calendar.Events() {
		item := Event{}
		for _, property := range event.Properties {
			if property.IANAToken == string(ics.PropertySummary) {
				item.Name = property.Value
			}
			if property.IANAToken == string(ics.PropertyDtstart) {
				item.StartDate = property.Value
			}
			if property.IANAToken == string(ics.PropertyDtend) {
				item.EndDate = property.Value
			}
			if property.BaseProperty.IANAToken == "X-APPLE-SPECIAL-DAY" {
				switch property.Value {
				case "ALTERNATE-WORKDAY":
					item.Remark = "补班"
				case "WORK-HOLIDAY":
					item.Remark = "假期"
				default:
					item.Remark = property.Value
				}
			}
		}
		events = append(events, item)
	}

	// 使用互斥锁保护数据更新
	mu.Lock()
	data = events
	lastUpdateTime = time.Now().Format("2006-01-02 15:04:05")
	mu.Unlock()

	return nil
}
