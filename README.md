# DaoscreenLock

电脑锁屏工具，一键全屏锁定电脑，防止他人操作。

## 功能特性

- **一键锁屏**：点击按钮，输入密码即可全屏锁定
- **密码保护**：锁屏后需输入密码才能解锁，密码仅本次有效，不保存
- **拦截快捷键**：锁屏状态下拦截 Win 键、Alt+Tab、Alt+F4、Ctrl+Shift+Esc（任务管理器）等快捷键
- **全屏锁定**：锁屏后窗口占满全屏，始终置顶
- **防止多开**：互斥体机制，程序只能运行一个实例
- **实时时间**：锁屏界面显示当前时间（精确到秒）、日期

## 系统要求

- Windows 10 / Windows 11（64位）
- WebView2 Runtime（Windows 11 已自带，Windows 10 可能需要安装）

## 使用方法

1. 运行 `DaoscreenLock.exe`
2. 点击「立即锁屏」按钮
3. 输入锁屏密码，点击「确定锁屏」
4. 解锁时输入密码，点击「解锁」按钮

## 编译方法

### 环境要求

- Go 1.22+
- Wails CLI v2.8.0
- Node.js（前端构建，本项目前端为纯 HTML，无需构建）

### 编译命令

```bash
# 设置环境变量
$env:GOPROXY="https://goproxy.cn,direct"
$env:CGO_ENABLED="0"

# 编译
wails build -platform windows/amd64
```

编译产物在 `build/bin/DaoscreenLock.exe`

## 项目结构

```
DaoscreenLock/
├── app.go              # 后端：锁屏/解锁/键盘钩子
├── main.go             # 入口：窗口配置/互斥体
├── go.mod              # Go 依赖
├── wails.json          # Wails 配置
├── appicon.png         # 图标（透明背景）
├── appicon_exe.png     # 图标（白色背景，用于 exe）
├── frontend/
│   └── dist/
│       ├── index.html  # 前端界面
│       └── appicon.png # 前端图标
└── build/
    └── windows/
        └── icon.ico    # 程序图标
```

## 技术栈

- **后端**：Go + Wails v2
- **前端**：原生 HTML / CSS / JavaScript
- **窗口**：WebView2
- **键盘钩子**：Windows 低级键盘钩子（WH_KEYBOARD_LL）

## 注意事项

- 锁屏密码仅保存在内存中，程序关闭后清除，不会写入磁盘
- Ctrl+Alt+Del 是 Windows 系统安全快捷键，无法被任何应用拦截
- 锁屏后窗口始终置顶，禁止切屏和任务管理器
- 程序窗口大小固定，不可调整

## 作者

- **作者**：Daomak
- **GitHub**：https://github.com/Daomak
- **Gitee**：https://gitee.com/daomak
- **邮箱**：daomak@qq.com

## 开源协议

MIT License
