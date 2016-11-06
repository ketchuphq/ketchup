var gulp = require('gulp');
var sass = require('gulp-sass');
var gutil = require('gutil');
var sourcemaps = require('gulp-sourcemaps');

gulp.task('css', () =>
  gulp.src('./src/css/app.sass')
    .pipe(sourcemaps.init())
    .pipe(sass().on('error', gutil.log))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('./build'))
);

gulp.task('css:watch', () =>
  gulp.watch('./src/css/*.sass', ['css'])
);