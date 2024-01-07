package domain

type Action string

const (
	CreateOperation Action = "create"
	DeleteOperation Action = "delete"
)

type Operation struct {
	Type     Action   `json:"action"`
	Relation Relation `json:"relation"`
}
