import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import Room from './components/Room.js'
import { LoginBanner } from './components/Components'
import { connect } from 'react-redux'



class App extends Component {
  
  constructor(props) {
    super(props);
    this.state = {}
  }
  
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2><LoginBanner isLoggedIn={this.props.isLoggedIn} name={this.props.name}/></h2>
        </div>
        <Room isLoggedIn={this.props.isLoggedIn}/>
        <p className="App-intro">
          We have no idea what we are doing
        </p>
        <div>
          <img src="/dog.jpg" alt="We have no idea what we are doing"/>
        </div>
      </div>
    );
  }
}

export default connect(props => props)(App);
