package domain

import "context"

type Usecase interface {
	Healthy(ctx context.Context) error
}
