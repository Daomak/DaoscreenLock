package main

import (
	"context"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// 全屏锁屏
func (a *App) LockScreen() {
	runtime.WindowFullscreen(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	installKeyboardHook()
}

// 退出锁屏
func (a *App) UnlockScreen() {
	uninstallKeyboardHook()
	runtime.WindowUnfullscreen(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
	runtime.WindowCenter(a.ctx)
}

// ==================== 低级键盘钩子 ====================

var (
	hookHandle  uintptr
	hookRunning bool
)

// 虚拟键码
const (
	VK_LWIN    = 0x5B
	VK_RWIN    = 0x5C
	VK_TAB     = 0x09
	VK_ESCAPE  = 0x1B
	VK_F4      = 0x73
	VK_F5      = 0x74
	VK_F11     = 0x7A
	VK_CONTROL = 0x11
	VK_SHIFT   = 0x10
	VK_MENU    = 0x12
	VK_DELETE  = 0x2E
	VK_LSHIFT  = 0xA0
	VK_RSHIFT  = 0xA1
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	getMessage          = user32.NewProc("GetMessageW")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessage     = user32.NewProc("DispatchMessageW")
)

const WH_KEYBOARD_LL = 13
const HC_ACTION = 0

// 钩子回调函数
func keyboardHookProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == HC_ACTION {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		vk := kb.VkCode

		ctrlDown := getAsyncKeyState(VK_CONTROL)&0x8000 != 0
		altDown := getAsyncKeyState(VK_MENU)&0x8000 != 0
		shiftDown := getAsyncKeyState(VK_SHIFT)&0x8000 != 0 || getAsyncKeyState(VK_LSHIFT)&0x8000 != 0 || getAsyncKeyState(VK_RSHIFT)&0x8000 != 0

		// 拦截 Win 键
		if vk == VK_LWIN || vk == VK_RWIN {
			return 1
		}

		// 拦截 Alt+Tab
		if altDown && vk == VK_TAB {
			return 1
		}

		// 拦截 Alt+F4
		if altDown && vk == VK_F4 {
			return 1
		}

		// 拦截 Ctrl+Shift+Esc (任务管理器)
		if ctrlDown && shiftDown && vk == VK_ESCAPE {
			return 1
		}

		// 拦截 Ctrl+Esc (开始菜单)
		if ctrlDown && vk == VK_ESCAPE {
			return 1
		}

		// 拦截 F5
		if vk == VK_F5 {
			return 1
		}

		// 拦截 F11
		if vk == VK_F11 {
			return 1
		}

		// 拦截 Esc
		if vk == VK_ESCAPE {
			return 1
		}
	}

	ret, _, _ := callNextHookEx.Call(hookHandle, uintptr(nCode), wParam, lParam)
	return ret
}

func getAsyncKeyState(vk int) uint16 {
	ret, _, _ := user32.NewProc("GetAsyncKeyState").Call(uintptr(vk))
	return uint16(ret)
}

// 安装键盘钩子
func installKeyboardHook() {
	if hookRunning {
		return
	}
	hookRunning = true

	go func() {
		callback := syscall.NewCallback(keyboardHookProc)
		moduleHandle, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
		handle, _, _ := setWindowsHookEx.Call(
			WH_KEYBOARD_LL,
			callback,
			moduleHandle,
			0,
		)
		hookHandle = handle

		var msg struct {
			HWND   uintptr
			UINT   uint32
			WPARAM uintptr
			LPARAM uintptr
			Time   uint32
			Pt     struct{ X, Y int32 }
		}
		for hookRunning {
			ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if ret == 0 || ret == ^uintptr(0) {
				break
			}
			translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}()
}

// 卸载键盘钩子
func uninstallKeyboardHook() {
	if !hookRunning {
		return
	}
	hookRunning = false
	if hookHandle != 0 {
		unhookWindowsHookEx.Call(hookHandle)
		hookHandle = 0
	}
}

// 获取当前时间
func (a *App) GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
