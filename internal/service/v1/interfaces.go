package v1service

type UserService interface {
	GetAllUsers(search string, page int, limit int)
	GetUserByUUID(uuid string)
	CreateUser()
	UpdateUser(uuid string)
	DeleteUser(uuid string)
}
