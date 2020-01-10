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

// Generate the codegen.json for gql-gen to dynamically generate .tsx files base on graphql file names

const scalars = {
  AWSEmail: 'string',
  AWSPhone: 'string',
  AWSTimestamp: 'number',
  AWSDateTime: 'string',
  AWSJSON: 'string',
};

const glob = require('glob');
const fs = require('fs');
const shell = require('shelljs');
const rimraf = require('rimraf');

const generateStruct = {};
const files = glob.sync('src/**/*.graphql');

// Delete _generated_ folder to start fresh
rimraf.sync('__generated__');

files.forEach(filePath => {
  const tempFilePath = filePath.split('/');

  tempFilePath.shift(); // remove "src"
  tempFilePath.shift(); // remove "graphql"

  tempFilePath.unshift('__generated__');

  const outputPath = tempFilePath.join('/').replace('.graphql', '.tsx');

  tempFilePath.pop(); // Remove filename
  // Recursively create folders
  shell.mkdir('-p', tempFilePath.join('/'));

  generateStruct[outputPath] = {
    documents: filePath,
    plugins: [{ typescript: { scalars } }],
  };
});

generateStruct['__generated__/schema.tsx'] = {
  plugins: [{ typescript: { scalars } }],
};

const codegenStruct = {
  schema: '../api/graphql/schema.graphql',
  overwrite: true,
  generates: generateStruct,
};

fs.writeFileSync('codegen.json', JSON.stringify(codegenStruct, null, 2));
