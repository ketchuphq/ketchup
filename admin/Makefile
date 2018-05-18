# shell is set so new path is picked up
SHELL := /bin/bash
export PATH := ./node_modules/.bin:$(PATH)

sources:
	gulp

production:
	gulp --production

prepare:
	go get -u github.com/jteeuwen/go-bindata/...
	yarn install
