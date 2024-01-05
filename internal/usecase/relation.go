package usecase

import (
	"strconv"
	"zanzibar-dag/domain"
	sqldomain "zanzibar-dag/domain/infra/sql"
	"zanzibar-dag/utils"
	"zanzibar-dag/utils/queue"

	"github.com/spf13/viper"
)

type RelationUsecase struct {
	RelationRepo sqldomain.RelationRepository
}

func NewRelationUsecase(relationRepo sqldomain.RelationRepository) *RelationUsecase {
	return &RelationUsecase{
		RelationRepo: relationRepo,
	}
}

func (u *RelationUsecase) GetAll() ([]domain.Relation, error) {
	tuples, err := u.RelationRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return tuples, nil
}

func (u *RelationUsecase) Query(query domain.Relation) ([]domain.Relation, error) {
	return u.RelationRepo.Query(query)
}

func (u *RelationUsecase) Create(relation domain.Relation) error {
	ok, err := u.Check(
		domain.Node{
			Namespace: relation.ObjectNamespace,
			Name:      relation.ObjectName,
			Relation:  relation.Relation,
		},
		domain.Node{
			Namespace: relation.SubjectNamespace,
			Name:      relation.SubjectName,
			Relation:  relation.SubjectRelation,
		},
	)
	if err != nil {
		return err
	}
	if ok {
		return domain.CauseCycleError{}
	}
	return u.RelationRepo.Create(relation)
}

func (u *RelationUsecase) Delete(relation domain.Relation) error {
	return u.RelationRepo.Delete(relation)
}

func (u *RelationUsecase) GetAllNamespaces() ([]string, error) {
	return u.RelationRepo.GetAllNamespaces()
}

func (u *RelationUsecase) Check(from domain.Node, to domain.Node) (bool, error) {
	// queryTimes := 0
	// defer func() {
	// 	fmt.Println(queryTimes)
	// }()
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return false, err
	}

	visited := utils.NewSet[domain.Node]()
	q := queue.NewQueue[domain.Node]()
	visited.Add(from)
	q.Push(from)

	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			node, _ := q.Pop()
			query := domain.Relation{
				SubjectNamespace: node.Namespace,
				SubjectName:      node.Name,
				SubjectRelation:  node.Relation,
			}
			tuples, err := u.RelationRepo.Query(query)
			// queryTimes += 1
			if err != nil {
				return false, err
			}

			for _, tuple := range tuples {
				if tuple.ObjectNamespace == to.Namespace && tuple.ObjectName == to.Name && tuple.Relation == to.Relation {
					return true, nil
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if !visited.Exist(child) {
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

	return false, nil
}

func (u *RelationUsecase) GetShortestPath(from domain.Node, to domain.Node) ([]domain.Relation, error) {
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return nil, err
	}
	visited := utils.NewSet[domain.Node]()
	type NodeItem struct {
		Cur  domain.Node
		Path []domain.Relation
	}
	firstNode := NodeItem{
		Cur:  from,
		Path: []domain.Relation{},
	}
	q := queue.NewQueue[NodeItem]()
	visited.Add(from)
	q.Push(firstNode)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			node, _ := q.Pop()
			query := domain.Relation{
				SubjectNamespace: node.Cur.Namespace,
				SubjectName:      node.Cur.Name,
				SubjectRelation:  node.Cur.Relation,
			}
			tuples, err := u.RelationRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, tuple := range tuples {
				if tuple.ObjectNamespace == to.Namespace && tuple.ObjectName == to.Name && tuple.Relation == to.Relation {
					return append(node.Path, tuple), nil
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if !visited.Exist(child) {
					visited.Add(child)
					copyPath := append(node.Path, tuple)
					q.Push(NodeItem{
						Cur:  child,
						Path: copyPath,
					})
				}
			}
		}
		depth++
		if depth >= maxDepth {
			break
		}
	}

	return nil, nil

	// maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	// if err != nil {
	// 	return nil, err
	// }
	// visited := utils.NewSet[domain.Node]()
	// type NodeItem struct {
	// 	Cur  domain.Node
	// 	Path []domain.Relation
	// }
	// finalPath := []domain.Relation{}
	// var dfsErr error
	// var dfs func(from, to domain.Node, path *[]domain.Relation, cur_depth int)
	// dfs = func(from, to domain.Node, path *[]domain.Relation, cur_depth int) {
	// 	if visited.Exist(from) || cur_depth >= maxDepth || len(finalPath) != 0 || dfsErr != nil {
	// 		return
	// 	}
	// 	visited.Add(from)
	// 	if from == to {
	// 		finalPath = *path
	// 		return
	// 	}
	// 	relations, err := u.RelationRepo.Query(domain.Relation{
	// 		SubjectNamespace: from.Namespace,
	// 		SubjectName:      from.Name,
	// 		SubjectRelation:  from.Relation,
	// 	})
	// 	if err != nil {
	// 		dfsErr = err
	// 		return
	// 	}
	// 	for _, relation := range relations {
	// 		*path = append(*path, relation)
	// 		child := domain.Node{
	// 			Namespace: relation.ObjectNamespace,
	// 			Name:      relation.ObjectName,
	// 			Relation:  relation.Relation,
	// 		}
	// 		dfs(child, to, path, cur_depth+1)
	// 		*path = (*path)[:len(*path)-1]
	// 	}
	// }
	// initPath := []domain.Relation{}
	// dfs(from, to, &initPath, 1)
	// return finalPath, nil
}

