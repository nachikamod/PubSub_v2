package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"overcompute.io/pubsub/pkg/config"
	"overcompute.io/pubsub/pkg/websocket/pool"
)

var (
	//wsupgrader = config.WsCorsUpgraderParser()
	Pool = *pool.NewPool()
)

type WSConf struct {
	conf *config.Config
}

func NewServer(conf *config.Config) *WSConf {
	return &WSConf{
		conf: conf,
	}
}

func (wsconf *WSConf) InitWebsocket() {

	r := gin.Default()

	r.GET("/ws/:uid", func(ctx *gin.Context) {
		wsconf.wshandler(ctx.Writer, ctx.Request, ctx.Param("uid"))
	})

	if err := r.Run(wsconf.conf.WS_PORT); err != nil {

		panic(err)

	}

}

// Websocket handler
func (wsconf *WSConf) wshandler(w http.ResponseWriter, r *http.Request, uid string) {
	upgrader := wsconf.conf.WsCorsUpgraderParser()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	// Check if uid is null/invalid if true then close the connection
	if uid == "" {
		log.Print("Invalid uid")
		conn.Close()
		return
	}

	// Add client to the connection pool
	Pool.AddClient(uid, conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {

				Pool.RemoveClient(conn)

				return
			}
			break
		}
	}
}
