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
		domain.SearchCondition{},
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

func (u *RelationUsecase) Check(subject domain.Node, object domain.Node, searchCondition domain.SearchCondition) (bool, error) {
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return false, err
	}

	visited := utils.NewSet[domain.Node]()
	q := queue.NewQueue[domain.Node]()
	visited.Add(subject)
	q.Push(subject)

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
				return false, err
			}

			for _, tuple := range tuples {
				if tuple.ObjectNamespace == object.Namespace && tuple.ObjectName == object.Name && tuple.Relation == object.Relation {
					return true, nil
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if !searchCondition.ShouldStop(child) && !visited.Exist(child) {
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

func (u *RelationUsecase) GetShortestPath(subject domain.Node, object domain.Node, searchCondition domain.SearchCondition) ([]domain.Relation, error) {
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
		Cur:  subject,
		Path: []domain.Relation{},
	}
	q := queue.NewQueue[NodeItem]()
	visited.Add(subject)
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
				if tuple.ObjectNamespace == object.Namespace && tuple.ObjectName == object.Name && tuple.Relation == object.Relation {
					return append(node.Path, tuple), nil
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if !searchCondition.ShouldStop(child) && !visited.Exist(child) {
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
	// var dfs func(subject, object domain.Node, path *[]domain.Relation, cur_depth int)
	// dfs = func(subject, object domain.Node, path *[]domain.Relation, cur_depth int) {
	// 	if visited.Exist(subject) || cur_depth >= maxDepth || len(finalPath) != 0 || dfsErr != nil {
	// 		return
	// 	}
	// 	visited.Add(subject)
	// 	if subject == object {
	// 		finalPath = *path
	// 		return
	// 	}
	// 	relations, err := u.RelationRepo.Query(domain.Relation{
	// 		SubjectNamespace: subject.Namespace,
	// 		SubjectName:      subject.Name,
	// 		SubjectRelation:  subject.Relation,
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
	// 		dfs(child, object, path, cur_depth+1)
	// 		*path = (*path)[:len(*path)-1]
	// 	}
	// }
	// initPath := []domain.Relation{}
	// dfs(subject, object, &initPath, 1)
	// return finalPath, nil
}

func (u *RelationUsecase) GetAllPaths(subject domain.Node, object domain.Node, searchCondition domain.SearchCondition) ([][]domain.Relation, error) {
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
		Cur:  subject,
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
				if tuple.ObjectNamespace == object.Namespace && tuple.ObjectName == object.Name && tuple.Relation == object.Relation {
					paths = append(paths, append(node.Path, tuple))
				}
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if searchCondition.ShouldStop(child) {
					continue
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

func (u *RelationUsecase) GetAllObjectRelations(subject domain.Node, searchCondition domain.SearchCondition, collectCondition domain.CollectCondition) ([]domain.Relation, error) {
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return nil, err
	}
	relations := utils.NewSet[domain.Relation]()
	visited := utils.NewSet[domain.Node]()
	q := queue.NewQueue[domain.Node]()
	visited.Add(subject)
	q.Push(subject)
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
				child := domain.Node{
					Namespace: tuple.ObjectNamespace,
					Name:      tuple.ObjectName,
					Relation:  tuple.Relation,
				}
				if collectCondition.ShouldCollect(child) {
					relations.Add(tuple)
				}
				if !searchCondition.ShouldStop(child) && !visited.Exist(child) {
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

func (u *RelationUsecase) GetAllSubjectRelations(object domain.Node, searchCondition domain.SearchCondition, collectCondition domain.CollectCondition) ([]domain.Relation, error) {
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
				parent := domain.Node{
					Namespace: tuple.SubjectNamespace,
					Name:      tuple.SubjectName,
					Relation:  tuple.SubjectRelation,
				}
				if collectCondition.ShouldCollect(parent) {
					relations.Add(tuple)
				}
				if !searchCondition.ShouldStop(parent) && !visited.Exist(parent) {
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
