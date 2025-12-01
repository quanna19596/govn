package repository

type UserRepository interface {
	FindAll()
	FindByUUID()
	Create()
	Update()
	Delete()
}
