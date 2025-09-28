package ws

import (
	"time"
	"walletStatus/internal/infra/logger"
)

func (hdl *Client) startListener() {
	hdl.wg.Add(1)
	defer hdl.wg.Done()

	for {
		select {
		case <-hdl.ctx.Done():
			return
		default:

			if hdl.stopped.Load() {
				return
			}

			for !hdl.connAlive.Load() {
				if hdl.stopped.Load() {
					return
				}

				time.Sleep(500 * time.Millisecond)
			}

			_, msg, err := hdl.conn.ReadMessage()
			if err != nil {
				logger.Errorf(hdl.ctx, "failed to read message %w", err)
				hdl.connAlive.Store(false)

				continue
			}

			hdl.Out <- msg
		}
	}
}
