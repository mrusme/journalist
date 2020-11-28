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

Please make sure to `export JOURNALIST_DB=~/.config/journalist.db` (or whatever location 
you would like to have the journalist database at).
