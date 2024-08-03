package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// sample struct > this will be stored in redis..
type UserWS struct {
	uacc uint16
	msg  chan []byte
	uid  uint32
}

type UserRequestWS struct {
	UID  uint32 `json:"uid"`
	UGID uint16 `json:"ugid"`
}

type WSockService struct {
	Multiplexer  http.ServeMux
	usersMu      sync.Mutex
	dataBuffer   uint64
	UserRegistry map[uint32]*UserWS
}

func newWSService() *WSockService {
	ws := &WSockService{
		dataBuffer:   30,
		UserRegistry: make(map[uint32]*UserWS),
	}

	ws.Multiplexer.Handler("/wsconn", ws.InitiateService)
}

// MESSAGE RETRIEVEL PROCESSORS----
func captureMessage(ws *websocket.Conn, ctx context.Context, msgchannel chan<- []byte) error {
	// Create a buffer to hold chunks of the incoming message
	_, reader, err := ws.Reader(ctx)
	if err != nil {
		fmt.Printf("error in Reader..")
	}
	buffer := make([]byte, 1024) // given - 1MB

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		buffer = make([]byte, 1024) //flushing the buffer

		msgchannel <- data
	}
	return nil
}

// got messaage processing
func processMessage(data []byte) error {
	// Example: just print the data
	fmt.Printf("Received: %s", data)
	return nil
}

//Message Senders-----

func (ws *WSockService) ServeHTTP(r *http.Request, w http.ResponseWriter) {
	ws.Multiplexer.ServeHTTP(w, r)
}

func (ws *WSockService) InitiateService(r *http.Request, w http.ResponseWriter) {
	defer r.Body.Close()
	fmt.Println("WSService starrted, loading user to the memory...")
	ws.startWSConnection(r.Context(), r, w)
}

func (ws *WSockService) startWSConnection(ctx context.Context, r *http.Request, w http.ResponseWriter) {

	var mu sync.Mutex
	var wsocket *websocket.Conn //tHis is the WS for the user..
	var closed bool

	//----extra\ting user details..
	var uData UserRequestWS
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&uData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//---- extraction done ---

	usr := &UserWS{
		uacc: uData.UGID,
		msg:  make(chan []byte, ws.dataBuffer),
		uid:  uData.UID,
	}

	ws.AddUser(usr) // here user foes to user registry in ws service
	defer ws.flushUser(uData.UID)

	wsocket, err = websocket.Accept(w, r, nil)
	defer wsocket.CloseNow()
	//ctx = wsocket.CloseRead(ctx) //TODO: check what happens here

	//CAPTURE - processing retrieval
	getMessage := make(chan []byte, 1024)
	go captureMessage(wsocket, ctx, getMessage)
	go func() {
		for msg := range getMessage {
			if err := processMessage(msg); err != nil {
				log.Printf("error processing message: %v", err)
				break
			}
		}
	}()

	//EMITTER - sending processor
	for {
		select {
		case MSG := <-usr.msg:
			err := sendDataThroughWS(ctx, time.Second*5, wsocket, MSG)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}

	}

}
func (ws *WSockService) AddUser(user *UserWS) {
	ws.usersMu.Lock()
	ws.UserRegistry[user.uid] = user
	ws.usersMu.Unlock()
}

func (ws *WSockService) flushUser(uid uint32) {
	ws.usersMu.Lock()
	delete(ws.UserRegistry, uid)
	ws.usersMu.Unlock()
}

func sendDataThroughWS(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}
