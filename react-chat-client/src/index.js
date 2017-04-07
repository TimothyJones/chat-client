import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';
import { Provider } from 'react-redux'
import LoginReducer from './reducers/LoginReducer'
import { createStore } from 'redux'

ReactDOM.render(
	<Provider store={createStore(LoginReducer)}>
	  <App />
  </Provider>,
  document.getElementById('root')
);
