const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require("html-webpack-plugin");
const HtmlWebpackHarddiskPlugin = require("html-webpack-harddisk-plugin");
const path = require('path')

const miniCssExtraPluginConfig = {
  loader: MiniCssExtractPlugin.loader,
  options: {
    esModule: false,
  },
};

module.exports = (env, argv) => ({
  entry: './src/index.tsx',

  output: {
    path: path.resolve(__dirname, 'dist/assets'),
    publicPath: '/assets/',
    filename: 'main-[contenthash].js',
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
      use: [miniCssExtraPluginConfig, 'css-loader'],
    }, {
      test: /\.styl/,
      use: [miniCssExtraPluginConfig, {
        loader: 'css-loader', options: {
        modules: {
          mode: 'local',
          localIdentName: argv.mode == 'production' ? '[hash:base64:5]' : '[path][name]__[local]',
        },
      }}, 'stylus-loader'],
    }, {
      test: /\.(svg|png|ttf|woff)/,
      use: [
        {
          loader: 'file-loader',
          options: {
            esModule: false
          }
        }
      ]
    }],
  },

  plugins: [
    new HtmlWebpackPlugin({
      title: "ASHIRT",
      alwaysWriteToDisk: true,
    }),
    new HtmlWebpackHarddiskPlugin({
      outputPath: path.resolve(__dirname, 'public')
    }),
    new MiniCssExtractPlugin({
      filename: 'main-[contenthash].css',
      chunkFilename: '[chunkhash].css',
    }),
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.styl'],
    alias: {
      src: path.resolve(__dirname, 'src'),
    },
  },

  devtool: 'cheap-source-map',

  devServer: {
    historyApiFallback: true,
    publicPath: "/assets/",
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
