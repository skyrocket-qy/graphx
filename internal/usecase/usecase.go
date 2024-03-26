package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/skyrocketOoO/go-utility/queue"
	"github.com/skyrocketOoO/go-utility/set"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"github.com/skyrocketOoO/zanazibar-dag/utils"
	"gorm.io/gorm"
)

type Usecase struct {
	sqlRepo        domain.SqlRepository
	PageStates     map[string]*PageState
	PageStatesLock sync.RWMutex
}

func NewUsecase(sqlRepo domain.SqlRepository) *Usecase {
	usecase := &Usecase{
		sqlRepo:        sqlRepo,
		PageStates:     map[string]*PageState{},
		PageStatesLock: sync.RWMutex{},
	}

	// go func(u *Usecase) {
	// 	for {
	// 		time.Sleep(time.Second * 10)

	// 		shouldDelete := []string{}
	// 		u.PageStatesLock.RLock()
	// 		for key, pageState := range u.PageStates {
	// 			if pageState.ExpiredTime.Before(time.Now()) {
	// 				shouldDelete = append(shouldDelete, key)
	// 			}
	// 		}
	// 		u.PageStatesLock.RUnlock()
	// 		u.PageStatesLock.Lock()
	// 		for _, key := range shouldDelete {
	// 			delete(u.PageStates, key)
	// 		}
	// 		u.PageStatesLock.Unlock()
	// 	}
	// }(usecase)

	return usecase
}

type PageState struct {
	LastRelID uint
	// if query after expired time without remove, keep query
	ExpiredTime time.Time
}

func (u *Usecase) Healthy(c context.Context) error {
	// do something check like db connection is established
	if err := u.sqlRepo.Ping(c); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Get(edge domain.Edge, options ...domain.PageOptions) (
	[]domain.Edge, string, error) {
	// if edge.ObjNs == "" && edge.ObjName == "" && edge.ObjRel == "" &&
	// 	edge.SbjNs == "" && edge.SbjName == "" && edge.SbjRel == "" {
	// var lastID uint
	// var pageSize int
	// if len(options) > 0 {
	// 	options := options[0]
	// 	pageSize = options.PageSize
	// 	u.PageStatesLock.RLock()
	// 	pageState, ok := u.PageStates[options.PageToken]
	// 	u.PageStatesLock.RUnlock()
	// 	if !ok {
	// 		return nil, "", fmt.Errorf("previouse page state not found")
	// 	}
	// 	lastID = pageState.LastRelID

	// 	u.PageStatesLock.Lock()
	// 	delete(u.PageStates, options.PageToken)
	// 	u.PageStatesLock.Unlock()
	// }
	// edges, lastID, err := u.sqlRepo.GetAll(domain.PageOptions{
	// 	LastID:   lastID,
	// 	PageSize: pageSize,
	// })
	// if err != nil {
	// 	return nil, "", err
	// }
	// pageState := PageState{
	// 	LastRelID:   lastID,
	// 	ExpiredTime: time.Now().Add(time.Minute * 5),
	// }
	// token, err := utils.GenerateRandomToken()
	// if err != nil {
	// 	return nil, "", err
	// }
	// u.PageStatesLock.Lock()
	// for _, ok := u.PageStates[token]; ok; {
	// 	token, err = utils.GenerateRandomToken()
	// 	if err != nil {
	// 		return nil, "", err
	// 	}
	// }
	// u.PageStates[token] = &pageState
	// u.PageStatesLock.Unlock()

	// return edges, token, nil
	// } else {
	edges, err := u.sqlRepo.Query(edge)
	return edges, "", err
	// }
}

func (u *Usecase) Create(edge domain.Edge, existOk bool) error {
	if err := utils.ValidateRel(edge); err != nil {
		return err
	}
	ok, err := u.Check(
		domain.Vertex{
			Ns:   edge.ObjNs,
			Name: edge.ObjName,
			Rel:  edge.ObjRel,
		},
		domain.Vertex{
			Ns:   edge.SbjNs,
			Name: edge.SbjName,
			Rel:  edge.SbjRel,
		},
		domain.SearchCond{},
	)
	if err != nil {
		return err
	}
	if ok {
		return domain.CauseCycleError{}
	}

	err = u.sqlRepo.Create(edge)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			if existOk {
				return nil
			} else {
				return err
			}
		}
		return err
	}
	return nil
}

