package usecase

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type Usecase struct {
	repo domain.DatabaseRepository
}

func NewUsecase(ormRepo domain.DatabaseRepository) *Usecase {
	return &Usecase{
		repo: ormRepo,
	}
}

func (u *Usecase) Healthy(c context.Context) error {
	// do something check like db connection is established
	if err := u.repo.Ping(c); err != nil {
		return err
	}

	return nil
}
