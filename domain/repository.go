package domain

import "context"

type OrmRepository interface {
	Ping(c context.Context) error
}
