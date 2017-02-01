var gulp = require('gulp');

gulp.task('html', () =>
  gulp.src('src/html/**/*.html')
    .pipe(gulp.dest('./build'))
);

gulp.task('images', () =>
  gulp.src('src/images/*.png')
    .pipe(gulp.dest('./build/images'))
);

gulp.task('html:watch', () =>
  gulp.watch([
    'src/html/**/*.html',
    'src/images/*.png'
  ], ['html', 'images'])
);