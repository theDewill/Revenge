package Revenge

import (
	"ssego/registries"
)

//Assuming this wil be used in the handler to broadcast the dep update to other users in the group

func broadcast_users(ugrp *registries.Usergroup, msgData interface{}) {
	for _, user := range ugrp.User_list {
		user.Sse_channel <- msgData
	}
}
