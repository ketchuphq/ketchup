let gulp = require('gulp');
let webpack = require('webpack');
let gutil = require('gulp-util');
let gulpTSLint = require('gulp-tslint');
let tslint = require('tslint');
let tsc = require('gulp-typescript');
let mocha = require('gulp-mocha');
let sourcemaps = require('gulp-sourcemaps');
let del = require('del');

let production = gutil.env.production
let webpackProdConfig = require('../webpack.config')
let webpackDevConfig = require('../webpack.config.dev')
let webpackConfig = production ? webpackProdConfig : webpackDevConfig
let webpackCache = {}
webpackConfig.cache = webpackCache

gulp.task('js:clean', () => del(['build/js/*']));

gulp.task('js:internal', () =>
  gulp.src([
    'node_modules/mithril/mithril.min.js',
  ])
    .pipe(gulp.dest('./build/vendor/'))
);

// webpack is instantiated outside the task for performance.
let webpackCompiler = webpack(webpackConfig)

gulp.task('js:webpack', (cb) => {
  webpackCompiler.run((err, stats) => {
    if (err) {
      throw new gutil.PluginError('webpack', err);
    }
    gutil.log('[webpack]\n' + stats.toString({
      colors: true,
      cachedAssets: false,
      chunks: false,
    }))
    cb()
  })
})

gulp.task('js:lint', () => {
  let program = tslint.Linter.createProgram('tsconfig.json')
  return gulp.src(['src/js/**/*.ts', 'src/js/**/*.tsx'])
    .pipe(gulpTSLint({ program, formatter: 'verbose' }))
    .pipe(gulpTSLint.report({ emitError: false }))
});

gulp.task('js', ['js:clean', 'js:internal', 'js:lint', 'js:webpack'])

gulp.task('js:watch', ['js:internal'], () =>
  gulp.watch('src/js/**/*.ts*', ['js:lint', 'js:webpack'])
);

let ts = tsc.createProject(
  'tsconfig.json', {
    module: 'commonjs',
    outDir: '.test/'
  }
)

gulp.task('js:test:compile', () => {
  return gulp.src('src/**/*.ts*')
    .pipe(sourcemaps.init())
    .pipe(ts())
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('.test/'))
})

gulp.task('js:test:mocha', ['js:test:compile'], () => {
  process.env.NODE_PATH = '.test/js/';
  return gulp.src(['.test/test/chai.js', '.test/**/*.test.js'], { read: false })
    .pipe(mocha({
      compilers: ['js:source-map-support/register']
    }))
})

gulp.task('js:test', ['js:test:compile', 'js:test:mocha'])

gulp.task('js:test:watch', () =>
  gulp.watch('src/js/**/*.ts*', ['js:test:mocha'])
);