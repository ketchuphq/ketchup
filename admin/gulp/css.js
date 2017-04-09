let gulp = require('gulp');
let sass = require('gulp-sass');
let gutil = require('gulp-util');
let sourcemaps = require('gulp-sourcemaps');
let cleanCSS = require('gulp-clean-css');
let del = require('del');
let production = gutil.env.production

gulp.task('css:clean', () => del(['build/css/*']));

gulp.task('css:internal', () =>
  gulp.src([
    'node_modules/quill/dist/quill.snow.css',
    'node_modules/codemirror/lib/codemirror.css',
    'node_modules/codemirror/theme/elegant.css'
  ])
    .pipe(gulp.dest('./build/vendor/'))
);


gulp.task('css:sass', () => {
  if (production) {
    return gulp.src('./src/css/app.sass')
      .pipe(sass({ includePaths: ['./bower_components'] })
        .on('error', sass.logError))
      .pipe(cleanCSS({ debug: true }, function (details) {
        let percent = details.stats.minifiedSize / details.stats.originalSize
        gutil.log(`${details.name} compressed: ${(percent * 100).toFixed(2)}%`)
      }))
      .pipe(gulp.dest('./build/css'))
  }

  return gulp.src('./src/css/app.sass')
    .pipe(sourcemaps.init())
    .pipe(sass({ includePaths: ['./bower_components'] })
        .on('error', sass.logError))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('./build/css'))
});

gulp.task('css', ['css:clean', 'css:internal', 'css:sass'])

gulp.task('css:watch', () =>
  gulp.watch([
    './src/css/*.sass',
    './src/css/**/*.sass'
  ], ['css'])
);