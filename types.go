package main

import "sync"

// Event 表示一个节日或节气事件
type Event struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Remark    string `json:"remark,omitempty"`
}

// 全局变量
var (
	data           []Event      // 节假日数据
	mu             sync.RWMutex // 用于保护数据访问的互斥锁
	lastUpdateTime string       // 最后更新时间
)

// JieQi 二十四节气映射表
var JieQi = map[string]bool{
	"立春": true, "雨水": true, "惊蛰": true, "春分": true,
	"清明": true, "谷雨": true, "立夏": true, "小满": true,
	"芒种": true, "夏至": true, "小暑": true, "大暑": true,
	"立秋": true, "处暑": true, "白露": true, "秋分": true,
	"寒露": true, "霜降": true, "立冬": true, "小雪": true,
	"大雪": true, "冬至": true, "小寒": true, "大寒": true,
}
