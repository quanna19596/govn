package repository

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	Create()
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}
