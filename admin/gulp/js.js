var gulp = require('gulp');
var webpack = require('webpack');
var gulpWebpack = require('webpack-stream');
var gutil = require('gulp-util');
var tslint = require('gulp-tslint');
var tsc = require('gulp-typescript');
var mocha = require('gulp-mocha');
var sourcemaps = require('gulp-sourcemaps');
var del = require('del');

let production = gutil.env.production
let webpackProdConfig = require('../webpack.config')
let webpackDevConfig = require('../webpack.config.dev')
let webpackCache = {}
let webpackConfig = production ? webpackProdConfig : webpackDevConfig
webpackConfig.cache = webpackCache

gulp.task('js:clean', () => del(['build/js/*']));

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

gulp.task('js', ['js:clean', 'js:internal', 'js:lint', 'js:webpack'])

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
    .pipe(sourcemaps.init())
    .pipe(ts())
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('.test/'))
})

gulp.task('js:test:mocha', () => {
  process.env.NODE_PATH = '.test/js/';
  return gulp.src(['.test/test/chai.js', '.test/**/*.test.js'], { read: false })
    .pipe(mocha({
      compilers: ['js:source-map-support/register']
    }))
})

gulp.task('js:test', ['js:test:compile', 'js:test:mocha'])

gulp.task('js:test:watch', () =>
  gulp.watch('src/js/**/*.ts*', ['js:test:compile', 'js:test:mocha'])
);