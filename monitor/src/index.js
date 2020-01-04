import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App.js';

const body = document.getElementsByTagName('body')[0];
const root = document.createElement('div');
root.id = 'root';
body.appendChild(root);

ReactDOM.render(<App />, root);
