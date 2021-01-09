![journalist](documentation/journalist.png)
-------------------------------------------

Journalist. A RSS aggregator.

[Download the latest version for macOS, Linux, FreeBSD, NetBSD, OpenBSD & Plan9 here](https://github.com/mrusme/journalist/releases/latest).


## Build

```sh
make
```

**Info**: This will build using the version 0.0.0. You can prefix the `make` 
command with `VERSION=x.y.z` and set `x`, `y` and `z` accordingly if you want 
the version in `journalist --help` to be a different one.


## Usage

Please make sure to 
`export JOURNALIST_DB="postgres://postgres:postgres@127.0.0.1:5432/journalist"` 
or however your PostgreSQL connection string might look.

You can change the log level using the `JOURNALIST_LOG_LEVEL` env variable,
e.g. `JOURNALIST_LOG_LEVEL=10` for `debug`. By default, the level is set to 
`warn`.
