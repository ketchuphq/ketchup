// https://github.com/Microsoft/TypeScript/issues/10178
// http://www.jbrantly.com/es6-modules-with-typescript-and-webpack/
// https://stackoverflow.com/questions/38168733/typescript-export-all-functions-in-a-namespace
// reexport codemirror type from namespace
import * as CM from 'CodeMirror';
declare global {
  var CodeMirror: typeof CM;
}
