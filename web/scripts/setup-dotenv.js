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

/**
 * Makes sure to load ENV vars from the corresponding dotenv file (based on a script param)
 */
const path = require('path');
const dotenv = require('dotenv');
const minimist = require('minimist');
const { spawn } = require('child_process');
const chalk = require('chalk');

const { environment, _: otherArgs } = minimist(process.argv.slice(2));
if (!environment) {
  throw new Error(
    chalk.red('No environment provided. Please add one through the "--environment" flag')
  );
}

dotenv.config({
  path: path.resolve(process.cwd(), `config/.env.${environment}`),
});

if (otherArgs.length) {
  const [command, ...commandArgs] = otherArgs;
  spawn(command, commandArgs, { stdio: 'inherit' });
}
