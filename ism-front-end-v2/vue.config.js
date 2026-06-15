let path = require('path')
const webpack = require('webpack')
const JavaScriptObfuscator = require('webpack-obfuscator')
const ThemeColorReplacer = require('webpack-theme-color-replacer')
const {getThemeColors, modifyVars} = require('./src/utils/themeUtil')
const {resolveCss} = require('./src/utils/theme-color-replacer-extend')
const CompressionWebpackPlugin = require('compression-webpack-plugin')
const TerserWebpackPlugin = require('terser-webpack-plugin');
const productionGzipExtensions = ['js', 'css']
//const isProd = process.env.NODE_ENV === 'production'
const isProd = false
const BuildMod = false
const assetsCDN = {
  // webpack build externals
  externals: {
    vue: 'Vue',
    'vue-router': 'VueRouter',
    vuex: 'Vuex',
    axios: 'axios',
    nprogress: 'NProgress',
    clipboard: 'ClipboardJS',
    '@antv/data-set': 'DataSet',
    'js-cookie': 'Cookies'
  },
  css: [
  ],
  js: [
    '//cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.min.js',
    '//cdn.jsdelivr.net/npm/vue-router@3.3.4/dist/vue-router.min.js',
    '//cdn.jsdelivr.net/npm/vuex@3.4.0/dist/vuex.min.js',
    '//cdn.jsdelivr.net/npm/axios@0.19.2/dist/axios.min.js',
    '//cdn.jsdelivr.net/npm/nprogress@0.2.0/nprogress.min.js',
    '//cdn.jsdelivr.net/npm/clipboard@2.0.6/dist/clipboard.min.js',
    '//cdn.jsdelivr.net/npm/@antv/data-set@0.11.4/build/data-set.min.js',
    '//cdn.jsdelivr.net/npm/js-cookie@2.2.1/src/js.cookie.min.js'
  ]
}

