package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SysUpdate struct {
	Msg   string
	Type  string
	Data  interface{}
	Event string
}

func (su *SysUpdate) Emit(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	message := fmt.Sprintf("event: %s\ndata: %s\n\n", su.Event, su.Msg)
	if _, err := c.Response().Write([]byte(message)); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}

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
			sendMsg := UpdateMessage{
				Msg:   "Hello, this is a periodic message!",
				Type:  "Notification",
				Data:  nil,
				Event: "message",
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
	e.GET("/sse", initSSE)
	e.Logger.Fatal(e.Start(":8080"))
}
