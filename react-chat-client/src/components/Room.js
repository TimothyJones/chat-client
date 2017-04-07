import React, { Component } from 'react';
import { Messages, UserList, MessageForm } from './Components'


class Room extends Component {

  render() {
    return (
      <div>
        <RoomBody/>
      </div>
    );
  }
}

function RoomBody(props) {
  if (props.loggedIn) {
        return (<div>
                <div style={{display: "flex", flexFlow: "row nowrap"}}>
                  <Messages/>
                  <UserList/>
                </div>
                <MessageForm/>
               </div>);
  } 
  return null
}

export default Room;
