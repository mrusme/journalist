![journalist](documentation/journalist.png)
-------------------------------------------

[![Release](https://github.com/mrusme/journalist/actions/workflows/release.yml/badge.svg)](https://github.com/mrusme/journalist/releases)
[![Docker](https://github.com/mrusme/journalist/actions/workflows/docker.yml/badge.svg)](https://hub.docker.com/r/mrusme/journalist)

Journalist. An RSS aggregator.


## What is `journalist`?

Journalist is an RSS aggregator that can sync subscriptions and read/unread
items across multiple clients without requiring a special client-side
integration. Clients can use Journalist by simply subscribing to its
personalized RSS feed.

Journalist aims to become a self-hosted alternative to services like Feedly,
Feedbin and others. It aims to offer a similar set of features like FreshRSS,
NewsBlur and Miniflux while being easier to set up/maintain and overall more
lightweight.

Find out more about Journalist [here](https://マリウス.com/journalist-v1/).


## Development

### Build

You can build Journalist yourself simply by running `make` in the repository
folder:

```sh
make
```

This will build a binary called `journalist`.


## Usage

Journalist ist a single binary service can be run on any Linux/Unix machine
by setting the required configuration values and launching the `journalist`
program.

### Configuration

Journalist will read its config either from a file or from environment
variables. Every configuration key available in the
example [`journalist.toml`](examples/etc/journalist.toml) can be exported as
environment variable, by separating scopes using `_` and prepend `JOURNALIST` to
it. For example, the following configuration:

```toml
[Server]
BindIP = "0.0.0.0"
```

... can also be specified as an environment variable:

```sh
export JOURNALIST_SERVER_BINDIP="0.0.0.0"
```

Journalist will try to read the `journalist.toml` file from one of the following
paths:

- `/etc/journalist.toml`
- `$XDG_CONFIG_HOME/journalist.toml`
- `$HOME/.config/journalist.toml`
- `$HOME/journalist.toml`
- `$PWD/journalist.toml`


### Database

Journalist requires a database to store users and subscriptions. Supported
database types are SQLite, PostgreSQL and MySQL. The database can be configured
using the `JOURNALIST_DATABASE_TYPE` and `JOURNALIST_DATABASE_CONNECTION` env,
or the `Database.Type` and `Database.Connection` config properties.

**WARNING:** If you do not specify a database configuration, Journalist will use
an in-memory SQLite database! As soon as Journalist shuts down, all data
inside the in-memory database is gone!


#### SQLite File Example

```toml
[Database]
Type = "sqlite3"
Connection = "file:my-database.sqlite?cache=shared&_fk=1"
```


#### PostgreSQL Example *(using Docker for PostgreSQL)*

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


#### MySQL Example

```toml
[Database]
Type = "mysql"
Connection = "mysqluser:mysqlpassword@tcp(mysqlhost:port)/database?parseTime=true"
```


### Deployment

#### Custom

All that's needed is a [configuration](#configuration) and Journalist can be
launched by e.g. running `./journalist` in a terminal.


#### Supervisor

To run Journalist via `supervisord`, create a config like this inside
`/etc/supervisord.conf` or `/etc/supervisor/conf.d/journalist.conf`:

```ini
[program:journalist]
command=/path/to/binary/of/journalist
process_name=%(program_name)s
numprocs=1
directory=/home/journalist
autostart=true
autorestart=unexpected
startsecs=10
startretries=3
exitcodes=0
stopsignal=TERM
stopwaitsecs=10
user=journalist
redirect_stderr=false
stdout_logfile=/var/log/journalist.out.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stdout_events_enabled=false
stderr_logfile=/var/log/journalist.err.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
stderr_events_enabled=false
```

**Note:** It is advisable to run Journalist under its own, dedicated daemon
user (`journalist` in this example), so make sure to either adjust `directory`
as well as `user` or create a user called `journalist`.


#### OpenBSD rc

As before, create a configuration file under `/etc/journalist.toml`.

Then copy the [example rc.d script](examples/etc/rc.d/journalist) to
`/etc/rc.d/journalist` and copy the binary to e.g.
`/usr/local/bin/journalist`. Last but not least, update the `/etc/rc.conf.local`
file to contain the following line:

```conf
journalist_user="_journalist"
```

It is advisable to run journalist as a dedicated user, hence create the
`_journalist` daemon account or adjust the line above according to your setup.

You can now run Journalist by enabling and starting the service:

```sh
rcctl enable journalist
rcctl start journalist
```


#### systemd

TODO


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

Alternatively a configuration TOML can be passed into the container like so:

```sh
docker run -it --rm --name journalist \
  -v /path/to/my/local/journalist.toml:/etc/journalist.toml \
  -p 0.0.0.0:8000:8000 \
  journalist:latest
```


#### Kubernetes

TODO


#### Render

Fork this repo into your GitHub account, adjust the
[`render.yaml`](render.yaml) accordingly and connect the forked repo [on
Render](https://dashboard.render.com/select-repo?type=blueprint).

Alternatively, you can also directly connect this public repo.


#### Heroku

TODO


#### DigitalOcean App Platform

Journalist can be deployed on [DigitalOcean's App
Platform](https://m.do.co/c/9d48825ddae1) for as little as ~$12 per month.

Fork this repo into your GitHub account, adjust the
[`.do/app.yaml`](.do/app.yaml) accordingly and connect the forked repo [on
DigitalOcean](https://cloud.digitalocean.com/apps/new).


#### DigitalOcean Function

Available soon.


#### Aamazon Web Services Lambda Function

TODO


#### Google Cloud Function

TODO


## API

TODO

