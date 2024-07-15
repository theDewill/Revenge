package registries

import (
	"fmt"
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

type User struct {
	User_id     int
	Sse_channel chan interface{}
}

// 2KB
type Usergroup struct {
	User_grp_identity int
	User_list         []User
}

// 2 x 10^6 KB = 2GB (1 million user)
type UserRegistry struct { // Imp - [ RegistryModifications, RegistryCommunications ]
	MuLock     sync.Mutex
	Validity   string
	usr_groups []Usergroup // the mapping for each user with respective sse pipeline
}

func NewUserRegitry() *UserRegistry {
	return &UserRegistry{
		Validity:   "valid",
		usr_groups: make([]Usergroup, 0, 100), //TODO: OPTIMIZE point -
	}
}

func (ureg *UserRegistry) Release() {
	ureg.MuLock.Lock()
	defer ureg.MuLock.Unlock()
	ureg.usr_groups = make([]Usergroup, 0, 100) // flushing all
}

// Expect that in the inital login group details and user details are retrived
func (ureg *UserRegistry) LoadUser(user_id int, user_grp_id int) *User {
	// ureg.MuLock.Lock()
	// defer ureg.MuLock.Unlock()

	//checking wether hes online :TODO: verfy

	if len(ureg.usr_groups) != 0 {

		if len(ureg.usr_groups[user_grp_id-1].User_list) != 0 {
			for _, user := range ureg.usr_groups[user_grp_id-1].User_list {
				if user.User_id == user_id {
					return &user
				} else {
					return &User{
						User_id:     user_id,
						Sse_channel: make(chan interface{}),
					}
				}
			}
		} else {
			return ureg.CreateUser(user_id, user_grp_id)
		}
	} else {

		ureg.AddUserGroup(Usergroup{
			User_grp_identity: 1, //TODO: db operation to capture users group
			User_list:         make([]User, 0, 10),
		})

		return ureg.CreateUser(1, 1)
	}
	return nil

}

func (ureg *UserRegistry) CreateUser(user_id int, user_grp_id int) *User {
	ureg.MuLock.Lock()
	defer ureg.MuLock.Unlock()
	usr := User{
		User_id:     user_id,
		Sse_channel: make(chan interface{}),
	}

	ureg.usr_groups[user_grp_id-1].User_list = append(ureg.usr_groups[user_grp_id-1].User_list, usr)
	fmt.Println("User created")
	return &usr
}

// TODO: implment more structures usr type without user id like here | 2nd phase
// mutext has applied in case haver to use go routines in future over uRegistry..
func (ureg *UserRegistry) AddUserGroup(user_grp Usergroup) {
	ureg.MuLock.Lock()
	defer ureg.MuLock.Unlock()

	ureg.usr_groups = append(ureg.usr_groups, user_grp)

}

func (ureg *UserRegistry) SendUpdates(senderID int, ugid int, msg messages.SSEprotocol) []User {
	ureg.MuLock.Lock()
	defer ureg.MuLock.Unlock()
	for _, usr := range ureg.usr_groups[ugid-1].User_list {
		if usr.User_id != senderID {
			msg.Emit(usr.Sse_channel)
		}
	}
	return nil
}
