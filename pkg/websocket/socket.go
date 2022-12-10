package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"overcompute.io/pubsub/pkg/utils"
	"overcompute.io/pubsub/pkg/websocket/pool"
)

var (
	wsupgrader = utils.WsCorsUpgraderParser()
	Pool       = *pool.NewPool()
)

// InitWebsocket godoc
//
//	@Summary		Pubsub Service API Document
//	@Description	Endpoint to connect to the websocket
//	@Param			uid	path	string	true	"User id"
//	@Success		200
//	@Router			/ws/{uid} [get]
func InitWebsocket() {

	r := gin.Default()

	r.GET("/ws/:uid", func(ctx *gin.Context) {
		wshandler(ctx.Writer, ctx.Request, ctx.Param("uid"))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(utils.ParseConfig().WS_PORT); err != nil {

		log.Fatalln("Error initializing socket server : ", err)

	}

}

func wshandler(w http.ResponseWriter, r *http.Request, uid string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	if uid == "" {
		log.Print("Invalid uid")
		conn.Close()
		return
	}

	Pool.AddClient(uid, conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "close") {

				Pool.RemoveClient(conn)

				return
			}
			break
		}
	}
}
