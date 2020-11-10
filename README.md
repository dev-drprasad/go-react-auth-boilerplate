## RepoBOOST

layout pattern: https://github.com/golang-standards/project-layout

### Local Development:

To install postres:

```
brew install postgresql

```

To run API server:

```
# By default runs at 8000 port
JWT_SECRET=abcdef go run cmd/*
```

To run web app:

```
cd web
npm i
npm start
```

### Deployment

Requirements:

- Docker 19: https://docs.docker.com/engine/install/
- Docker 1.27: Compose https://docs.docker.com/compose/install/

To start whole stack:

```
COMPOSE_PROJECT_NAME=repoboost DB_USER=repoboost  DB_NAME=repoboost DB_PASSWORD=repoboost JWT_SECRET=abcdef PORT=80  docker-compose  -f deployments/docker-compose.yml up -d
```

To stop whole stack:

```
COMPOSE_PROJECT_NAME=repoboost-test DB_USER=repoboost  DB_NAME=repoboost DB_PASSWORD=repoboost JWT_SECRET=abcdef PORT=80  docker-compose  -f deployments/docker-compose.yml stop
```
