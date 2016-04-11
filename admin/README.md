# Naga Frontend

This module serves static assets. The static assets are built using Typescript and SASS into the `./build` folder, and then `go-bindata` is used to generate a Go file with all assets embedded into it. The Naga module in this folder is responsible for mapping HTTP requests to the embedded files.

The Makefile contains rules for building the `bindata.go` file. It combines
sub-makefiles from the scripts folder.

The following are various configs files:

|file|description|
|----|-----------|
|`bower.json`|bower dependencies
|`package.json`|npm dependencies
|`tsconfig.json`|Typescript
|`tslint.json`|linter configs
|`typings.json`|Typescript external types

## Usage

```
make prepare # installs go-bindata, and node and bower deps.
make sources # builds all sources into ./build
make test    # builds tests into ./test and runs them
```
