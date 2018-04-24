const gulp = require('gulp');
const webpack = require('webpack');
const gutil = require('gulp-util');
const gulpTSLint = require('gulp-tslint');
const tslint = require('tslint');
const del = require('del');
const prettier = require('prettier');

const production = gutil.env.production;
const webpackProdConfig = require('../webpack.config');
const webpackDevConfig = require('../webpack.config.dev');
const webpackConfig = production ? webpackProdConfig : webpackDevConfig;
const webpackCache = {};
webpackConfig.cache = webpackCache;

gulp.task('js:clean', () => del(['build/js/*']));

// webpack is instantiated outside the task for performance.
let webpackCompiler = webpack(webpackConfig);

gulp.task('js:webpack', (cb) => {
  webpackCompiler.run((err, stats) => {
    if (err) {
      throw new gutil.PluginError('webpack', err);
    }
    gutil.log(
      '[webpack]\n' +
        stats.toString({
          colors: true,
          cachedAssets: false,
          chunks: false,
        })
    );
    cb();
  });
});

gulp.task('js:lint', () => {
  let program = tslint.Linter.createProgram('tsconfig.json');
  return gulp
    .src(['src/js/**/*.ts', 'src/js/**/*.tsx'])
    .pipe(gulpTSLint({program, formatter: 'verbose'}))
    .pipe(gulpTSLint.report({emitError: false}));
});

gulp.task('js:fmt', () => {
  return gulp.src(['src/js/**/*.ts', 'src/js/**/*.tsx']).pipe((dest) => {
    console.log(dest);
    retirm
  });
});

gulp.task('js', ['js:clean', 'js:lint', 'js:webpack']);

gulp.task('js:watch', [], () =>
  gulp.watch('src/js/**/*.ts*', ['js:lint', 'js:webpack'])
);
