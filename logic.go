package main

import (
	"sort"
	"strings"
)

// getJieQiList 获取节气列表
func getJieQiList(input []Event) []Event {
	if len(input) == 0 {
		return nil
	}

	output := make([]Event, 0, len(input)/4) // 预分配大约1/4的容量，因为节气通常占总数据的1/4左右
	for _, datum := range input {
		if _, ok := JieQi[datum.Name]; ok {
			output = append(output, datum)
		}
	}

	if len(output) == 0 {
		return nil
	}

	// 根据日期进行排序
	sort.Slice(output, func(i, j int) bool {
		return output[i].StartDate < output[j].StartDate
	})

	return output
}

// getHolidays 获取节假日列表
func getHolidays(input []Event) []Event {
	if len(input) == 0 {
		return nil
	}

	output := make([]Event, 0, len(input)/2) // 预分配大约1/2的容量
	for _, datum := range input {
		if datum.Remark != "" {
			event := datum
			if idx := strings.Index(event.Name, "（"); idx != -1 {
				event.Name = event.Name[:idx]
			}
			output = append(output, event)
		}
	}

	if len(output) == 0 {
		return nil
	}

	// 根据日期进行排序
	sort.Slice(output, func(i, j int) bool {
		return output[i].StartDate < output[j].StartDate
	})

	return output
}

// filterByYear 按年份过滤事件列表
func filterByYear(input []Event, year string) []Event {
	if len(input) == 0 || len(year) != 4 {
		return nil
	}

	output := make([]Event, 0, len(input)/3) // 预分配大约1/3的容量，假设每年数据大约占总数据的1/3
	for _, datum := range input {
		if strings.HasPrefix(datum.StartDate, year) || strings.HasPrefix(datum.EndDate, year) {
			output = append(output, datum)
		}
	}

	return output
}
