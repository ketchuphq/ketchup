var gulp = require('gulp');

require('./gulp/css.js');
require('./gulp/js.js');
require('./gulp/html.js');
require('./gulp/bindata.js');

gulp.task('default', ['css', 'js', 'html', 'bindata'])
gulp.task('watch', [
  'css', 'js', 'html', 'bindata',
  'css:watch', 'js:watch', 'html:watch', 'bindata:watch'
])
