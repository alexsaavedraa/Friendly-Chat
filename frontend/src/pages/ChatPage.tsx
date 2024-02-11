import React from "react"
import { useEffect, useState } from "react";
import { connect, sendMsg } from "../api/index.ts";
import styles from "@chatscope/chat-ui-kit-styles/dist/default/styles.min.css";
import { 
  MainContainer, 
  ChatContainer, 
  MessageList, 
  MessageInput
} from '@chatscope/chat-ui-kit-react';
import MessageBox from "../components/MessageBox.tsx"

interface ChatPageState {
    chatHistory: any[]; // Assuming chatHistory can be an array of any type
  }

class ChatPage extends React.Component<ChatPageState> {
    constructor(props) {
        super(props);
        this.state = {
          chatHistory: [],
        }
        console.log(styles)
    }


    componentDidMount() {
        const userDataString = localStorage.getItem("user");
        const userData = userDataString ? JSON.parse(userDataString) : null;
        const { username, token } = userData || {};
        connect(username, token, (msg: any) => {
          console.log("New Message")
          this.setState(prevState => ({
            chatHistory: [...this.state.chatHistory, JSON.parse(msg.data)]
          }))
          console.log(this.state.chatHistory);
        });
      }

  handleSendMessage(msg: string) {
    sendMsg(msg);
  };
render() {
return (
    <div className="chatContainer" style={{display: "flex", justifyContent: "center", alignItems:"center", flexDirection: "column"}}>
        <div className={"titleContainer"}>
            <div>Nimble Chat</div>
        </div>
        <MainContainer style={{height: "90vh", width: "50vw", minWidth: "450px"}}>
            <ChatContainer style={{overflow: "auto"}}>       
            <MessageList >
                {
                this.state.chatHistory.map((msg: any, i: number) => {
                    return (<MessageBox key={i} message={msg}/> )
                })
                }
                </MessageList>
            <MessageInput placeholder="Type message here" attachButton={false} onSend={this.handleSendMessage}/>        
            </ChatContainer>
        </MainContainer>
    </div>
    );
};
};

export default ChatPage;