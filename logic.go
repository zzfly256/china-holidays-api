package main

import (
	"sort"
	"strings"
)

func getJieQiList(input []Event) (output []Event) {
	for _, datum := range input {
		if _, ok := JieQi[datum.Name]; ok {
			output = append(output, datum)
		}
	}

	// 根据日期进行排序
	sort.Slice(output, func(i, j int) bool {
		return output[i].StartDate < output[j].StartDate
	})

	return
}

func getHolidays(input []Event) (output []Event) {
	for _, datum := range input {
		if datum.Remark != "" {
			nameAry := strings.Split(datum.Name, "（")
			if len(nameAry) > 1 {
				datum.Name = nameAry[0]
			}
			output = append(output, datum)
		}
	}

	// 根据日期进行排序
	sort.Slice(output, func(i, j int) bool {
		return output[i].StartDate < output[j].StartDate
	})

	return
}

func filterByYear(input []Event, year string) (output []Event) {
	for _, datum := range input {
		if strings.HasPrefix(datum.StartDate, year) || strings.HasPrefix(datum.EndDate, year) {
			output = append(output, datum)
		}
	}

	return
}
