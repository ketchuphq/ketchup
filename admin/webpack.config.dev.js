let webpack = require('webpack');
module.exports = {
  entry: {
    vendor: [
      'store/dist/store.modern',
      'date-fns/format',
      'lodash-es/debounce',
      'lodash-es/cloneDeep',
      'lodash-es/isEqual'
    ],
    app: 'app.ts'
  },
  devtool: 'source-map',
  output: {
    filename: '[name].js',
    publicPath: '/admin/js/'
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js'],
    modules: ['src/js', 'node_modules']
  },
  module: {
    rules: [
      { test: /\.js$/, loader: 'source-map-loader', enforce: 'pre' },
      { test: /\.tsx?$/, loader: 'ts-loader' },
      { test: /\.css$/, loader: 'style-loader!css-loader' }
    ],
  },
  externals: {
    'mithril': 'm',
  },
  plugins: [
    new webpack.optimize.CommonsChunkPlugin({
      name: 'vendor'
    })
  ],
  cache: {}
}