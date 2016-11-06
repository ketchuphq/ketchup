var gulp = require('gulp');
var webpack = require('webpack');
var gutil = require('gutil');
var tslint = require('gulp-tslint');

gulp.task('js', (cb) => {
  webpack({
    entry: 'app.ts',
    devtool: 'source-map',
    output: {
      filename: 'app.js',
      path: 'build',
    },
    resolve: {
      extensions: ['', '.ts', '.js'],
      modulesDirectories: ['src/js', 'node_modules']
    },
    module: {
      loaders: [{ test: /\.ts$/, loader: 'ts-loader' }],
      preLoaders: [{ test: /\.js$/, loader: 'source-map-loader' }],
    },
    externals: {
      'mithril': 'm'
    }
  }, (err, stats) => {
    if (err) {
      throw new gutil.PluginError('webpack', err)
    }
    gutil.log('[webpack]', stats.toString({
      chunks: false,
      colors: true
    }))
    cb()
  })
});

gulp.task('js:lint', () =>
  gulp.src('src/js/**/*.ts')
    .pipe(tslint({
      formatter: 'verbose'
    }))
    .pipe(tslint.report({
      emitError: false
    }))
);

gulp.task('js:watch', () =>
  gulp.watch('src/js/**/*.ts', ['js:lint', 'js'])
);
