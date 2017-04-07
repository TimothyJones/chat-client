import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import Room from './components/Room.js'
import {LoginBanner} from './components/Components'

class App extends Component {
  
  constructor(props) {
    super(props);
    this.join = this.join.bind(this)
    this.state = {}
  }
  
  join(name) {
    this.setState(prevState => {
      console.log("Button");
      return {name: name, isLoggedIn: true} 
    })
  }

  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2><LoginBanner join={this.join} isLoggedIn={this.state.isLoggedIn} name={this.state.name}/></h2>
        </div>
        <Room isLoggedIn={this.state.isLoggedIn}/>
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

export default App;
