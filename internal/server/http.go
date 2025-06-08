package server

import (
	accounterv1 "accounter_go/api/accounter/v1"
	v1 "accounter_go/api/helloworld/v1"
	"accounter_go/internal/conf"
	"accounter_go/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

// CORS中间件
func corsMiddleware() khttp.FilterFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 允许的源
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 允许的HTTP方法
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// 允许的请求头
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			// 允许携带认证信息
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// 预检请求的缓存时间
			w.Header().Set("Access-Control-Max-Age", "86400")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, accounter *service.AccounterService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
		khttp.Filter(corsMiddleware()), // 添加CORS中间件
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	accounterv1.RegisterAccounterHTTPServer(srv, accounter)
	
	return srv
}
