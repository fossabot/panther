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
