package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

func ValidateRelation(rel domain.Relation) error {
	fmt.Printf("%+v\n", rel)
	if rel.ObjectNamespace == "" || rel.ObjectName == "" || rel.Relation == "" ||
		rel.SubjectNamespace == "" || rel.SubjectName == "" {
		return domain.RequestBodyError{}
	}
	return ValidateReservedWord(rel)
}

func ValidateNode(node domain.Node, isSubject bool) error {
	if node.Namespace == "" || node.Name == "" {
		return domain.RequestBodyError{}
	}
	if !isSubject && node.Relation == "" {
		return domain.RequestBodyError{}
	}
	return ValidateReservedWord(node)
}

func ValidateReservedWord(st interface{}) error {
	value := reflect.ValueOf(st)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		str := field.Interface().(string)
		if strings.Contains(str, "%") {
			return domain.RequestBodyError{}
		}

	}
	return nil
}