module.exports = {
  parallel: false, // Windows 下多 worker 并发读文件会导致 EMFILE，强制关闭
  runtimeCompiler: true,
  devServer: {
    writeToDisk: true,
    watchOptions: {
      poll: 2000,        // 改大轮询间隔，进一步降低句柄占用
      aggregateTimeout: 500,
      ignored: [
        /node_modules/,
        /public[\\/]static[\\/]MapStorage/,
        /MapStorage/,
      ],
    },
    proxy: {
      '/api': { //此处要与 /services/api.js 中的 API_PROXY_PREFIX 值保持一致
        target: "http://127.0.0.1:8081/",
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      }
    },
    open: true, // 自动打开浏览器
    port: 7080 // 端口号
  },
  pluginOptions: {
    'style-resources-loader': {
      preProcessor: 'less',
      patterns: [path.resolve(__dirname, "./src/theme/theme.less")],
    }
  },
  configureWebpack: config => {
    config.entry.app = ["babel-polyfill", "whatwg-fetch", "./src/main.js"];
    config.performance = {
      hints: false
    }
    config.plugins.push(
        new ThemeColorReplacer({
          fileName: 'css/theme-colors-[contenthash:8].css',
          matchColors: getThemeColors(),
          injectCss: true,
          resolveCss
        })
    )
	config.optimization = {
        sideEffects: true,
        usedExports: true, // 开启 tree shaking
        concatenateModules:true,
        splitChunks: {
            chunks: "all",//async 按需加载的异步块 initial（初始块） all（所有块）
            minSize: 30000, // 模块的最小体积
            minChunks: 1, // 模块的最小被引用次数
            maxAsyncRequests: 5, // 按需加载的最大并行请求数
            maxInitialRequests: 3, // 一个入口最大并行请求数
            automaticNameDelimiter: '~', // 文件名的连接符
            name: true,
            cacheGroups: { // 缓存组
                vendor: {
                    test: /[\\/]node_modules[\\/]/,
                    name: "vendor",
                    priority: -10,
                    enforce: true
                }
            }
        },
		minimize: BuildMod,
		minimizer: [
            new TerserWebpackPlugin({
              parallel: true,
              terserOptions: {
                compress: {
                  drop_console: true,
                },
                output: {
                  comments: true,
                },
              },
              extractComments: false,
            })
		],
	}
    if(BuildMod) {
      config.plugins.push(
          new JavaScriptObfuscator({
            compact: true,
            controlFlowFlattening: false,
            deadCodeInjection: false,
            debugProtection: false,
            debugProtectionInterval: 0,
            disableConsoleOutput: true,
            identifierNamesGenerator: 'hexadecimal',
            log: false,
            numbersToExpressions: false,
            renameGlobals: false,
            selfDefending: true,
            simplify: true,
            splitStrings: false,
            stringArray: true,
            stringArrayCallsTransform: false,
            stringArrayEncoding: [],
            stringArrayIndexShift: true,
            stringArrayRotate: true,
            stringArrayShuffle: true,
            stringArrayWrappersCount: 1,
            stringArrayWrappersChainedCalls: true,
            stringArrayWrappersParametersMaxCount: 2,
            stringArrayWrappersType: 'variable',
            stringArrayThreshold: 0.75,
            unicodeEscapeSequence: false
          }, [  '**/js/chunk-**.**.js',
              '**/js/runtime.**.js',
              '**/js/vendor.**.js']),
      )
    }
    // Ignore all locale files of moment.js
    config.plugins.push(new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/))
    // 生产环境下将资源压缩成gzip格式
    if (isProd) {
      // add `CompressionWebpack` plugin to webpack plugins
      config.plugins.push(new CompressionWebpackPlugin({
        algorithm: 'gzip',
        test: new RegExp('\\.(' + productionGzipExtensions.join('|') + ')$'),
        threshold: 0,
        minRatio: 0.8
      }))
    }
    // if prod, add externals
    if (isProd) {
      config.externals = assetsCDN.externals
    }
  },
  chainWebpack: config => {
    // 排除 MapStorage — 防止 Windows EMFILE（Too many open files）
    // Vue CLI copy 插件：args[0] 是数组 [{ from, to, ignore }]
    try {
      config.plugin('copy').tap(args => {
        const patterns = args[0]  // ← 是数组，不是对象
        if (!Array.isArray(patterns)) return args
        patterns.forEach(p => {
          if (!p || typeof p !== 'object') return
          if (!p.ignore) p.ignore = []
          if (!p.ignore.includes('MapStorage/**')) {
            p.ignore.push('MapStorage/**')
          }
        })
        return args
      })
    } catch (e) {
      console.warn('[vue.config] 配置 copy 插件忽略 MapStorage 失败:', e.message)
    }

    // 生产环境下关闭css压缩的 colormin 项，因为此项优化与主题色替换功能冲突
    if (isProd) {
      config.plugin('optimize-css')
          .tap(args => {
            args[0].cssnanoOptions.preset[1].colormin = false
            return args
          })
    }
    // 生产环境下使用CDN
    if (isProd) {
      config.plugin('html')
          .tap(args => {
            args[0].cdn = assetsCDN
            return args
          })
    }
      config.plugins.delete('prefetch') // 关闭预加载
      config.plugins.delete('preload')
      // 彻底禁用 copy 插件，防止 Windows EMFILE（MapStorage 文件太多）
      // macOS/其他平台不需要，注释掉即可恢复 public/ → dist/ 复制
      // config.plugins.delete('copy')
  },
  css: {
    loaderOptions: {
      less: {
        lessOptions: {
          modifyVars: modifyVars(),
          javascriptEnabled: true
        }
      }
    }
  },
  publicPath: process.env.VUE_APP_PUBLIC_PATH,
  outputDir: 'dist',
  assetsDir: 'static',
  productionSourceMap: false,
  lintOnSave: false
}