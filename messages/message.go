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

type UpdateMessage struct { // Implementations - [ SSEprotocol ]
	Msg   string
	Type  string
	Data  interface{}
	Event string
}

type TempMessage struct {
	MsgType string
	Msg     string
}

func (tm *TempMessage) Emit(c echo.Context) error {
	c.Response().WriteHeader(http.StatusOK)

	message := fmt.Sprintf("type: %s\n Msg: %s\n", tm.MsgType, tm.Msg)
	if _, err := c.Response().Write([]byte(message)); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}

func (su *UpdateMessage) Emit(c echo.Context) error {

	c.Response().WriteHeader(http.StatusOK)

	message := fmt.Sprintf("event: %s\ndata: %s\n\n", su.Event, su.Msg)
	if _, err := c.Response().Write([]byte(message)); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}
