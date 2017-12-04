# antho

antho (short for anthology) is a library and package manager for [jsonnet](http://jsonnet.org/).

The goal is to provide a way to `package`, `distribute`, and `easily develop` libraries.

As convention we use the file `main.libsonnet` as our entry point. Anything that imports a library can assume that file exists.

## Building

We rely on two tools for the build process:

* [glide](https://github.com/Masterminds/glide) - Package management
* [gox](https://github.com/mitchellh/gox) - Cross compilation

### Requirements

```bash
go get github.com/mitchellh/gox
brew install glide
```

### Tests

```bash
make
```

### Binaries

```bash
make bin
```
