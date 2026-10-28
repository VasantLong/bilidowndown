# headless 构建（本 fork 的改动）

> 上游是 [iuroc/bilidown](https://github.com/iuroc/bilidown)。第三方 fork `piaodbg/bilidowndown` 只作 `upstream`（只读）参照；本 fork 的改动只进**自己的远端**：`origin` = `git@github.com-personal:VasantLong/bilidowndown.git`（GitHub 公开 fork，2026-10-28 建），`gitea` = `gitea:vasant/bilidowndown.git`（私有镜像，同一份提交）。两条远端都含 `main` 与 `feat/headless-mode`；不要往 `upstream` 推。
> 相对上游只加了一件事：**可无托盘、无 CGO 地作为后台 HTTP 服务运行**。上游 README 描述的是托盘玩法，别混用。

## 为什么需要

上游 `server/main.go` 直接 `systray.Run(...)`，而 `getlantern/systray` 在 Linux 编译期需要
pkg-config 的 `ayatana-appindicator3-0.1`（`libayatana-appindicator3-dev`，只能 `sudo apt` 装）。
在权限受限的服务器/容器里装不了 → 编译直接失败，也没法当守护进程跑。

## 实现（commit `9a40269`，分支 `feat/headless-mode`）

用 Go build tag 把「托盘 / 无托盘」拆成两个文件，`main.go` 只留三个钩子：

| 文件                      | build tag           | 内容                                                                                                              |
| ------------------------- | ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `server/tray_gui.go`      | `!headless`（默认） | `runApp() = systray.Run(onReady, nil)`、`setupTray()`（图标/标题/菜单）、`afterServerStart()`（打开浏览器）       |
| `server/tray_headless.go` | `headless`          | `runApp() = onReady()`，`setupTray()` 空实现，`afterServerStart()` 只打日志（不自动开浏览器）                     |
| `server/main.go`          | —                   | `main()` → `checkFFmpeg()` + `runApp()`；`onReady()` → `setupTray()` + 建表 + 起 HTTP 服务 + `afterServerStart()` |

钩子名（`runApp` / `setupTray` / `afterServerStart`）是**契约**：两个 tag 文件必须同时实现，否则编译报 undefined。
`systray` 的 import 已从 `main.go` 拿掉，只在 `tray_gui.go` 里。

## 构建

```bash
cd server
CGO_ENABLED=0 go build -tags headless -o bilidown     # 产物无 CGO、无 X11/GTK 依赖
```

- SQLite 走 `modernc.org/sqlite`（纯 Go），关掉 CGO 不影响数据层
- 不加 `-tags headless` 仍可照常编译出托盘版本（需装上面那个 dev 包）

## 运行

- **所有路径相对 cwd**（`static/` 前端、`data.db`、`download/` 默认下载目录）——必须在正确的工作目录里启动
- 监听 `*:8098`（`server/main.go` 的 `HTTP_PORT`），**HTTP 无鉴权**；对外暴露请自己加反向代理/鉴权
- headless 下不会自动打开浏览器，日志打印一次入口 URL
- `/api/quit` 是 `os.Exit(0)` 干净退出；systemd 的 `Restart=on-failure` 不会把它拉起来

## 部署（本机实况）

在 minihub 上以 user service 常驻：`~/apps/bilidown`（二进制 + `static/` + `rebuild.sh`）+
`~/.config/systemd/user/bilidown.service`，端口 8098 的 reachability 由 tailnet ACL 按端口收口。
完整 runbook（重建、补丁回放、排障）见 pi-sync KB `knowledge-base/infrastructure/minihub-服务器.md`「Bilidown」。
