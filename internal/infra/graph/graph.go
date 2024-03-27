package graph

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type GraphInfra struct {
	sqlRepo domain.SqlRepository
}

func NewGraphInfra(sqlRepo domain.SqlRepository) *GraphInfra {
	return &GraphInfra{
		sqlRepo: sqlRepo,
	}
}

func (g *GraphInfra) Check(c context.Context, sbj domain.Vertex,
	obj domain.Vertex, searchCond domain.SearchCond) (
	found bool, err error) {
	return false, domain.ErrNotImplemented{}
}

func (g *GraphInfra) GetPassedVertices(c context.Context, start domain.Vertex,
	isSbj bool, searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) (vertices []domain.Vertex, err error) {
	return nil, domain.ErrNotImplemented{}
}

func (g *GraphInfra) GetTree(subject domain.Vertex, maxDepth int) (
	*domain.TreeNode, error) {
	return nil, domain.ErrNotImplemented{}
}
