const REFRESH_MS = 2000;
const tbody = document.getElementById("stock-body");
const statusDot = document.getElementById("status-dot");
const statusText = document.getElementById("status-text");
const updatedAt = document.getElementById("updated-at");
const stockCount = document.getElementById("stock-count");
const searchInput = document.getElementById("search");
const toggleBtn = document.getElementById("toggle-btn");
const contentPanel = document.getElementById("content-panel");

let allStocks = [];
let timer = null;
let visible = false;

function formatNum(n) {
  return Number(n).toLocaleString("en-US");
}

function changeClass(change) {
  if (change > 0) return "up";
  if (change < 0) return "down";
  return "flat";
}

function formatChange(change) {
  const sign = change > 0 ? "+" : "";
  return `${sign}${change.toFixed(2)}`;
}

function formatRatio(ratio) {
  const sign = ratio > 0 ? "+" : "";
  return `${sign}${ratio.toFixed(2)}%`;
}

function renderStocks(stocks) {
  const sorted = [...stocks].sort((a, b) => {
    const na = a.name.toLowerCase();
    const nb = b.name.toLowerCase();
    return na < nb ? -1 : na > nb ? 1 : 0;
  });
  const keyword = searchInput.value.trim().toLowerCase();
  const filtered = keyword
    ? sorted.filter(
        (s) =>
          s.pinyin.toLowerCase().includes(keyword) ||
          s.code.toLowerCase().includes(keyword) ||
          s.name.toLowerCase().includes(keyword)
      )
    : sorted;

  stockCount.textContent = `${filtered.length} rows`;

  if (filtered.length === 0) {
    tbody.innerHTML = '<tr><td colspan="12" class="loading">No matching results</td></tr>';
    return;
  }

  tbody.innerHTML = filtered
    .map((s) => {
      const cls = changeClass(s.change);
      return `
        <tr>
          <td>${s.pinyin}</td>
          <td class="num">${s.code}</td>
          <td class="num ${cls}">${s.price}</td>
          <td class="num ${cls}">${formatChange(s.change)}</td>
          <td class="num ${cls}">${formatRatio(s.ratio)}</td>
          <td class="num">${s.lowest}</td>
          <td class="num">${s.highest}</td>
          <td class="num">${s.sell1}</td>
          <td class="num">${formatNum(s.sell1_vol)}</td>
          <td class="num">${s.buy1}</td>
          <td class="num">${formatNum(s.buy1_vol)}</td>
          <td class="num">${formatNum(s.volume)}</td>
        </tr>`;
    })
    .join("");
}

function setStatus(ok, text) {
  statusDot.className = "dot " + (ok ? "ok" : "err");
  statusText.textContent = text;
}

function stopPolling() {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
}

function hidePanel() {
  visible = false;
  contentPanel.classList.add("hidden");
  toggleBtn.classList.remove("active");
  stopPolling();
  setStatus(false, "Idle");
  updatedAt.textContent = "";
}

function showPanel() {
  visible = true;
  contentPanel.classList.remove("hidden");
  toggleBtn.classList.add("active");
  tbody.innerHTML = '<tr><td colspan="12" class="loading">Loading...</td></tr>';
  fetchStocks();
  stopPolling();
  timer = setInterval(fetchStocks, REFRESH_MS);
}

async function fetchStocks() {
  if (!visible) return;

  try {
    const res = await fetch("/api/stocks");
    if (!res.ok) throw new Error("Request failed");
    const data = await res.json();
    allStocks = data.stocks;
    renderStocks(allStocks);
    updatedAt.textContent = `Updated: ${data.updated_at}`;
    setStatus(true, "Connected");
  } catch {
    setStatus(false, "Fetch failed");
    if (allStocks.length === 0) {
      tbody.innerHTML = '<tr><td colspan="12" class="error">Unable to fetch quotes. Check your network or try again later.</td></tr>';
    }
  }
}

toggleBtn.addEventListener("click", () => {
  if (visible) return;
  showPanel();
});

searchInput.addEventListener("input", () => renderStocks(allStocks));

document.addEventListener("visibilitychange", () => {
  if (document.hidden) {
    hidePanel();
  }
});

hidePanel();
