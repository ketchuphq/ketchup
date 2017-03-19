var gulp = require('gulp');
var webpack = require('webpack');
var gulpWebpack = require('webpack-stream');
var gutil = require('gutil');
var tslint = require('gulp-tslint');
var tsc = require('gulp-typescript');
var mocha = require('gulp-mocha')

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
    .pipe(gulp.dest('build/js/'))
)

gulp.task('js:lint', () =>
  gulp.src(['src/js/**/*.ts', 'src/js/**/*.tsx'])
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

gulp.task('js:test:compile', () => {
  let ts = tsc.createProject(
    'tsconfig.json', {
      module: 'commonjs',
      outDir: '.test/'
    }
  )
  return gulp.src('src/**/*.ts*')
    .pipe(ts())
    .pipe(gulp.dest('.test/'))
})

gulp.task('js:test:mocha', () => {
  process.env.NODE_PATH = '.test/js/';
  return gulp.src(['.test/test/chai.js', '.test/**/*.test.js'], { read: false })
    .pipe(mocha())
    .on('error', gutil.log)
})

gulp.task('js:test', ['js:test:compile', 'js:test:mocha'])

gulp.task('js:test:watch', () =>
  gulp.watch('src/js/**/*.ts*', ['js:test:compile', 'js:test:mocha'])
);