package ws

import (
	"time"
	"walletStatus/internal/infra/logger"

	"github.com/gorilla/websocket"
)

func (hdl *Client) startWriter() {
	hdl.wg.Add(1)
	defer hdl.wg.Done()

	for {
		select {
		case <-hdl.ctx.Done():
			return
		case payload := <-hdl.In:
			for !hdl.connAlive.Load() {
				time.Sleep(500 * time.Millisecond)
			}

			if err := hdl.conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
				logger.Errorf(hdl.ctx, "failed to send message", err)
			}
		}
	}
}
