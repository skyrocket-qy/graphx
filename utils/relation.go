package utils

import (
	"zanzibar-dag/domain"
	sqldomain "zanzibar-dag/domain/infra/sql"
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
