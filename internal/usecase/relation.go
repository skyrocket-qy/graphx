package usecase

import (
	"errors"
	"fmt"
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

func (u *RelationUsecase) GetAll() ([]string, error) {
	tuples, err := u.RelationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	relations := []string{}
	for _, tuple := range tuples {
		relations = append(relations, utils.RelationToString(tuple))
	}

	return relations, nil
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
		return errors.New("create cycle detected")
	}
	return u.RelationRepo.Create(relation)
}

func (u *RelationUsecase) Delete(relation domain.Relation) error {
	return u.RelationRepo.Delete(relation)
}

func (u *RelationUsecase) Check(from domain.Node, to domain.Node) (bool, error) {
	queryTimes := 0
	defer func() {
		fmt.Println(queryTimes)
	}()
	depth := 0
	maxDepth, err := strconv.Atoi(viper.GetString("main.max-search-depth"))
	if err != nil {
		return false, err
	}

	visited := utils.NewSet[domain.Relation]()

	firstQuery := domain.Relation{
		SubjectNamespace: from.Namespace,
		SubjectName:      from.Name,
		SubjectRelation:  from.Relation,
	}
	q := queue.NewQueue[domain.Relation]()
	visited.Add(firstQuery)
	q.Push(firstQuery)

	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			query, _ := q.Pop()
			tuples, err := u.RelationRepo.Query(query)
			queryTimes += 1
			if err != nil {
				return false, err
			}

			for _, tuple := range tuples {
				if tuple.ObjectNamespace == to.Namespace && tuple.ObjectName == to.Name && tuple.Relation == to.Relation {
					return true, nil
				}
				nextQuery := domain.Relation{
					SubjectNamespace: tuple.ObjectNamespace,
					SubjectName:      tuple.ObjectName,
					SubjectRelation:  tuple.Relation,
				}
				if !visited.Exist(nextQuery) {
					visited.Add(nextQuery)
					q.Push(nextQuery)
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

func (u *RelationUsecase) GetShortestPath(from domain.Node, to domain.Node) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (u *RelationUsecase) GetAllPaths(from domain.Node, to domain.Node) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (u *RelationUsecase) GetObjectRelations(from domain.Node) ([]string, error) {
	objectRelations := utils.NewSet[string]()
	visited := utils.NewSet[domain.Relation]()

	firstQuery := domain.Relation{
		SubjectNamespace: from.Namespace,
		SubjectName:      from.Name,
		SubjectRelation:  from.Relation,
	}
	q := queue.NewQueue[domain.Relation]()
	visited.Add(firstQuery)
	q.Push(firstQuery)

	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			query, _ := q.Pop()
			tuples, err := u.RelationRepo.Query(query)
			if err != nil {
				return nil, err
			}

			for _, tuple := range tuples {
				objectRelations.Add(utils.RelationToString(tuple))
				nextQuery := domain.Relation{
					SubjectNamespace: tuple.ObjectNamespace,
					SubjectName:      tuple.ObjectName,
					SubjectRelation:  tuple.Relation,
				}
				if !visited.Exist(nextQuery) {
					visited.Add(nextQuery)
					q.Push(nextQuery)
				}
			}
		}
	}

	return objectRelations.ToSlice(), nil
}

func (u *RelationUsecase) ClearAllRelations() error {
	return u.RelationRepo.DeleteAll()
}
