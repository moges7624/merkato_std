package permission

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll() ([]Permission, error) {
	return s.store.getAll()
}

func (s *Service) GetAllForUser(userID int64) ([]Permission, error) {
	return s.store.getAllForUser(userID)
}

func (s *Service) AddForUser(req AddPermissionForUserRequest) error {
	return s.store.addForUser(req.UserID, req.Permissions)
}
