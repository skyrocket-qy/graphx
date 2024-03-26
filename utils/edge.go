package utils

import (
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"github.com/skyrocketOoO/zanazibar-dag/internal/infra/repository/sql"
)

func EdgeToString(edge domain.Edge) string {
	res := edge.ObjNs + ":" + edge.ObjName + "#" + edge.ObjRel
	res += "@" + edge.SbjNs + ":" + edge.SbjName
	if edge.SbjRel != "" {
		res += "#" + edge.SbjRel
	}

	return res
}

func ConvertRelation(in sql.Edge) domain.Edge {
	return domain.Edge{
		ObjNs:   in.ObjNs,
		ObjName: in.ObjName,
		ObjRel:  in.ObjRel,
		SbjNs:   in.SbjNs,
		SbjName: in.SbjName,
		SbjRel:  in.SbjRel,
	}
}
