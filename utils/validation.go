package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

func ValidateRel(rel domain.Edge) error {
	fmt.Printf("%+v\n", rel)
	if rel.ObjNs == "" || rel.ObjName == "" || rel.ObjRel == "" ||
		rel.SbjNs == "" || rel.SbjName == "" {
		return domain.ErrRequestBody{}
	}
	return ValidateReservedWord(rel)
}

func ValidateVertex(vertex domain.Vertex, isSubject bool) error {
	if vertex.Ns == "" || vertex.Name == "" {
		return domain.ErrRequestBody{}
	}
	if !isSubject && vertex.Rel == "" {
		return domain.ErrRequestBody{}
	}
	return ValidateReservedWord(vertex)
}

func ValidateReservedWord(st interface{}) error {
	value := reflect.ValueOf(st)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		str := field.Interface().(string)
		if strings.Contains(str, "%") {
			return domain.ErrRequestBody{}
		}

	}
	return nil
}
