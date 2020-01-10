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

/* eslint-disable no-console  */
const express = require('express');
const path = require('path');

// construct a mini server
const app = express();

// allow static assets to be served from the /dist folder
app.use(express.static(path.join(path.resolve(), 'dist')));

// Instantly reply to health checks from our ALB
app.get('/healthcheck', (req, res) => {
  res.sendStatus(200);
});

// Resolve all other requests to the index.html file
app.get('*', (req, res) => {
  res.sendFile(path.join(path.resolve(), 'dist/index.html'));
});

// initialize server
app.listen(process.env.SERVER_PORT, () => {
  console.log(`Listening on port ${process.env.SERVER_PORT}`);
});
