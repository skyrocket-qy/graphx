package utils

import (
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	sqldomain "github.com/skyrocketOoO/zanazibar-dag/domain/infra/sql"
)

func RelationToString(tuple domain.Relation) string {
	res := tuple.ObjectNamespace + ":" + tuple.ObjectName + "#" + tuple.Relation
	res += "@" + tuple.SubjectNamespace + ":" + tuple.SubjectName
	if tuple.SubjectRelation != "" {
		res += "#" + tuple.SubjectRelation
	}

	return res
}

func ConvertRelation(in sqldomain.Relation) domain.Relation {
	return domain.Relation{
		ObjectNamespace:  in.ObjectNamespace,
		ObjectName:       in.ObjectName,
		Relation:         in.Relation,
		SubjectNamespace: in.SubjectNamespace,
		SubjectName:      in.SubjectName,
		SubjectRelation:  in.SubjectRelation,
	}
}
