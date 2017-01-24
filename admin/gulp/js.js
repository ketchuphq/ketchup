var gulp = require('gulp');
var webpack = require('webpack');
var gulpWebpack = require('webpack-stream');
var gutil = require('gutil');
var tslint = require('gulp-tslint');

let webpackCache = {}
let webpackConfig = require('../webpack.config')

webpackConfig.cache = webpackCache

let prod = false
if (prod) {
  webpackConfig.plugins = [
    new webpack.optimize.UglifyJsPlugin({
      compress: { warnings: false }
    })
  ]
}

gulp.task('js:internal', () =>
  gulp.src([
    'node_modules/mithril/mithril.min.js',
  ])
    .pipe(gulp.dest('./build/vendor/'))
);

gulp.task('js:webpack', () =>
  gulp.src('src/app.ts')
    .pipe(gulpWebpack(webpackConfig, webpack))
    .on('error', function(err) { this.emit('end') })
    .pipe(gulp.dest('build/'))
)

gulp.task('js:lint', () =>
  gulp.src('src/js/**/*.ts')
    .pipe(tslint({
      formatter: 'verbose'
    }))
    .pipe(tslint.report({
      emitError: false
    }))
);

gulp.task('js', ['js:internal', 'js:lint', 'js:webpack'])

gulp.task('js:watch', () =>
  gulp.watch('src/js/**/*.ts*', ['js'])
);
