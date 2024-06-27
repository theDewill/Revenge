//here is the all message event structs with templates

package messages

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEprotocol interface {
	Emit(c echo.Context) error
}

type UpdateMessage struct {
	Msg   string
	Type  string
	Data  interface{}
	Event string
}

func (su *UpdateMessage) Emit(c echo.Context) error {
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
