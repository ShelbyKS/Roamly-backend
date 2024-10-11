package user

type Service interface {
}

type serviceImpl struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &serviceImpl{
		storage: storage,
	}
}
