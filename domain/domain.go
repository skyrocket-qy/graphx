package domain

type Relation struct {
	ObjectNamespace  string `json:"object_namespace"`
	ObjectName       string `json:"object_name"`
	Relation         string `json:"relation"`
	SubjectNamespace string `json:"subject_namespace"`
	SubjectName      string `json:"subject_name"`
	SubjectRelation  string `json:"subject_relation"`
}

type Node struct {
	Namespace string
	Name      string
	Relation  string
}

type ErrResponse struct {
	Error string `json:"error"`
}

type DataResponse struct {
	Data []Relation `json:"data"`
}
