package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 设置静态文件目录
	staticDir := "static"

	// 先注册API路由（更具体的路径）
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 代理到后端API服务器
		http.Error(w, "API服务请求，请确保后端服务在8000端口运行", http.StatusBadGateway)
	})

	// 然后注册静态文件服务器（通用路径）
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	port := ":30000"
	fmt.Printf("🌐 Web界面服务器启动在端口%s\n", port)
	fmt.Printf("📍 请访问: http://localhost%s\n", port)
	fmt.Printf("🔧 API服务请确保在8000端口运行\n")
	fmt.Println("📁 静态文件目录:", staticDir)

	log.Fatal(http.ListenAndServe(port, nil))
}
