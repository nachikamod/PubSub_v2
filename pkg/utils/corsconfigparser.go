package utils

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gorilla/websocket"
)

var conf = *ParseConfig()

// Parse config from env
// In case of development mode allow all origins to access else parse origins from env

func CorsConfigParser() cors.Config {

	if conf.Mode {
		return cors.Config{
			AllowAllOrigins: true,
		}
	}

	return cors.Config{
		AllowOrigins: conf.AllowedOrigins,
	}
}

func WsCorsUpgraderParser() websocket.Upgrader {
	if conf.Mode {
		return websocket.Upgrader{
			ReadBufferSize:  conf.WS_ReadBufferSize,
			WriteBufferSize: conf.WS_WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	}

	return websocket.Upgrader{
		ReadBufferSize:  conf.WS_ReadBufferSize,
		WriteBufferSize: conf.WS_ReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {

			var origin = r.Header.Get("origin")

			// Check if origin is in allowed list
			for _, allowOrigin := range conf.AllowedOrigins {
				if origin == allowOrigin {
					return true
				}
			}

			return false
		},
	}
}
