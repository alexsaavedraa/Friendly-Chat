import React from "react"
import {Component} from "react";
import { connect, sendMsg } from "./api/index.ts";
import styles from '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css';
import { 
  MainContainer, 
  ChatContainer, 
  ConversationHeader,
  MessageList, 
  Message, 
  MessageInput } 
  from '@chatscope/chat-ui-kit-react';

interface AppState {
  chatHistory: any[];
  user_id: string;
}

class App extends Component<{}, AppState> {
  constructor(props) {
    console.log(styles)
    super(props);
    this.state = {
      chatHistory: [],
      user_id: "Alex"
    }
  }

  componentDidMount() {
    connect((msg: any) => {
      let msg_data = JSON.parse(msg.data)
      let msg_model = {
        message: msg_data.body,
        sentTime: "Just now",
        sender: "someone",
        direction: "incoming"
      }
      this.setState(prevState => ({
        chatHistory: [...prevState.chatHistory, msg_model]
      }))
      console.log(this.state);
    });
  }

  send(msg: string) {
    sendMsg(msg);
  };

  render() {
    return (
      <div className="App">
        <MainContainer>
          <ChatContainer>       
          <ConversationHeader>
            <ConversationHeader.Content userName="Nimble Chat" info="Chat app" />     
          </ConversationHeader>
            <MessageList >
              {
                this.state.chatHistory.map((msg, i) => {
                  return (<Message key={i} model={msg} />)
                })
              }
              </MessageList>
            <MessageInput placeholder="Type message here" attachButton={false} onSend={this.send}/>        
          </ChatContainer>
        </MainContainer>
      </div>
    );
  };
};

export default App;