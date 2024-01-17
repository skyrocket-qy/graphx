# zanazibar-dag

A permission architecture based on Zanzibar and DAG, which can be used to build lots of permission system

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
