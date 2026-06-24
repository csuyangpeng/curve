# Stock Quote Monitor

Full-stack implementation based on `aa.py`: fetches real-time quotes from the Tencent Finance API and displays pinyin name, current price, change, change %, low/high, and volume. Auto-refreshes every 2 seconds.

## Getting Started

```bash
go mod tidy
go run ./backend
```

Open http://localhost:8000 in your browser.

## API

| Endpoint | Description |
|----------|-------------|
| `GET /api/stocks` | Fetch real-time quotes for all monitored stocks |
| `GET /api/codes` | List monitored stock codes |

## Project Layout

```
curve/
├── aa.py              # Original script
├── go.mod
├── backend/
│   ├── main.go        # HTTP server
│   └── stock.go       # Quote fetching and parsing
└── frontend/
    ├── index.html
    ├── style.css
    └── app.js
```
# curve
