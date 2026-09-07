package main

import (
	"embed"
	"log"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// 互斥体
var mutexHandle syscall.Handle

func createMutex(name string) bool {
	utf16Name, _ := syscall.UTF16PtrFromString(name)
	handle, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("CreateMutexW").Call(
		0,
		0,
		uintptr(unsafe.Pointer(utf16Name)),
	)
	if handle == 0 {
		return false
	}
	err := syscall.GetLastError()
	if err == syscall.ERROR_ALREADY_EXISTS {
		syscall.CloseHandle(syscall.Handle(handle))
		return false
	}
	mutexHandle = syscall.Handle(handle)
	return true
}

func main() {
	// 互斥体防止多开
	if !createMutex("DaoscreenLock_Mutex_2026") {
		log.Println("程序已在运行中")
		return
	}
	defer syscall.CloseHandle(mutexHandle)

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "DaoscreenLock - 电脑锁屏工具",
		Width:  620,
		Height: 520,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 247, G: 248, B: 250, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
