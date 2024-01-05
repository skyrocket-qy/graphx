package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"zanzibar-dag/domain"
)

func ValidateRelation(rel domain.Relation) error {
	fmt.Printf("%+v\n", rel)
	if rel.ObjectNamespace == "" || rel.ObjectName == "" || rel.Relation == "" ||
		rel.SubjectNamespace == "" || rel.SubjectName == "" {
		return errors.New("invalid relation: some attr can't be empty")
	}
	return ValidateReservedWord(rel)
}

func ValidateNode(node domain.Node, isSubject bool) error {
	if node.Namespace == "" || node.Name == "" {
		return errors.New("invalid node: some attr can't be empty")
	}
	if !isSubject && node.Relation == "" {
		return errors.New("invalid node: some attr can't be empty")
	}
	return ValidateReservedWord(node)
}

func ValidateReservedWord(st interface{}) error {
	value := reflect.ValueOf(st)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		str := field.Interface().(string)
		if strings.Contains(str, "%") {
			return errors.New("has reserved word")
		}

	}
	return nil
}
