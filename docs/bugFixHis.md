# Bug修复历史记录

## 2025年 - ngrok跨域访问问题

### 问题描述
- **时间**: 2025年6月27日
- **现象**: 通过ngrok暴露的地址 `https://e749-111-243-106-98.ngrok-free.app/` 访问web页面时，页面显示"加载中..."而没有显示任何数据记录
- **对比**: 直接访问 `http://localhost:30000/` 时数据正常显示

### 问题原因分析

#### 1. API地址配置问题
- **位置**: `web/static/index.html` 第349行
- **问题代码**:
  
  ```javascript
  const API_BASE_URL = window.location.protocol + '//' + window.location.hostname + ':8000';
  ```
- **问题**: 硬编码API请求到8000端口，当通过ngrok访问时，浏览器会尝试从 `https://e749-111-243-106-98.ngrok-free.app:8000` 获取数据

#### 2. 跨域访问限制
- **根本原因**: 浏览器的同源策略阻止了从 `https://e749-111-243-106-98.ngrok-free.app` 向 `https://e749-111-243-106-98.ngrok-free.app:8000` 发送请求
- **端口差异**: ngrok只暴露了30000端口，但API请求指向8000端口

#### 3. Web服务器代理功能不完整
- **位置**: `web/server.go`
- **问题**: API路由处理函数只是返回错误信息，没有真正代理请求到后端服务
- **原代码**:
  
  ```go
  http.Error(w, "API服务请求，请确保后端服务在8000端口运行", http.StatusBadGateway)
  ```

### 解决方案

#### 1: 修改Web服务器实现反向代理

**修改文件**: `web/server.go`

**修改内容**:

1. 创建反向代理：
   ```go
   // 创建反向代理到后端API服务器
   backendURL, err := url.Parse("http://localhost:8000")
   if err != nil {
       log.Fatal("无法解析后端URL:", err)
   }
   
   // 创建反向代理
   proxy := httputil.NewSingleHostReverseProxy(backendURL)
   ```

2. 修改API路由处理：
   ```go
   // 使用反向代理转发请求到后端
   proxy.ServeHTTP(w, r)
   ```

#### 2: 修改前端API地址配置

**修改文件**: `web/static/index.html`

**修改内容**:
```javascript
// 修改前
const API_BASE_URL = window.location.protocol + '//' + window.location.hostname + ':8000';

// 修改后
const API_BASE_URL = window.location.protocol + '//' + window.location.host;
```

### 修复效果

修复后的架构：
```
用户浏览器 → ngrok → localhost:30000 (web服务器) → localhost:8000 (API服务器)
```

### 技术要点

1. **反向代理**: 使用Go的 `httputil.NewSingleHostReverseProxy` 实现API请求转发
2. **CORS处理**: 在web服务器层面设置跨域头，避免浏览器阻止请求
3. **端口统一**: 前端不再直接访问8000端口，而是通过30000端口的代理功能

### 相关文件

- `web/server.go` - Web服务器配置
- `web/static/index.html` - 前端页面和JavaScript代码
- `configs/config.yaml` - 后端API服务器配置（8000端口）

### 验证方法

1. 启动服务：
   ```bash
   ./start.sh
   ```

2. 启动ngrok：
   ```bash
   ngrok http 30000
   ```

3. 访问ngrok提供的地址，确认数据正常加载

### 注意事项

- 确保后端API服务器在8000端口正常运行
- ngrok只暴露30000端口即可，不需要暴露8000端口
- 如果修改了端口配置，需要同步更新相关文件

## 问题：Mac合盖后ngrok隧道断开，公网无法访问本地服务

### 现象
- 使用 `ngrok http 30000` 暴露本地web服务到互联网
- Mac合上盖子后，ngrok隧道变为offline，公网无法访问
- 打开盖子后，ngrok恢复正常

### 原因
Mac合盖后会进入睡眠，导致ngrok进程挂起或网络断开，隧道失效。

---

### 解决方案：使用launchd服务管理ngrok

### 具体操作步骤

1. **停止当前ngrok进程**
   ```bash
   pkill ngrok
   ```

2. **创建launchd配置文件**
   在 `~/Library/LaunchAgents/` 目录下新建 `com.ngrok.tunnel.plist` 文件，内容如下：
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
   <plist version="1.0">
   <dict>
       <key>Label</key>
       <string>com.ngrok.tunnel</string>
       <key>ProgramArguments</key>
       <array>
           <string>/opt/homebrew/bin/ngrok</string>
           <string>http</string>
           <string>30000</string>
       </array>
       <key>RunAtLoad</key>
       <true/>
       <key>KeepAlive</key>
       <true/>
       <key>StandardOutPath</key>
       <string>/tmp/ngrok.log</string>
       <key>StandardErrorPath</key>
       <string>/tmp/ngrok.error.log</string>
       <key>WorkingDirectory</key>
       <string>/Users/yhj19/developer/MY_Go/accounter_go</string>
   </dict>
   </plist>
   ```
   > 注意：ngrok路径和WorkingDirectory需根据实际情况修改。

3. **加载launchd服务**
   ```bash
   launchctl bootstrap gui/$UID ~/Library/LaunchAgents/com.ngrok.tunnel.plist
   ```

4. **验证服务是否运行**
   ```bash
   ps aux | grep ngrok
   tail -f /tmp/ngrok.log
   ```

5. **管理命令**
   - 停止服务：
     ```bash
     launchctl bootout gui/$UID ~/Library/LaunchAgents/com.ngrok.tunnel.plist
     ```
   - 重新加载：
     ```bash
     launchctl bootstrap gui/$UID ~/Library/LaunchAgents/com.ngrok.tunnel.plist
     ```

### 效果
- 合盖后ngrok依然在后台运行，公网访问不受影响。
- 用户登录时自动启动，无需手动干预。
