const path = require('path');
const DirContentReplacementPlugin = require('./build-utils/dir-content-replacement-plugin');
const openSourceWebpackConfig = require('../../web/webpack.config.js');

module.exports = {
  ...openSourceWebpackConfig,
  plugins: [
    ...openSourceWebpackConfig.plugins,
    new DirContentReplacementPlugin({
      dir: path.resolve(__dirname, 'src'),
      mapper: filePath => filePath.replace('/enterprise', ''),
    }),
  ],
};
