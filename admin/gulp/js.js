var gulp = require('gulp');
var webpack = require('webpack');
var gulpWebpack = require('webpack-stream');
var gutil = require('gutil');
var tslint = require('gulp-tslint');

let webpackConfig = {
  entry: 'app.ts',
  devtool: 'source-map',
  output: {
    filename: 'app.js',
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js'],
    modules: ['src/js', 'node_modules']
  },
  module: {
    rules: [
      { test: /\.js$/, loader: 'source-map-loader', enforce: 'pre' },
      { test: /\.tsx?$/, loader: 'ts-loader' }
    ],
  },
  externals: {
    'mithril': 'm',
    'quill': 'Quill',
    'CodeMirror': 'CodeMirror'
  }
}

gulp.task('js:internal', () =>
  gulp.src([
    'node_modules/mithril/mithril.min.js',
    'node_modules/quill/dist/quill.min.js',

    'node_modules/codemirror/lib/codemirror.js',
    'node_modules/codemirror/mode/gfm/gfm.js',
    'node_modules/codemirror/mode/markdown/markdown.js',
    'node_modules/codemirror/addon/display/placeholder.js',
    'node_modules/codemirror/addon/mode/overlay.js'
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
