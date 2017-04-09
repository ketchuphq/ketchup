var gulp = require('gulp');

require('./gulp/css.js');
require('./gulp/js.js');
require('./gulp/html.js');
require('./gulp/bindata.js');

gulp.task('default', ['bindata'])
gulp.task('watch', ['bindata:watch'])
