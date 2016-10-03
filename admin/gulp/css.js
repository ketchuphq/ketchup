var gulp = require('gulp');
var sass = require('gulp-sass');
var gutil = require('gutil');

gulp.task('css', () =>
  gulp.src('./src/css/app.sass')
    .pipe(sass().on('error', gutil.log))
    .pipe(gulp.dest('./build'))
);

gulp.task('css:watch', () =>
  gulp.watch('./src/css/*.sass', ['css'])
);