func (u *RelationUsecase) GetAllPaths(from domain.Node, to domain.Node) ([][]domain.Relation, error) {
	paths := [][]domain.Relation{}
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return nil, err
	}
	type NodeItem struct {
		Cur  domain.Node
		Path []domain.Relation
	}
	firstNode := NodeItem{
		Cur:  from,
		Path: []domain.Relation{},
	}
	q := queue.NewQueue[NodeItem]()
	q.Push(firstNode)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			node, _ := q.Pop()
			query := domain.Relation{
				SubjectNamespace: node.Cur.Namespace,
				SubjectName:      node.Cur.Name,
				SubjectRelation:  node.Cur.Relation,
			}
			tuples, err := u.RelationRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, tuple := range tuples {
				if tuple.ObjectNamespace == to.Namespace && tuple.ObjectName == to.Name && tuple.Relation == to.Relation {
					paths = append(paths, append(node.Path, tuple))
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				copyPath := append(node.Path, tuple)
				q.Push(NodeItem{
					Cur:  child,
					Path: copyPath,
				})

			}
		}
		depth++
		if depth >= maxDepth {
			break
		}
	}

	return paths, nil
}

func (u *RelationUsecase) GetAllObjectRelations(from domain.Node) ([]domain.Relation, error) {
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return nil, err
	}
	relations := utils.NewSet[domain.Relation]()
	visited := utils.NewSet[domain.Node]()
	q := queue.NewQueue[domain.Node]()
	visited.Add(from)
	q.Push(from)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			node, _ := q.Pop()
			query := domain.Relation{
				SubjectNamespace: node.Namespace,
				SubjectName:      node.Name,
				SubjectRelation:  node.Relation,
			}
			tuples, err := u.RelationRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, tuple := range tuples {
				relations.Add(tuple)
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if !visited.Exist(child) {
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

	return relations.ToSlice(), nil
}

func (u *RelationUsecase) GetAllSubjectRelations(object domain.Node) ([]domain.Relation, error) {
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return nil, err
	}
	relations := utils.NewSet[domain.Relation]()
	visited := utils.NewSet[domain.Node]()
	q := queue.NewQueue[domain.Node]()
	visited.Add(object)
	q.Push(object)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			node, _ := q.Pop()
			query := domain.Relation{
				ObjectNamespace: node.Namespace,
				ObjectName:      node.Name,
				Relation:        node.Relation,
			}
			tuples, err := u.RelationRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, tuple := range tuples {
				relations.Add(tuple)
				parent := domain.Node{
					Namespace: tuple.SubjectNamespace,
					Name:      tuple.SubjectName,
					Relation:  tuple.SubjectRelation,
				}
				if !visited.Exist(parent) {
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

	return relations.ToSlice(), nil
}

func (u *RelationUsecase) ClearAllRelations() error {
	return u.RelationRepo.DeleteAll()
}
