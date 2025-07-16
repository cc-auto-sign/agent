// Package api 提供API服务相关功能
package api

import (
	"encoding/json"
	"net/http"
)

// handleAuthMiddleware 中间件：验证安全密钥
func (s *Server) handleAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求头中的安全密钥
		secureKey := r.Header.Get("X-Secure-Key")
		if secureKey == "" {
			// 如果请求头中没有安全密钥，尝试从表单或JSON正文中获取
			if r.Method == http.MethodPost {
				if err := r.ParseForm(); err == nil {
					secureKey = r.FormValue("secure_key")
				}

				if secureKey == "" && r.Header.Get("Content-Type") == "application/json" {
					var data map[string]interface{}
					decoder := json.NewDecoder(r.Body)
					if err := decoder.Decode(&data); err == nil {
						if key, ok := data["secure_key"].(string); ok {
							secureKey = key
						}
					}
					// 由于我们已经读取了body，需要重新设置它以便后续处理
					r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 限制请求体大小为1MB
				}
			}
		}

		// 验证安全密钥
		if secureKey != s.config.GetSecureKey() {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "无效的安全密钥",
			})
			return
		}

		next(w, r)
	}
}
