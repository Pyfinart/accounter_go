#!/bin/bash

echo "🚀 启动个人记账系统..."

# 检查是否已编译
if [ ! -f "build/accounter" ]; then
    echo "📦 编译应用..."
    go build -o build/accounter ./cmd/accounter
fi

# 启动API服务器（后台运行）
echo "🔧 启动API服务器 (端口8000)..."
./build/accounter -conf ./configs/config.yaml &
API_PID=$!

# 等待API服务器启动
sleep 3

# 启动Web界面服务器（后台运行）
echo "🌐 启动Web界面服务器 (端口30000)..."
cd web && go run server.go &
WEB_PID=$!

echo ""
echo "✅ 系统启动完成！"
echo ""
echo "📊 Web界面: http://localhost:30000"
echo "🔌 API服务: http://localhost:8000"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待用户中断
trap "echo ''; echo '🛑 正在停止服务...'; kill $API_PID $WEB_PID 2>/dev/null; exit" INT

# 保持脚本运行
wait 