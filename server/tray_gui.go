//go:build !headless

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/getlantern/systray"
)

// runApp 以桌面托盘模式启动（默认构建）。
func runApp() {
	systray.Run(onReady, nil)
}

// setupTray 初始化托盘图标、标题与菜单。
func setupTray() {
	setIcon()
	setTitle()
	setMenuItem()
}

// afterServerStart 在 HTTP 服务启动后打开浏览器。
func afterServerStart() {
	openBrowser(urlLocalUnix)
}

// setIcon 设置托盘图标
func setIcon() {
	var path string
	if runtime.GOOS == "windows" {
		path = "static/favicon.ico"
	} else {
		path = "static/favicon-32x32.png"
	}
	systray.SetIcon(mustReadFile(path))
}

// mustReadFile 返回文件字节内容
func mustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("os.ReadFile:", err)
	}
	return data
}

// setTitle 设置托盘标题和工具提示
func setTitle() {
	title := "Bilidown"
	tooltip := fmt.Sprintf("%s 视频解析器 %s (port:%d)", title, VERSION, HTTP_PORT)
	// only available on Mac and Windows.
	systray.SetTooltip(tooltip)
}

// setMenuItem 设置托盘菜单
func setMenuItem() {
	openBrowserItemText := fmt.Sprintf("打开主界面 (port:%d)", HTTP_PORT)
	openBrowserItem := systray.AddMenuItem(openBrowserItemText, openBrowserItemText)
	go func() {
		for {
			<-openBrowserItem.ClickedCh
			openBrowser(urlLocalUnix)
		}
	}()

	aboutItemText := "Github 项目主页"
	aboutItem := systray.AddMenuItem(aboutItemText, aboutItemText)
	go func() {
		for {
			<-aboutItem.ClickedCh
			openBrowser("https://github.com/iuroc/bilidown")
		}
	}()

	exitItemText := "退出应用"
	exitItem := systray.AddMenuItem(exitItemText, exitItemText)
	go func() {
		for {
			<-exitItem.ClickedCh
			log.Printf("Bilidown has exited.")
			systray.Quit()
		}
	}()
}