func (u *Usecase) Delete(edge domain.Edge) error {
	if err := utils.ValidateRel(edge); err != nil {
		return err
	}
	return u.sqlRepo.Delete(edge)
}

func (u *Usecase) DeleteByQueries(queries []domain.Edge) error {
	return u.sqlRepo.DeleteByQueries(queries)
}

func (u *Usecase) BatchOperation(operations []domain.Operation) error {
	for _, operation := range operations {
		if err := utils.ValidateRel(operation.Edge); err != nil {
			return err
		}
	}
	return u.sqlRepo.BatchOperation(operations)
}

func (u *Usecase) GetAllNs() ([]string, error) {
	return u.sqlRepo.GetAllNs()
}

func (u *Usecase) Check(sbj domain.Vertex, obj domain.Vertex,
	searchCond domain.SearchCond) (bool, error) {
	if err := utils.ValidateVertex(obj, false); err != nil {
		return false, err
	}
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return false, err
	}
	visited := set.NewSet[domain.Vertex]()
	q := queue.NewQueue[domain.Vertex]()
	visited.Add(sbj)
	q.Push(sbj)

	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			query := domain.Edge{
				SbjNs:   vertex.Ns,
				SbjName: vertex.Name,
				SbjRel:  vertex.Rel,
			}
			edges, err := u.sqlRepo.Query(query)
			if err != nil {
				return false, err
			}

			for _, edge := range edges {
				if edge.ObjNs == obj.Ns && edge.ObjName == obj.Name &&
					edge.ObjRel == obj.Rel {
					return true, nil
				}
				child := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if !searchCond.ShouldStop(child) &&
					!visited.Exist(child) {
					visited.Add(child)
					q.Push(child)
				}
			}
		}
	}

	return false, nil
}

func (u *Usecase) GetShortestPath(sbj domain.Vertex, obj domain.Vertex,
	searchCond domain.SearchCond) ([]domain.Edge, error) {
	if err := utils.ValidateVertex(obj, false); err != nil {
		return nil, err
	}
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return nil, err
	}
	visited := set.NewSet[domain.Vertex]()
	type NodeItem struct {
		Cur  domain.Vertex
		Path []domain.Edge
	}
	firstNode := NodeItem{
		Cur:  sbj,
		Path: []domain.Edge{},
	}
	q := queue.NewQueue[NodeItem]()
	visited.Add(sbj)
	q.Push(firstNode)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			query := domain.Edge{
				SbjNs:   vertex.Cur.Ns,
				SbjName: vertex.Cur.Name,
				SbjRel:  vertex.Cur.Rel,
			}
			edges, err := u.sqlRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, edge := range edges {
				if edge.ObjNs == obj.Ns && edge.ObjName == obj.Name &&
					edge.ObjRel == obj.Rel {
					return append(vertex.Path, edge), nil
				}
				child := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if !searchCond.ShouldStop(child) &&
					!visited.Exist(child) {
					visited.Add(child)
					copyPath := append(vertex.Path, edge)
					q.Push(NodeItem{
						Cur:  child,
						Path: copyPath,
					})
				}
			}
		}
	}

	return nil, nil
}

func (u *Usecase) GetAllPaths(sbj domain.Vertex, obj domain.Vertex,
	searchCond domain.SearchCond) ([][]domain.Edge, error) {
	if err := utils.ValidateVertex(obj, false); err != nil {
		return nil, err
	}
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return nil, err
	}
	paths := [][]domain.Edge{}
	type NodeItem struct {
		Cur  domain.Vertex
		Path []domain.Edge
	}
	firstNode := NodeItem{
		Cur:  sbj,
		Path: []domain.Edge{},
	}
	q := queue.NewQueue[NodeItem]()
	q.Push(firstNode)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			query := domain.Edge{
				SbjNs:   vertex.Cur.Ns,
				SbjName: vertex.Cur.Name,
				SbjRel:  vertex.Cur.Rel,
			}
			edges, err := u.sqlRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, edge := range edges {
				if edge.ObjNs == obj.Ns && edge.ObjName == obj.Name &&
					edge.ObjRel == obj.Rel {
					paths = append(paths, append(vertex.Path, edge))
				}
				child := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if searchCond.ShouldStop(child) {
					continue
				}
				copyPath := append(vertex.Path, edge)
				q.Push(NodeItem{
					Cur:  child,
					Path: copyPath,
				})

			}
		}
	}

	return paths, nil
}

