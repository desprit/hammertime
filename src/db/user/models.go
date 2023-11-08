package user_storage

import "desprit/hammertime/src/config"

const (
	UserD = "алешенька"
	UserM = "настенька"
)

type User struct {
	ID    int64
	Name  string
	Token string
}

var Users = []User{
	{
		ID:    0,
		Name:  UserD,
		Token: config.GetConfig().HAMMER_TOKEN_D,
	},
	{
		ID:    1,
		Name:  UserM,
		Token: config.GetConfig().HAMMER_TOKEN_M,
	},
}

var UserMapByName = map[string]User{
	UserD: Users[0],
	UserM: Users[1],
}

var UserMapByID = map[int64]User{
	Users[0].ID: Users[0],
	Users[1].ID: Users[1],
}
