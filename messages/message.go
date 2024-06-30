//here is the all message event structs with templates

package messages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEprotocol interface {
	Emit(c echo.Context) error
}

type ViewMessage struct {
	action       string `json:"action"`
	component    string `json:"component"`
	update_value string `json:"update_value"`
}

type UpdateMessage struct { // Implementations - [ SSEprotocol ]
	Msg   string      `json:"msg"`
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
	Event string      `json:"event"`
}

type TempMessage struct {
	Type string `json:"msgType"`
	Msg  string `json:"msg"`
}

func (tm *TempMessage) Emit(c echo.Context) error {

	c.Response().WriteHeader(http.StatusOK)

	msgJson, err := json.Marshal(tm)
	if err != nil {
		fmt.Println("Error in marshalling json")
	}
	message := fmt.Sprintf("data: {\"jsonContent\": %s}\n\nevent: %s\n\n", string(msgJson), tm.Type)
	if _, err := c.Response().Write([]byte(message)); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}

func (su *UpdateMessage) Emit(c echo.Context) error {

	c.Response().WriteHeader(http.StatusOK)

	msgJson, err := json.Marshal(su)
	if err != nil {
		fmt.Println("Error in marshalling json")
	}

	message := fmt.Sprintf("data: {\"jsonContent\": %s}\n\nevent: %s\n\n", string(msgJson), su.Type)
	if _, err := c.Response().Write([]byte(message)); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}
