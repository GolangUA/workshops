package calendar

type Repository interface {
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
