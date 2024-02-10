import React from "react"
import { useEffect, useState } from "react";
import { connect, sendMsg } from "../api/index.ts";
import styles from "@chatscope/chat-ui-kit-styles/dist/default/styles.min.css";
import { 
  MainContainer, 
  ChatContainer, 
  ConversationHeader,
  MessageList, 
  Message, 
  MessageInput } 
  from '@chatscope/chat-ui-kit-react';

const ChatPage = (props) => {
    const { username } = props
    console.log(styles)
    const [chatHistory, setChatHistory] = useState<any>([])

    useEffect (() => {
        connect((msg: any) => {
        let msg_data = JSON.parse(msg.data)
        let msg_model = {
            message: msg_data.body,
            sentTime: "Just now",
            sender: "someone",
            direction: "incoming"
        }
        setChatHistory([...chatHistory, msg_model])
        console.log(chatHistory);
        });
    })

  function handleSendMessage(msg: string) {
    sendMsg(msg);
  };

return (
    <div className="chatContainer" style={{display: "flex", justifyContent: "center", alignItems:"center", flexDirection: "column"}}>
        <div className={"titleContainer"}>
            <div>Nimble Chat</div>
        </div>
        <MainContainer style={{height: "90vh", width: "50vw", minWidth: "450px"}}>
            <ChatContainer style={{overflow: "auto"}}>       
            <MessageList >
                {
                chatHistory.map((msg: any, i: number) => {
                    return (<Message key={i} model={msg} />)
                })
                }
                </MessageList>
            <MessageInput placeholder="Type message here" attachButton={false} onSend={handleSendMessage}/>        
            </ChatContainer>
        </MainContainer>
    </div>
    );
};

export default ChatPage;