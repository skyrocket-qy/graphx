package domain

import "context"

type Usecase interface {
	Healthy(ctx context.Context) error

	Get(relation Edge, options ...PageOptions) (relations []Edge, token string, err error)
	Create(relation Edge, existOk bool) error
	Delete(relation Edge) error
	DeleteByQueries(queries []Edge) error
	BatchOperation(operations []Operation) error

	GetAllNs() ([]string, error)
	Check(subject Vertex, object Vertex, searchCond SearchCond) (bool, error)
	GetShortestPath(subject Vertex, object Vertex, searchCond SearchCond) ([]Edge, error)
	GetAllPaths(subject Vertex, object Vertex, searchCond SearchCond) ([][]Edge, error)
	GetAllObjRels(subject Vertex, searchCond SearchCond, collectCond CollectCond, maxDepth int) ([]Edge, error)
	GetAllSbjRels(object Vertex, searchCond SearchCond, collectCond CollectCond, maxDepth int) ([]Edge, error)
	GetTree(subject Vertex, maxDepth int) (*TreeNode, error)

	ClearAllEdges() error
}
