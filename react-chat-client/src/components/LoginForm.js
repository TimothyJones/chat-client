import React, { Component } from 'react';

class LoginForm extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
        <input type="text" name="name" placeholder="What is your name"/>
        <button name="submit" onClick={this.props.join}>Join</button>
      </div>
    );
  }
}

export default LoginForm;