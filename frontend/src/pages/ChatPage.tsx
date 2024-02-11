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

const ChatPage = (props) => {
    const { username } = props
    console.log(styles)
    const [chatHistory, setChatHistory] = useState<any>([
        {
            body: "test message that is so long that it ends up taking multiple lines. We just want to see if the line wrapping is going to end up changing or warping how the upvote and downvote buttons work. Ideally, the lines should warp and there should be a sweet little rectangle where the voting takes place.",
            timestamp: new Date(),
            username: "Bob",
        },
        {
            body: "Yet another message. Now this one is a bit shorter.",
            timestamp: new Date(),
            username: "Joe",
        }
    ])

    useEffect (() => {
        connect((msg: any) => {
        let msg_data = JSON.parse(msg.data)
        setChatHistory([...chatHistory, msg_data])
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
                    return (<MessageBox key={i} message={msg}/> )
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