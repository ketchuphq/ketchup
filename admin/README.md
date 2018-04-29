# Admin Frontend

This module serves static assets. The static assets are built using Typescript
and SASS into the `./build` folder, and then `go-bindata` is used to generate a
Go file with all assets embedded into it. The Naga module in this folder is
responsible for mapping HTTP requests to the embedded files.

## Development

The frontend is built in React, with Yarn for dependency management, Jest for tests, and Gulp for build tooling.

Install all dependencies with `yarn`. Then, refer to the following cheat sheet for available commands.

```
npx gulp          # compile css, js, html, and bindata
npx gulp watch    # as above, but re-run when files change
npx gulp css      # compile sass to css
npx gulp js       # compile typescript to javascript
npx gulp html     # copy html and images to build dir
npx jest          # run jest tests
npx jest --watch  # as above, but re-run when files change
npm run fmt       # run prettier to reformat code
```
