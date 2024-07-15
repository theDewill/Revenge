package db

// user groups will be the partiotons in all the users - docs and assistants
type User struct {
	User_id     int
	User_grp_id int
}

type Usergroup struct {
	User_grp_id int
}
