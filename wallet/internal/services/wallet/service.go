package wallet

type Repository interface {
	GetBalance(walletID string) (int64, error)
	GetCreditLimit(walletID string) (int64, error)
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetPaymentCapacity Returns sum of wallet balance and credit allowance
func (s *Service) GetPaymentCapacity(walletID string) (int64, error) {
	balance, err := s.repo.GetBalance(walletID)
	if err != nil {
		return 0, err
	}

	credit, err := s.repo.GetCreditLimit(walletID)
	if err != nil {
		return 0, err
	}

	return balance + credit, nil
}
