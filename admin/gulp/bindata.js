let gulp = require('gulp')
let shell = require('gulp-shell')

const buildPaths = [
  'build',
  'build/vendor',
  'build/js',
  'build/css',
  'build/images'
]
let bindata = () =>
  gulp.src('build/*', { read: false })
    .pipe(shell([
      `go-bindata -pkg admin -prefix build ${buildPaths.join(' ')}`
    ], {
        env: {
          PATH: `${process.env.GOPATH}/bin`
        }
      }))

gulp.task('bindata', ['css', 'js', 'html', 'images'], bindata)
gulp.task('bindata:partial', ['css', 'js:internal', 'js:lint', 'js:webpack', 'html', 'images'], bindata)

gulp.task('bindata:watch', ['bindata'], () =>
  gulp.watch('./src/**', ['bindata:partial'])
)