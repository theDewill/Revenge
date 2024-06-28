package main

import (
	"ssego/messages"
	"time"

	"github.com/labstack/echo/v4"
)

func initSSE(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Example: Emit an event every 5 seconds
			sendMsg := messages.TempMessage{
				Msg:     "Hello, this is a periodic message!",
				MsgType: "periodic",
			}

			if err := sendMsg.Emit(c); err != nil {
				return err
			}
		case <-c.Request().Context().Done():
			return nil
		}

	}
}

func main() {
	e := echo.New()
	e.GET("/startSSE", initSSE)
	e.Logger.Fatal(e.Start(":8080"))
}
