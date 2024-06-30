package registries

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

func (ureg *UserRegistry) Release() {
	ureg.muLock.Lock()
	defer ureg.muLock.Unlock()
	ureg.sse_map = make(map[string]echo.Context)
}

// TODO: implment more structures usr type without user id like here | 2nd phase
// mutext has applied in case haver to use go routines in future over uRegistry..
func (ureg *UserRegistry) AddUser(user string, sseContext echo.Context) {
	ureg.muLock.Lock()
	defer ureg.muLock.Unlock()
	ureg.sse_map[user] = sseContext
}

func (ureg *UserRegistry) Notify(user string, msg messages.SSEprotocol) error { // msg - msg structs in messages that impl -- SSEprotocol must be given

	if err := msg.Emit(ureg.sse_map[user]); err != nil {
		return err
	}
	return nil
}

func (ureg *UserRegistry) RecordSession(user string) {
	ureg.muLock.Lock()
	defer ureg.muLock.Unlock()
	delete(ureg.sse_map, user)
}
