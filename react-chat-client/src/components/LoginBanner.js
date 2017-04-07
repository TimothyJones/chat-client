import React, { Component } from 'react';
import { LoginForm } from './Components';

class LoginBanner extends Component {
  constructor(props) {
    super(props);
  }

  render() {
  	console.log(this.props)
  	if(this.props.isLoggedIn) {
  		return <div>
  					Welcome, <span style={{fontStyle: "italic"}}>{this.props.name}</span>
  				</div>
  	} else {
		return <div>
				 <LoginForm join={this.props.join}/>
			   </div>
  	}
  }
}

export default LoginBanner;