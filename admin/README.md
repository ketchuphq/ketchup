# Admin Frontend

This module serves static assets. The static assets are built using Typescript
and SASS into the `./build` folder, and then `go-bindata` is used to generate a
Go file with all assets embedded into it. The Naga module in this folder is
responsible for mapping HTTP requests to the embedded files.

## Development

Note: make sure you have `./node_modules/.bin` on your path, otherwise you will
need to run `./node_modules/.bin/gulp` instead of `gulp` below:

```
gulp         # compile css, js, html, and bindata
gulp watch   # recompile css, js, html, and bindata on files changes
gulp css     # compile sass to css
gulp js      # compile typescript to javascript
gulp html    # copy html and images to build dir
```
