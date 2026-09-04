# Curve Quote Monitor

一个基于 Go 和原生前端实现的实时行情监控工具。项目从 `aa.py` 演进而来，后端从腾讯财经接口获取股票、指数和 ETF 行情，前端以表格形式展示价格、涨跌、涨跌幅、最高价、最低价和成交量，并支持按拼音、代码或名称搜索。

## 功能特性

- 实时获取腾讯财经行情数据
- 默认每 2 秒刷新一次
- 支持 A 股、港股、指数和 ETF 代码
- 自动将中文名称转换为拼音，便于搜索和展示
- 提供静态前端页面，无需额外前端构建步骤
- 提供 JSON API，方便后续扩展或接入其他客户端

## 技术栈

- Go 1.25+
- 标准库 `net/http`
- `github.com/mozillazg/go-pinyin`
- 原生 HTML、CSS、JavaScript

## 快速开始

在项目根目录执行：

```bash
go mod tidy
go run ./backend
```

服务启动后访问：

```text
http://localhost:8000
```

页面默认处于空闲状态，点击右上角按钮后开始拉取行情并自动刷新；浏览器标签页切到后台时会停止轮询。

## API

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/api/stocks` | 获取当前监控列表的实时行情 |
| `GET` | `/api/codes` | 获取当前监控的股票代码列表 |

`/api/stocks` 返回示例：

```json
{
  "updated_at": "2026-07-08 15:04:05",
  "count": 1,
  "stocks": [
    {
      "name": "East Money Information",
      "pinyin": "east money information",
      "code": "300059",
      "price": 20.12,
      "open": 19.8,
      "change": 0.32,
      "ratio": 1.62,
      "high": "20.35",
      "low": "19.72",
      "volume": 123456
    }
  ]
}
```

## 配置监控代码

监控列表定义在 `backend/stock.go` 的 `stockNameMapping` 中。新增代码时，按接口需要使用市场前缀：

- 上海：`sh`，例如 `sh601138`
- 深圳：`sz`，例如 `sz300059`
- 港股：`hk`，例如 `hk01347`

## 项目结构

```text
curve/
├── aa.py                 # 原始脚本
├── go.mod
├── go.sum
├── CHANGELOG.md
├── README.md
├── backend/
│   ├── main.go           # HTTP 服务、路由和静态文件服务
│   ├── stock.go          # 行情获取、解析和名称处理
│   └── cmd_test/         # 后端测试相关目录
└── frontend/
    ├── index.html        # 页面结构
    ├── style.css         # 页面样式
    └── app.js            # 轮询、搜索和渲染逻辑
```

## 常见问题

### 页面提示 Fetch failed

通常是后端无法访问腾讯财经接口，或本机网络不可用。可以先确认服务终端是否仍在运行，再检查网络连通性。

### 修改代码后页面没有变化

前端文件由后端静态服务直接读取。刷新浏览器即可看到最新内容；如仍未生效，可以尝试强制刷新。

### 如何调整刷新频率

修改 `frontend/app.js` 中的 `REFRESH_MS` 常量即可，单位为毫秒。