func (u *Usecase) GetAllObjRels(sbj domain.Vertex,
	searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) ([]domain.Edge, error) {
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return nil, err
	}
	depth := 0
	edges := set.NewSet[domain.Edge]()
	visited := set.NewSet[domain.Vertex]()
	q := queue.NewQueue[domain.Vertex]()
	visited.Add(sbj)
	q.Push(sbj)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			query := domain.Edge{
				SbjNs:   vertex.Ns,
				SbjName: vertex.Name,
				SbjRel:  vertex.Rel,
			}
			qEdges, err := u.sqlRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, edge := range qEdges {
				child := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if collectCond.ShouldCollect(child) {
					edges.Add(edge)
				}
				if !searchCond.ShouldStop(child) &&
					!visited.Exist(child) {
					visited.Add(child)
					q.Push(child)
				}
			}
		}
		depth++
		if depth >= maxDepth {
			break
		}
	}

	return edges.ToSlice(), nil
}

func (u *Usecase) GetAllSbjRels(obj domain.Vertex,
	searchCond domain.SearchCond,
	collectCond domain.CollectCond,
	maxDepth int) ([]domain.Edge, error) {
	if err := utils.ValidateVertex(obj, false); err != nil {
		return nil, err
	}
	depth := 0
	edges := set.NewSet[domain.Edge]()
	visited := set.NewSet[domain.Vertex]()
	q := queue.NewQueue[domain.Vertex]()
	visited.Add(obj)
	q.Push(obj)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			query := domain.Edge{
				ObjNs:   vertex.Ns,
				ObjName: vertex.Name,
				ObjRel:  vertex.Rel,
			}
			qEdges, err := u.sqlRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, edge := range qEdges {
				parent := domain.Vertex{
					Ns:   edge.SbjNs,
					Name: edge.SbjName,
					Rel:  edge.SbjRel,
				}
				if collectCond.ShouldCollect(parent) {
					edges.Add(edge)
				}
				if !searchCond.ShouldStop(parent) &&
					!visited.Exist(parent) {
					visited.Add(parent)
					q.Push(parent)
				}
			}
		}
		depth++
		if depth >= maxDepth {
			break
		}
	}

	return edges.ToSlice(), nil
}

func (u *Usecase) GetTree(sbj domain.Vertex, maxDepth int) (
	*domain.TreeNode, error) {
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return &domain.TreeNode{}, err
	}
	depth := 0
	head := &domain.TreeNode{
		Ns:       sbj.Ns,
		Name:     sbj.Name,
		Rel:      sbj.Rel,
		Children: []domain.TreeNode{},
	}
	q := queue.NewQueue[*domain.TreeNode]()
	visited := set.NewSet[**domain.TreeNode]()
	q.Push(head)
	for !q.IsEmpty() {
		depth++
		if depth > maxDepth {
			break
		}
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			parent, err := q.Pop()
			if err != nil {
				return head, err
			}
			if visited.Exist(&parent) {
				continue
			}
			visited.Add(&parent)
			query := domain.Edge{
				SbjNs:   parent.Ns,
				SbjName: parent.Name,
				SbjRel:  parent.Rel,
			}
			edges, err := u.sqlRepo.Query(query)
			if err != nil {
				return head, err
			}
			for _, r := range edges {
				treeNode := domain.TreeNode{
					Ns:       r.ObjNs,
					Name:     r.ObjName,
					Rel:      r.ObjRel,
					Children: []domain.TreeNode{},
				}
				parent.Children = append(parent.Children, treeNode)
				q.Push(&treeNode)
			}
		}
	}

	return head, nil
}

func (u *Usecase) ClearAllEdges() error {
	return u.sqlRepo.DeleteAll()
}
