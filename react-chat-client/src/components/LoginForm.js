import React, { Component } from 'react'
import { connect } from 'react-redux'


class LoginForm extends Component {

  constructor (props) {
    super(props)
    this.onJoinClick = this.onJoinClick.bind(this)
    this.checkName =   this.checkName.bind(this)

    this.state =       { valid: false, name: null }
  }

  onJoinClick (event) {
    event.preventDefault()
    if (!this.state.valid) {
      return
    }

    this.props.dispatch({type: "JOIN_ROOM", name: this.state.name})
    //this.props.join(this.state.name)
  }

  checkName (event) {
    const name = event.target.value
    const valid = name && name.length > 0
    this.setState({ valid, name })

    // if the enter key was pressed and the form is valid, submit it
    if (valid && event.type === 'keydown' && event.keyCode === 13) {
	    this.props.dispatch({type: "JOIN_ROOM", name: name})
     // this.props.join(name)
    }
  }


  render() {
    return (
      <div>
        <input type="text" name="name" placeholder="What is your name" onKeyDown={this.checkName}
          onChange={this.checkName}/>
        <button name="submit" onClick={this.onJoinClick}>Join</button>
      </div>
    );
  }
}

export default connect()(LoginForm);