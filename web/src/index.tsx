import React from 'react';
import ReactDOM from 'react-dom';
import './config';
import history from './history';
import App from './app';

ReactDOM.render(<App history={history} />, document.getElementById('root'));
