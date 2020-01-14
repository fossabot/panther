/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/* eslint-disable no-param-reassign */
const fs = require('fs');
const klawSync = require('klaw-sync');

class DirContentReplacementPlugin {
  constructor({ dir, mapper }) {
    this.dir = dir;
    this.mapper = mapper;
  }

  apply(compiler) {
    if (!fs.existsSync(this.dir)) {
      return;
    }

    const filePaths = klawSync(this.dir, {
      nodir: true,
      traverseAll: true,
    }).map(f => f.path);

    const originToDestinationFilePathMapping = {};
    filePaths.forEach(filePath => {
      originToDestinationFilePathMapping[this.mapper(filePath)] = filePath;
    });
    const originAbsFilePaths = Object.keys(originToDestinationFilePathMapping);

    compiler.hooks.normalModuleFactory.tap('PantherEnterpriseReplacementPlugin', nmf => {
      nmf.hooks.afterResolve.tap('PantherEnterpriseReplacementPlugin', result => {
        if (!result) {
          return undefined;
        }

        if (originAbsFilePaths.includes(result.resource)) {
          result.resource = originToDestinationFilePathMapping[result.resource];
        }

        return result;
      });
    });
  }
}

module.exports = DirContentReplacementPlugin;
