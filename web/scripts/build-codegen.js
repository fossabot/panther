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
