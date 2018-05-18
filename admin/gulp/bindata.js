const gulp = require('gulp');
const shell = require('gulp-shell');

const buildPaths = ['build', 'build/vendor', 'build/js', 'build/css', 'build/images'];
const buildCmd = `$(go env GOPATH)/bin/go-bindata -pkg admin -prefix build ${buildPaths.join(' ')}`;
const bindata = () => gulp.src('build/*', {read: false}).pipe(shell([buildCmd]));

gulp.task('bindata', ['css', 'js', 'html', 'images'], bindata);
gulp.task('bindata:partial', ['css', 'js:lint', 'js:webpack', 'html', 'images'], bindata);

gulp.task('bindata:watch', ['bindata'], () => gulp.watch('./src/**', ['bindata:partial']));
