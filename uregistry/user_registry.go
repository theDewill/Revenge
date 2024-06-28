package uregistry

import (
	"ssego/messages"
	"sync"

	"github.com/labstack/echo/v4"
)

type RegistryModifications interface {
	AddUser(user string, sseContext *echo.Context)
	DisconnectUser(user string)
}

type RegistryCommunications interface {
	Notify(user string)
}

type UserRegistry struct { // Imp - [ RegistryModifications, RegistryCommunications ]
	muLock   sync.Mutex
	validity string
	sse_map  map[string]echo.Context // the mapping for each user with respective sse pipeline
}

func NewUserRegitry() *UserRegistry {
	ureg := UserRegistry{
		validity: "valid",
		sse_map:  make(map[string]echo.Context),
	}
	return &ureg
}

// TODO: implment more structures usr type without user id like here | 2nd phase
func (ureg *UserRegistry) AddUser(user string, sseContext echo.Context) {
	ureg.muLock.Lock()
	defer ureg.muLock.Unlock()
	ureg.sse_map[user] = sseContext
}

func (ureg *UserRegistry) Notify(user string, msg messages.UpdateMessage) error { // pass the message  object constructed from the message package struct

	if err := msg.Emit(ureg.sse_map[user]); err != nil {
		return err
	}
	return nil
}
