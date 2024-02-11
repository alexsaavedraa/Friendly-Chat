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
          chatMessageTooLong: false
        }
        console.log(styles)
    }

    componentWillReceiveProps = (nextProps) => {
      if (nextProps.loggedIn !== this.props.loggedIn) {
        this.props.navigate("/login")
      }
    }

    componentDidMount() {
        if (!this.props.loggedIn) {
          this.props.navigate("/login")
        }
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

    handleSendMessage = (msg: string) => {
        sendMsg(JSON.stringify({"category": "message", "body": msg}));
        this.setState({
          charsMessage: `${256} characters remaining`,
          chatMessageTooLong: false
        })
    };

    handleInputChange = (currInput: string) => {
      console.log(this.state.chatMessageTooLong)
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
                        this.state.chatHistory.map((msg: any, i: number) => {
                            return (
                                <MessageBox key={i} message={msg}/> 
                            )
                        }) }
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
            <div style={{color: this.state.chatMessageTooLong?"red":"gray"}}>
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