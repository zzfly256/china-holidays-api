# China Holidays API

[![Go Version](https://img.shields.io/badge/Go-1.23.4-blue.svg)](https://golang.org/)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-green.svg)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

中国节假日 API 服务，提供节假日、调休和二十四节气信息查询。

## 功能特性

- 支持查询法定节假日及调休安排
- 支持查询二十四节气
- 支持按年份筛选（支持近三年数据）
- 自动更新数据
- 轻量级，易部署

## API 文档

### 1. 健康检查

```
GET /health
```

响应示例：
```json
{
    "code": 200,
    "message": "ok. data updated at [timestamp]"
}
```

### 2. 获取节假日信息

```
GET /get_holidays
```

查询参数：
- `type`: 数据类型（可选）
  - `1`或空: 法定节假日及调休信息（默认）
  - `2`: 二十四节气
- `year`: 年份（可选，仅支持近三年）

响应示例：
```json
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "name": "元旦",
            "startDate": "2024-01-01",
            "endDate": "2024-01-01",
            "remark": "节假日"
        }
        // ...
    ]
}
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/zzfly256/china-holidays-api.git
```

2. 安装依赖
```bash
go mod download
```

3. 运行服务
```bash
go run .
```

服务将在 `http://localhost:8080` 启动

## 部署

支持多种部署方式：

1. 直接运行
```bash
go build && ./china-holidays-api
```

2. Docker 部署（需自行构建镜像）
```bash
docker build -t china-holidays-api .
docker run -p 8080:8080 china-holidays-api
```

## 开源协议

本项目采用 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。
