var gulp = require('gulp');
var sass = require('gulp-sass');
var gutil = require('gutil');
var sourcemaps = require('gulp-sourcemaps');
var cleanCSS = require('gulp-clean-css');

gulp.task('css', () =>
  gulp.src('./src/css/app.sass')
    .pipe(sourcemaps.init())
    .pipe(sass().on('error', gutil.log))
    .pipe(cleanCSS({ debug: true }, function (details) {
      let percent = details.stats.minifiedSize / details.stats.originalSize
      gutil.log(`${details.name} compressed: ${(percent * 100).toFixed(2)}%`)
    }))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('./build'))
);

gulp.task('css:watch', () =>
  gulp.watch([
    './src/css/*.sass',
    './src/css/**/*.sass'
  ], ['css'])
);