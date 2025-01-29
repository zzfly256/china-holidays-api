package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	ical "github.com/arran4/golang-ical"
)

type Event struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Remark    string `json:"remark"`
}

var (
	mu             sync.RWMutex
	data           []Event
	lastUpdateTime string

	JieQi = map[string]bool{
		"立春": true,
		"雨水": true,
		"惊蛰": true,
		"春分": true,
		"清明": true,
		"谷雨": true,
		"立夏": true,
		"小满": true,
		"芒种": true,
		"夏至": true,
		"小暑": true,
		"大暑": true,
		"立秋": true,
		"处暑": true,
		"白露": true,
		"秋分": true,
		"寒露": true,
		"霜降": true,
		"立冬": true,
		"小雪": true,
		"大雪": true,
		"冬至": true,
		"小寒": true,
		"大寒": true,
	}
)

// updateData 更新数据
func updateData(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	st := time.Now()
	eventList, err := parseCalendar(ctx, "https://calendars.icloud.com/holidays/cn_zh.ics/")
	if err != nil {
		log.Println("更新数据出错:", err.Error())
		return
	}

	data = eventList
	duration := time.Since(st)
	log.Println("更新数据成功，耗时:", duration.String())
	lastUpdateTime = time.Now().Format(time.DateTime)
	return
}

// parseCalendar 解析日历数据
func parseCalendar(ctx context.Context, icsUrl string) (result []Event, err error) {
	cal, err := ical.ParseCalendarFromUrl(icsUrl, ctx)
	if err != nil {
		return nil, fmt.Errorf("get data from remote error: %w", err)
	}

	for _, event := range cal.Events() {
		item := Event{}
		for _, property := range event.Properties {
			if property.BaseProperty.IANAToken == "SUMMARY" {
				item.Name = property.Value
			}
			if property.BaseProperty.IANAToken == "DTSTART" {
				item.StartDate = property.Value
			}
			if property.BaseProperty.IANAToken == "DTEND" {
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
		result = append(result, item)
	}

	return
}
