var gulp = require('gulp')
var shell = require('gulp-shell')

gulp.task('bindata', () =>
  gulp.src('build/*', { read: false })
    .pipe(shell([
      'go-bindata -pkg admin -prefix build build build/vendor build/js build/css'
    ], {
        env: {
          PATH: `${process.env.GOPATH}/bin`
        }
      }))
)

gulp.task('bindata:watch', () =>
  gulp.watch('./build/**', ['bindata'])
);