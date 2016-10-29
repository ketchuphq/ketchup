var gulp = require('gulp');

require('./gulp/css.js');
require('./gulp/js.js');
require('./gulp/html.js');
require('./gulp/bindata.js');

gulp.task('default', ['css', 'js', 'html'])
gulp.task('watch', [
  'css', 'js', 'html',
  'css:watch', 'js:watch', 'html:watch', 'bindata:watch'
])
