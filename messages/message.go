//here is the all message event structs with templates

package messages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEprotocol interface {
	Emit(msgChannel chan interface{}) error
}

// Message Items [2]
type ViewMessage struct {
	Type string `json:"msgType"`
	Data string `json:"update_value"`
}

type UpdateMessage struct { // Implementations - [ SSEprotocol ]
	Msg        string      `json:"msg"`
	Type       string      `json:"type"`
	Data       interface{} `json:"data"`
	Event      string      `json:"event"`
	Action     string      `json:"action"`
	Dependents string      `json:"dependents"`
	Component  string      `json:"component"`
}

type EventChanger struct {
	tmp interface{}
}

func (EV *EventChanger) ChangeEvent(newEvent string) {
	EV.tmp = newEvent
}

type TempMessage struct {
	Type string `json:"msgType"`
	Msg  string `json:"msg"`
}

func (tm *TempMessage) EmitTmp(c echo.Context) error {
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
	fmt.Print("Message sent")
	return nil
}

// fixed
func (su *UpdateMessage) Emit(msgChannel chan interface{}) error {

	msgJson, err := json.Marshal(su)
	if err != nil {
		fmt.Println("Error in marshalling json")
	}
	message := fmt.Sprintf("data: {\"jsonContent\": %s}\n\nevent: %s\n\n", string(msgJson), su.Type)

	msgChannel <- message

	return nil
}
