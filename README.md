# antho

antho (short for _anthology_) is a library and package manager for [jsonnet](http://jsonnet.org/).

The goal is to provide a way to `package`, `distribute`, and `easily develop` libraries.

As convention we use the file `main.libsonnet` as our entry point. Anything that imports a library can assume that file exists.

## Commands

### Print the JSonnet Path

```bash
# Current directory
antho jpath

# Specific library
antho jpath LIBRARY
```

### Package a library

```bash
antho pack LIBRARY
```

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

Create binaries for the current os.

```bash
make dev
```

Create all binaries for all platforms.

```bash
make bin
```
