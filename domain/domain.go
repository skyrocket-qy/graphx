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
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Relation  string `json:"relation"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type RelationsResponse struct {
	Relations []Relation `json:"data"`
}

type StringsResponse struct {
	Data []string `json:"data"`
}
