# Ketchup

[![Go Report Card](https://goreportcard.com/badge/github.com/ketchuphq/ketchup)](https://goreportcard.com/report/github.com/ketchuphq/ketchup)
[![Coverage Status](https://coveralls.io/repos/github/ketchuphq/ketchup/badge.svg?branch=master)](https://coveralls.io/github/ketchuphq/ketchup?branch=master)
[![CircleCI](https://circleci.com/gh/ketchuphq/ketchup.svg?style=svg)](https://circleci.com/gh/ketchuphq/ketchup)


[Ketchup](https://ketchuphq.com) is a developer-first CMS.

- **BYO (bring your own) version-controlled templates**. Develop your templates outside of Ketchup, and pull them in with built-in support for git.
- **Customizable editor**. Templates declare content fields for frontend to render, so your editorial team can make edits easily.
- **[First class JSON API](https://ketchuphq.com/docs/api)**. Use Ketchup as a headless CMS by fetching content programmatically to power your own static websites or single-page webapps.
- **Automatic HTTPS**. Let's Encrypt certificate provisioning and renewal built-in.
- **Easy to deploy**. One binary, database included.
- Other features include file uploads, Markdown and WYSIWYG editors, content preview, etc.

### Table of contents

- [Usage](#usage)
- [Development](#development)
  - [1. Getting started](#1-getting-started)
  - [2. Compiling](#2-compiling)
    - [Watching](#watching)
  - [3. Development](#3-development)
    - [Protobuf API](#protobuf-api)
  - [4. Releasing new versions](#4-releasing-new-versions)
- [Extending Ketchup](#extending-ketchup)
- [Architectural decisions](#architectural-decisions)
- [Changelog](#changelog)
  - [0.3.0 - Relaxed React](#030---relaxed-react)
  - [0.2.0 - Tranquil Themes](#020---tranquil-themes)
  - [0.1.0 - Initial Release!](#010---initial-release)
- [License](#license)

## Usage

- If you’re interested in _using_ Ketchup, please check out the documentation at [ketchuphq.com/docs](https://ketchuphq.com/docs) (which is itself built with Ketchup &mdash; dogfood!).
- If you want to download the latest release, visit the [Releases page](https://github.com/ketchuphq/ketchup/releases).
- If you want to learn how to develop your own theme, read the [Theme Documentation](https://ketchuphq.com/docs/themes).
- If you'd like to see some screenshots, see [here](https://ketchuphq.com/docs/screenshots).

This readme will focus on how to build and develop Ketchup.

## Development

### 1. Getting started

These instructions assume you're using OS X, but they should work on Linux as well. The only exception to this is that Homebrew is used to install `protoc` on OS X. You won't need `protoc` unless you're planning to modify the API (see the Protobuf section below).

To get started, you’ll first need to install the following dependencies, which are used to compile the backend and frontend respectively.

- **Go >=1.9**
- **Node >= 8** with [NPM](https://www.npmjs.com) and [Yarn](https://yarnpkg.com/)

Then, run the following to download frontend and backend dependencies:

```
make prepare
```

### 2. Compiling

Now you should be able to compile the `ketchup` binary:

```
make
```

This will first run the frontend compilation (a Gulp task), then embed the frontend assets into `admin/bindata.go`, and finally run `go build`. The result will be a `ketchup` binary in the top-level.

If you want to modify the API, see below for instructions for how to update the proto files and regenerate the corresponding Go structs and Typescript classes.

#### Watching

For ongoing development, it can be a hassle to keep recompiling. There's a `gulp watch` task which you can run in the `admin` folder to recompile Typescript and SASS and output `bindata.go` on changes to frontend code. See [admin/README.md](https://github.com/ketchuphq/ketchup/blob/master/admin/README.md) for more details.

For backend code, there's `./scripts/dev-watch.sh`, which will start a Ketchup server and recompile+restart it on changes to `.go` files.

### 3. Development

#### Protobuf API

Protobufs are used to describe the API as well as being the serialization for data stored in BoltDB.

To change the API, you’ll need to install the protobuf compiler, `protoc`. If you have Homebrew, running `make prepare-protos` will install it using `brew`, as well as a custom plugin for generating Go output.

You should never edit `*.pb.go` and `./admin/src/js/lib/api.ts` directly. Instead, you should edit the relevant `.proto` protobuf file, and then regenerate those files:

```
make protos
```

After regenerating, you can recompile the frontend and backend to view your changes.

### 4. Releasing new versions

This section is for completeness; you probably won't have to do this.

```
make release-nr # dry run, only outputs to ./dist

# tag a release, goreleaser uses the latest tag
git tag -a v0.3.0 && git push origin v0.3.0
make release
```

Ketchup uses [`goreleaser`](https://github.com/goreleaser/goreleaser) to create and release new builds. The `goreleaser.yml` config file is dynamically generated in order to interpolate `$GOPATH` into the config, which is used to remove the `$GOPATH` that would otherwise appear in stack traces.

The version of the release is read from the latest git tag.


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
- **React**, you already know what this is.
- **Protobufs**, autogenerating Go structs and Typescript interfaces for the API is kind of magical ([shameless plug](https://github.com/octavore/pbts)). Version 2 is used because it is nice to be able to differentiate between missing fields and set fields.
- **BoltDB**, embeddable and needs no additional setup.

## Changelog

### 0.3.0 - Relaxed React

- Migrated frontend to React.
- Added live preview when editing content.
- Added drag-and-drop file uploads.
- Added `users:list` CLI to list all users from command line.
- Added `/admin/logout` endpoint for logging out.
- Increased server test coverage significantly.
- Fixed [#11](https://github.com/ketchuphq/ketchup/issues/11)
- 0.3.1 - Fixed [blank compose page](https://github.com/ketchuphq/ketchup/issues/13)

### 0.2.0 - Tranquil Themes

**Major changes**

- Upgrade mithril to 1.1.
  - Upgrade component class hierarchy to conform to 1.1 style.
  - Kudos to mithril for [#1922](https://github.com/MithrilJS/mithril.js/pull/1922)!
- Lots of progress on theme installation
  - Install themes via git; update with git pull; and check latest git version.
- Rework theme store, now supports nested templates via `{{ template "helper.html"}}` (closes [#6](https://github.com/ketchuphq/ketchup/issues/6))
- Reworked rendering to add additional context in templates, e.g. `{{ .Page.Content }}`
  - `.Page` context: `Content`, `PublishedAt`, `Theme`, `Template`, `Route`
  - `.Site` context: `Pages` (sorted in published order).
  - Page functions from [sprig](https://github.com/Masterminds/sprig): `date`, `dateModify`, `dateInZone`, `now`, and also custom `dateParseMillis`.
- Add support for global data variables and specifying global data variables in top-level `placeholders` field in `theme.json`.
- Remove bower, update deps.

### 0.1.0 - Initial Release!

## License

[ASLv2](https://www.apache.org/licenses/LICENSE-2.0.html)
