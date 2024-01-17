# zanazibar-dag

Zanzibar-dag is a DAG(directed acyclic graph) based with Google's Zanzibar format, which can be the infrastructure below permission app.

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
