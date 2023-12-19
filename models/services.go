package models

// Services contains the services which need db connection
type Services struct {
	User    UserService
	Content ContentService
}

// NewServices initialises all services with a single db connection
func NewServices(cfgs ...ServicesConfig) (*Services, error) {

	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Close closes the database connections.
func (s *Services) Close() error {
	err := s.User.CloseDB()
	if err != nil {
		return err
	}
	err = s.Content.CloseDB()
	return err
}
