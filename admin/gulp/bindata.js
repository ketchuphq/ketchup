var gulp = require('gulp')
var shell = require('gulp-shell')

const buildPaths = [
  'build',
  'build/vendor',
  'build/js',
  'build/css',
  'build/images'
]
gulp.task('bindata', () =>
  gulp.src('build/*', { read: false })
    .pipe(shell([
      `go-bindata -pkg admin -prefix build ${buildPaths.join(' ')}`
    ], {
        env: {
          PATH: `${process.env.GOPATH}/bin`
        }
      }))
)

gulp.task('bindata:watch', () =>
  gulp.watch('./build/**', ['bindata'])
);