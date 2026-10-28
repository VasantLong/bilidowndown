//go:build headless

package main

import "log"

// runApp 以无界面（headless）模式启动：不依赖托盘与图形库，仅提供 HTTP 服务。
func runApp() {
	onReady()
}

// setupTray headless 构建下没有托盘。
func setupTray() {}

// afterServerStart headless 构建下不自动打开浏览器。
func afterServerStart() {
	log.Printf("headless mode: open %s in a browser", urlLocalUnix)
}
