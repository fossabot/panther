/* eslint-disable no-console  */
const express = require('express');
const path = require('path');

// construct a mini server
const app = express();

// allow static assets to be served from the /dist folder
app.use(express.static(path.join(path.resolve(), 'dist')));

// resolve all requests to the index.html file
app.get('*', (req, res) => {
  res.sendFile(path.join(path.resolve(), 'dist/index.html'));
});

// initialize server
app.listen(process.env.SERVER_PORT, () => {
  console.log(`Listening on port ${process.env.SERVER_PORT}`);
});
