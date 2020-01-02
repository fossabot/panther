const { execSync } = require('child_process');

execSync('node_modules/webpack/bin/webpack.js', { stdio: 'inherit' });
