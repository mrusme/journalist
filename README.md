![journalist](documentation/journalist.png)
-------------------------------------------

[![Release](https://github.com/mrusme/journalist/actions/workflows/release.yml/badge.svg)](https://github.com/mrusme/journalist/releases)
[![Docker](https://github.com/mrusme/journalist/actions/workflows/docker.yml/badge.svg)](https://hub.docker.com/r/mrusme/journalist)

Journalist. An RSS aggregator.


## What is `journalist`?

[Get more information here](https://マリウス.com/journalist-v1/).

## Build

```sh
make
```


## Usage

TODO

### Database

Journalist requires a database to store users and subscriptions. Supported
database types are SQLite, PostgreSQL and MySQL. The database can be configured
using the `JOURNALIST_DATABASE_TYPE` and `JOURNALIST_DATABASE_CONNECTION` env,
or the `Database.Type` and `Database.Connection` config properties.

#### Docker PostgreSQL Example

Run the database:

```sh
docker run -it --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=journalist \
  -p 127.0.0.1:5432:5432 \
  -d postgres:alpine
```

Configure `Database.Type` and `Database.Connection`:

```toml
[Database]
Type = "postgres"
Connection = "host=127.0.0.1 port=5432 dbname=journalist user=postgres password=postgres"
```


### Configuration

TODO

### Deployment


#### Docker

Official images are available on Docker Hub at 
[mrusme/journalist](https://hub.docker.com/r/mrusme/journalist) 
and can be pulled using the following command:

```sh
docker pull mrusme/journalist
```

GitHub release versions are available as Docker image tags (e.g. `0.0.1`). 
The `latest` image tag contains the latest code of the `master` branch.

It's possible to build journalist locally as a Docker container like this:

```sh
docker build -t journalist:latest . 
```

It can then be run using the following command:

```sh
docker run -it --rm --name journalist \
  -e JOURNALIST_... \
  -e JOURNALIST_... \
  -p 0.0.0.0:8000:8000 \
  journalist:latest
```


#### DigitalOcean App Platform

You can use the following App Spec to deploy `journalist` for as little as $12
per month on [DigitalOcean's App Platform](https://m.do.co/c/9d1b223a47bc). 
Fork this repo into your GitHub account, connect that with DO and 
replace `$$ACCOUNT$$` with your account's name in the App Spec:

```yaml
databases:
- engine: PG
  name: journalist
  num_nodes: 1
  size: db-s-dev-database
  version: "12"
name: journalist
region: nyc
services:
- dockerfile_path: Dockerfile
  envs:
  - key: JOURNALIST_DATABASE_TYPE
    scope: RUN_TIME
    value: "postgres"
  - key: JOURNALIST_DATABASE_CONNECTION
    scope: RUN_TIME
    value: "${journalist.DATABASE_URL}"
  github:
    branch: master
    deploy_on_push: true
    repo: $$ACCOUNT$$/journalist
  http_port: 8000
  instance_count: 1
  instance_size_slug: basic-xxs
  name: journalist
  routes:
  - path: /
```

