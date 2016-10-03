var gulp = require('gulp');

gulp.task('html', () =>
  gulp.src('src/html/**/*.html')
    .pipe(gulp.dest('./build'))
);