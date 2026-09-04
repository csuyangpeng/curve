# Changelog

本文件记录 core 部署包的版本变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

---

## [2.7.7] - 2026-07-02

**发布环境：** `10.18.11.52`  
**发布时间：** 2026-07-02 10:29 ~ 10:32

### Added

- 发布包 `core_v2.7.7.tar.gz`，包含完整 `deploy` 部署目录。

### Changed

- 更新以下核心网二进制组件（来源：`../../patch/artifacts/binary/`）：
  - `n2proc`
  - `amf`
  - `smf`
  - `upfcp`
  - `upfdp`

### Removed

### 发布步骤

| 时间 | 操作 |
|------|------|
| 10:29:25 | 复制二进制：`n2proc`, `amf`, `smf`, `upfcp`, `upfdp` |

### 产物

- `core_v2.7.7.tar.gz`
