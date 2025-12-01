package service

type UserService interface {
	GetAllUsers()
	GetUserByUUID()
	CreateUser()
	UpdateUser()
	DeleteUser()
}
