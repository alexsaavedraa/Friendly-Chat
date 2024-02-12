import React from "react"
import { connect, sendMsg } from "../api/index.ts";
import styles from "@chatscope/chat-ui-kit-styles/dist/default/styles.min.css";
import { 
  MainContainer, 
  ChatContainer, 
  MessageList, 
  MessageInput
} from '@chatscope/chat-ui-kit-react';
import MessageBox from "../components/MessageBox.tsx"
import { useNavigate } from "react-router-dom";

interface ChatPageState {
    chatHistory: any[]; 
    charsMessage: string
    chatMessageTooLong: boolean
    messageScores: any
  }

  interface ChatPageProps {
    loggedIn: boolean;
    setLoggedIn: Function; 
    navigate: any;
  }

class ChatPage extends React.Component<ChatPageProps, ChatPageState> {
    constructor(props) {
        super(props);
        this.state = {
          chatHistory: [],
          charsMessage: "256 characters remaining",
          chatMessageTooLong: false,
          messageScores: {}        
        }
        const style = styles
    }

    componentWillReceiveProps = (nextProps) => {
      if (nextProps.loggedIn !== this.props.loggedIn) {
        this.props.navigate("/login")
      }
    }

    connectCallback = (msg:any) => {
          let msg_data = JSON.parse(msg.data)
          let msg_id = msg_data.MessageID;
          if (msg_data.category=="votes") {
            if (msg_id) {
              this.setState(prevState => ({
                messageScores: {...this.state.messageScores, [msg_id] : msg_data.body}
              }))
            }
          }
          else { 
            this.setState(prevState => ({
              chatHistory: [...this.state.chatHistory, msg_data],
              messageScores: {...this.state.messageScores, [msg_id] : "0"}
            }))
          }
        }


    componentDidMount() {
        if (this.props.loggedIn) {
          console.log("Chat Component logged in ", this.props.loggedIn)
          const userDataString = localStorage.getItem("user");
          const userData = userDataString ? JSON.parse(userDataString) : null;
          const { username, token } = userData || {};
          connect(username, token, this.connectCallback)
        }
      }

    handleSendMessage = (msg: string) => {
        sendMsg(JSON.stringify({"category": "message", "body": msg}));
        this.setState({
          charsMessage: `${256} characters remaining`,
          chatMessageTooLong: false
        })
    };

    handleInputChange = (currInput: string) => {
      if (256>=currInput.length) {
        this.setState({
          charsMessage: `${256-currInput.length} characters remaining`,
          chatMessageTooLong: false
        })
      } else {

        this.setState({
          charsMessage: `${currInput.length - 256} characters over limit`,
          chatMessageTooLong: true
        })
      }
    }

    

    render() {
        return (
            <div className="chatContainer" style={{display: "flex", justifyContent: "center", alignItems:"center", flexDirection: "column"}}>
                <MainContainer style={{height: "90vh", width: "50vw", minWidth: "450px"}}> 
                    <ChatContainer style={{overflow: "auto"}}>       
                    <MessageList >
                      <MessageList.Content >
                          {
                            this.state.chatHistory.map((msg: any) => {
                                return (
                                    <MessageBox key={msg.MessageID} 
                                                message={msg} 
                                                messageScore={this.state.messageScores[msg?.MessageID]}/> 
                                )
                            }) 
                          }
                        </MessageList.Content>
                    </MessageList>
                <MessageInput sendDisabled={this.state.chatMessageTooLong}
                              placeholder="Type message here" 
                              attachButton={false} 
                              onSend={this.handleSendMessage}
                              onChange={this.handleInputChange}
                              /> 
                </ChatContainer>
            </MainContainer>
            <div overflow-wrap={"break-word"} 
                 style={{color: this.state.chatMessageTooLong?"red":"gray"}}>
              {this.state.charsMessage}
            </div>
        </div>
        );
    };
};

export default function(props) {
  const navigate = useNavigate();
  return <ChatPage {...props} navigate={navigate} />;
}