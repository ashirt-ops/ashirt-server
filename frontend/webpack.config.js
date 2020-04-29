const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const path = require('path')
const _ = require('lodash')

module.exports = (env, argv) => ({
  entry: _.pick({
    'main': './src/entrypoints/main.ts',
    'archive_viewer': './src/entrypoints/archive_viewer.ts',
  }, env),

  output: {
    path: {
      'main': path.resolve(__dirname, 'dist'),
      'archive_viewer': path.resolve(__dirname, '../backend/archive_templates/assets'),
    }[env],
    publicPath: {
      'main': '/assets/',
      'archive_viewer': 'assets/',
    }[env],
    chunkFilename: '[chunkhash].js',
  },

  module: {
    rules: [{
      test: /\.tsx?$/,
      loader: 'ts-loader',
      exclude: /node_modules/,
      // This option cannot be set in tsconfig because it breaks frontend unit tests
      options: {compilerOptions: {module: 'esnext'}},
    }, {
      test: /\.css$/,
      use: [styleLoader(env), 'css-loader'],
    }, {
      test: /\.styl/,
      use: [
        styleLoader(env),
        {loader: 'css-loader', options: {
          modules: {
            mode: 'local',
            localIdentName: argv.mode == 'production' ? '[hash:base64:5]' : '[path][name]__[local]',
          },
        }},
        'stylus-loader',
      ],
    }, {
      test: /\.(svg|png|ttf|woff)/,
      use: 'file-loader',
    }],
  },

  plugins: [
    new MiniCssExtractPlugin({
      filename: '[name].css',
      chunkFilename: '[chunkhash].css',
    }),
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.styl'],
    alias: {
      src: path.resolve(__dirname, 'src'),
    },
  },

  devtool: {
    'main': 'cheap-source-map',
    'archive_viewer': 'none',
  }[env],

  devServer: {
    historyApiFallback: true,
    proxy: {
      '/web': {target: process.env.WEB_BACKEND_ORIGIN},
    },
    headers: {
      // Set the same security headers for local development as in production to catch bugs
      // before deployment due to blocked content.
      //
      // This should be kept in sync with headers set in nginx.conf
      //
      // Only exceptions are `Strict-transport-security` and `Expect-CT`, which cannot be set
      // locally because we don't use TLS for local development
      'X-Frame-Options': 'DENY',
      'X-Content-Type-Options': 'nosniff',
      'Referrer-Policy': 'strict-origin-when-cross-origin',
      'Content-Security-Policy': [
        // Default to none
        "default-src 'none'",

        // These do not fallback to default-src
        "base-uri 'none'",
        "form-action 'none'",
        "frame-ancestors 'none'",
        "sandbox allow-scripts allow-same-origin allow-forms allow-popups",

        // Allow xhr/fonts/images/scripts/css from self
        "connect-src 'self'",
        "font-src 'self'",
        "img-src 'self' data:",
        "script-src 'self'",
        "style-src 'self' 'unsafe-inline'",
      ].join(';'),
    },
  }
})

// A Webpack style loader takes raw css and packages it up in a way that it can be loaded in the app.
// It must be the final step in any css rules (webpack loads rules from right to left).
//
// MiniCssExtractPlugin does this by writing content addressable .css files to disk and including the hash
// in the webpack bundle
//
// style-loader does this by including the raw css string in the javascript bundle and appending a `style` tag
// to the head of the document on load
//
// We don't use style-loader for main because extracting css will allow the browser to load styles in parallel
// with javascript, allowing the app to display sooner. Since the archive loads all assets off disk this is an
// optimization that doesn't need to happen in the archive viewer.
//
// MiniCssExtractPlugin doesn't work well for local archives (it doesn't handle relative public paths) so we
// switch to style-loader for local archives.
function styleLoader(env) {
  return {
    'main': MiniCssExtractPlugin.loader,
    'archive_viewer': 'style-loader',
  }[env]
}
