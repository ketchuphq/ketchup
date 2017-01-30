var gulp = require('gulp');
var sass = require('gulp-sass');
var gutil = require('gutil');
var sourcemaps = require('gulp-sourcemaps');
var cleanCSS = require('gulp-clean-css');

gulp.task('css:internal', () =>
  gulp.src([
    'node_modules/quill/dist/quill.snow.css',
    'node_modules/codemirror/lib/codemirror.css',
    'node_modules/codemirror/theme/elegant.css'
  ])
    .pipe(gulp.dest('./build/vendor/'))
);


gulp.task('css:sass', () =>
  gulp.src('./src/css/app.sass')
    .pipe(sourcemaps.init())
    .pipe(sass({ includePaths: ['./bower_components'] })
      .on('error', sass.logError))
    // .pipe(cleanCSS({ debug: true }, function (details) {
    //   let percent = details.stats.minifiedSize / details.stats.originalSize
    //   gutil.log(`${details.name} compressed: ${(percent * 100).toFixed(2)}%`)
    // }))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('./build/css'))
);

gulp.task('css', ['css:internal', 'css:sass'])

gulp.task('css:watch', () =>
  gulp.watch([
    './src/css/*.sass',
    './src/css/**/*.sass'
  ], ['css'])
);