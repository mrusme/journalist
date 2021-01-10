![journalist](documentation/journalist.png)
-------------------------------------------

Journalist. A RSS aggregator.

[Download the latest version for macOS, Linux, FreeBSD, NetBSD, OpenBSD & Plan9 here](https://github.com/mrusme/journalist/releases/latest).

***WARNING: `journalist` is highly experimental software and not ready for use. 
Don't rely on it and expect changes in data structures, with no 
possibility to migrate existing data at the moment.***

## Build

```sh
make
```

**Info**: This will build using the version 0.0.0. You can prefix the `make` 
command with `VERSION=x.y.z` and set `x`, `y` and `z` accordingly if you want 
the version in `journalist --help` to be a different one.


## Usage

<iframe src="https://player.vimeo.com/video/498907228" width="640" height="400" frameborder="0" allow="autoplay; fullscreen" allowfullscreen></iframe>

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
for macOS/iOS). Simply specify `http://localhost:8000` (or the machine you're
running `journalist server` on) and either use `nobody` and `nobody` as 
credentials or – if you've subscribed to feeds using custom credentials – use
your own.
