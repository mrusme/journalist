![journalist](documentation/journalist.png)
-------------------------------------------

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/mrusme/journalist)

Journalist. An RSS aggregator.

[Download the latest version for macOS, Linux, FreeBSD, NetBSD, OpenBSD & Plan9 here](https://github.com/mrusme/journalist/releases/latest).

***WARNING: `journalist` is highly experimental software and not ready for use. 
Don't rely on it and expect changes in data structures, with no 
possibility to migrate existing data at the moment.***

## What is `journalist`?

[Get more information here](https://マリウス.com/journalist-an-rss-aggregator/).

## Repository

This repository contains the source code of `journalist`. The code is being 
actively developed in the `develop` branch and only merged into `master` and 
tagged with a version as soon as it's stable enough for a release. 

If you intend to create **PRs**, please do so **against develop**.

## Build

```sh
make
```

**Info**: This will build using the version 0.0.0. You can prefix the `make` 
command with `VERSION=x.y.z` and set `x`, `y` and `z` accordingly if you want 
the version in `journalist --help` to be a different one.


## Usage

<iframe src="https://player.vimeo.com/video/498907228" width="640" height="400" frameborder="0" allow="autoplay; fullscreen" allowfullscreen></iframe>
 \
 

Please make sure to 
`export JOURNALIST_DB="postgres://postgres:postgres@127.0.0.1:5432/journalist"` 
or however your PostgreSQL connection string might look.

You can change the log level using the `JOURNALIST_LOG_LEVEL` env variable,
e.g. `JOURNALIST_LOG_LEVEL=10` for `debug`. By default, the level is set to 
`warn`.

### Database

Journalist requires you to have your own PostgreSQL database running somewhere.
Running it can be as easy as this, in case you're using Docker:

```sh
docker run -it --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=journalist \
  -p 127.0.0.1:5432:5432 \
  -d postgres:alpine
```

There are also plenty of cloud platforms where you can easily launch a fully-
managed PostgreSQL instance. For example on 
[DigitalOcean](https://m.do.co/c/9d1b223a47bc) you can get a PostgreSQL single
node cluster for as little as $15 per month.


### CLI

The `journalist` binary is a daemon as well as a CLI client for managing 
subscriptions and service configuration. There is no way (yet) to manage
subscriptions through the Fever API / a RSS client connecting to it.


#### Listing subscriptions

You can list all feeds a user is subscribed to using the `subscriptions`
command:

```sh
JOURNALIST_LOG_LEVEL=10 \
JOURNALIST_DB=postgres://postgres:postgres@127.0.0.1:5432/journalist \
journalist subscriptions
```

By default this command would list subscriptions for the user `nobody` 
(password: `nobody`). It's possible to specify `-u` (username) and `-p` 
(password) flags in order to list subscriptions for an individual 
account.


#### Subscribing to a feed

You can subscribe to a feed by using the `subscribe` command:

```sh
JOURNALIST_LOG_LEVEL=10 \
JOURNALIST_DB=postgres://postgres:postgres@127.0.0.1:5432/journalist \
journalist subscribe https://xn--gckvb8fzb.com/index.xml -g "Cool People"
```

`-g` adds the feed to a custom group, in this case `Cool People`. Groups are 
automatically created when specified via `-g`.

By default this command would subscribe as the user `nobody` 
(password: `nobody`). It's possible to specify `-u` (username) and `-p` 
(password) flags in order to subscribe to a feed under an individual account.


#### Unsubscribing from a feed

You can unsubscribe from a feed by using the `unsubscribe` command:

```sh
JOURNALIST_LOG_LEVEL=10 \
JOURNALIST_DB=postgres://postgres:postgres@127.0.0.1:5432/journalist \
journalist unsubscribe https://xn--gckvb8fzb.com/index.xml
```

If the feed was the last one in its group, the group is also being removed.

By default this command would unsubscribe as the user `nobody` 
(password: `nobody`). It's possible to specify `-u` (username) and `-p` 
(password) flags in order to unsubscribe from a feed under an individual 
account.


#### Running the RSS server

In order to be able to connect using any Fever API capable client you'll need
to run `journalist` in server mode:

```sh
JOURNALIST_LOG_LEVEL=10 \
JOURNALIST_DB=postgres://postgres:postgres@127.0.0.1:5432/journalist \
journalist server
```

You can then connect to it using your favourite Fever API client (e.g. Reeder
for macOS/iOS). Simply specify `http://localhost:8000/fever/` (or the machine you're
running `journalist server` on) and either use `nobody` and `nobody` as 
credentials or – if you've subscribed to feeds using custom credentials – use
your own.

#### Environment Variables

**General (CLI & server)**
- `JOURNALIST_LOG_LEVEL`: The log level, `0` being the lowest, `10` the highest
- `JOURNALIST_DB`: The PostgreSQL connection string

**Server only**
- `JOURNALIST_SERVER_BINDIP`: The IP to bind the server to, default: `0.0.0.0`
- `JOURNALIST_SERVER_PORT`: The port the server should run on, default: `8000`
- `JOURNALIST_SERVER_REFRESH`: The refresh interval (in seconds) at which the 
  server should update subscriptions, default: `0` (disabled)
- `JOURNALIST_SERVER_API_FEVER`: The Fever API, boolean value, default: `true` 
  (enabled)
- `JOURNALIST_SERVER_API_GREADER`: The Google Reader API, boolean value, 
  default: `false` (disabled) *NOT YET AVAILABLE*

### Docker

Official images are available on Docker Hub at 
[mrusme/journalist](https://hub.docker.com/r/mrusme/journalist) 
and can be pulled using the following command:

```sh
docker pull mrusme/journalist
```

GitHub release versions are available as Docker image tags (e.g. `0.0.1`). 
The `latest` image tag contains the latest code of the `main` branch, while the
`develop` tag contains the latest code of the `develop` branch.

It's possible to build journalist locally as a Docker container like this:

```sh
docker build -t journalist:latest . 
```

It can then be run using the following command:

```sh
docker run -it --rm --name journalist \
  -e JOURNALIST_LOG_LEVEL=10 \
  -e JOURNALIST_DB="postgres://postgres:postgres@172.17.0.2:5432/journalist" \
  -p 0.0.0.0:8000:8000
  journalist:latest
```
