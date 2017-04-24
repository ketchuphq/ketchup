# Ketchup

[Ketchup](https://ketchuphq.com) is a CMS that aims to be...

- Easy to install and run
- Pretty snappy
- Extensible for making custom builds
- Easy to create and use your own templates
- Automatic HTTPS via Let’s Encrypt

## Usage

If you’re interested in _using_ Ketchup, please check out the documentation at [ketchuphq.com/docs](https://ketchuphq.com/docs) (which is itself built with Ketchup &mdash; dogfood!).

If you want to download the latest release, visit the [Releases page](https://github.com/ketchuphq/ketchup/releases).

If you'd like to see some screenshots, see [here](https://ketchuphq.com/docs/screenshots).

This readme will focus on how to build and develop Ketchup.

## Development

### Getting started

These instructions assume you're using OS X, but they should work on Linux as well. The only exception to this is that Homebrew is used to install `protoc` on OS X. You won't need `protoc` unless you're planning to modify the API (see the Protobuf section below).


To get started, you’ll first need to install **Go** and **Node**. They are used to compile the backend and frontend respectively.

Then, run the following to download frontend and backend dependencies:

```
make prepare-admin prepare-vendor
```

### Compiling

Now you should be able to compile the `ketchup` binary:

```
make
```

This will first run the frontend compilation (a Gulp task), then embed the frontend assets into `admin/bindata.go`, and finally run `go build`. The result will be a `ketchup` binary in the top-level.


If you want to modify the API, see below for instructions for how to update the proto files and regenerate the corresponding Go structs and Typescript classes.

### Development

For ongoing development, it can be a hassle to keep recompiling. There's a `gulp watch` task which you can run in the `admin` folder to recompile Typescript and SASS and output `bindata.go` on changes to frontend code.

For backend code, there's `./scripts/dev-watch.sh`, which will start a Ketchup server and recompile/restart it on changes to `.go` files.

### Releasing new versions

This section is for completeness; you probably won't have to do this.

```
make release-nr # dry run, only outputs to ./dist

# tag a release, goreleaser uses the latest tag
git tag -a v0.2.0 && git push origin v0.2.0
make release
```

Ketchup uses [`goreleaser`](https://github.com/goreleaser/goreleaser) to create and release new builds. The `goreleaser.yml` config file is dynamically generated in order to interpolate `$GOPATH` into the config, which is used to remove the `$GOPATH` that would otherwise appear in stack traces.

The version of the release is read from the latest git tag.


## Protobuf API

Protobufs are used to describe the API as well as the serialization for data stored in BoltDB.

To change the API, you’ll need to install the protobuf compiler, `protoc`. If you have Homebrew, running `make prepare` will install it using `brew`, as well as a custom plugin for generating Go output.

### Workflow

You should never edit `*.pb.go` and `./admin/src/js/lib/api.ts` directly. Instead, you should edit the relevant `.proto` protobuf file, and then regenerate those files:

```
make protos
```

After regenerating, you can recompile the frontend and backend to view your changes.

## Extending Ketchup

This is a work in progress. The long-term goal is to make it easy to swap out, configure, or add modules, using a custom `main.go` file. The following features are intended to be pluggable:

- authentication
- database
- template rendering engine
- additional API endpoints
- admin/authoring interface

The module system is documented [here](https://github.com/octavore/naga).

## Architectural decisions

A brief run through of the things without which this project would not exist:

- **Go**, powerful in its simplicity.
- **Typescript**, types are great.
- **SASS**, who even writes CSS.
- **Mithril**, simple yet batteries-included.
- **Protobufs**, autogenerating Go structs and Typescript interfaces for the API is kind of magical ([shameless plug](https://github.com/octavore/pbts)). Version 2 is used because it is nice to be able to differentiate between missing fields and set fields.
- **BoltDB**, embeddable and needs no additional setup.

## License

[ASLv2](https://www.apache.org/licenses/LICENSE-2.0.html)
