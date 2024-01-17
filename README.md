# zanazibar-dag

Zanzibar-dag is a DAG(directed acyclic graph) based with Google's Zanzibar format, which can be the infrastructure below permission app.

## Relation

The `Relation` struct represents a relationship like edge in DAG between objects and subjects. It is defined as follows:

```go
// This means: Subject has a relation on Object
type Relation struct {
    ObjectNamespace  string
    ObjectName       string 
    Relation         string 
    SubjectNamespace string 
    SubjectName      string 
    SubjectRelation  string 
}
```

## How to use

1. Run postgres on docker(without docker, see ./docker-compose.yaml to get config)

```bash
docker compose up -d postgres
```

2. Run the main server

```bash
go run .
```

## Example

[HRBAC](https://github.com/skyrocketOoO/hrbac/tree/main)

## Reserved words

%

## Development benchmark

[Link](https://docs.google.com/spreadsheets/d/1qZiRE_kkno1mM0LzWiUnvX4cuYQRnep2NcNb4fPud-k/edit#gid=0)
