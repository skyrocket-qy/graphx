package usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	sqldom "github.com/skyrocketOoO/zanazibar-dag/domain/infra/sql"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"
)

func TestGetAllNamespaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRelationRepo := sqldom.NewMockRelationRepository(ctrl)
	mockRelationRepo.EXPECT().GetAllNamespaces().Return([]string{"foo", "bar"}, nil)

	usecaseRepo := usecase.NewRelationUsecase(mockRelationRepo)

	nss, err := usecaseRepo.GetAllNamespaces()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if nss[0] != "foo" && nss[1] != "bar" {
		t.Errorf("Unexpected error: %s", "nss should be 'foo' or 'bar")
	}
}
