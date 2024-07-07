import React from "react"
import { connect, sendMsg } from "../api/index.ts";
import styles from "@chatscope/chat-ui-kit-styles/dist/default/styles.min.css";

import { 
  MainContainer, 
  ChatContainer, 
  MessageList, 
  MessageInput,
  Loader
} from '@chatscope/chat-ui-kit-react';
import MessageBox from "../components/MessageBox.tsx"
import { useNavigate } from "react-router-dom";
import { endpoint_base } from "../config.ts";

interface ChatPageState {
    chatHistory: any[]; 
    charsMessage: string
    chatMessageTooLong: boolean
    messageScores: any
    loadingHistory: boolean
  }

  interface ChatPageProps {
    loggedIn: boolean;
    setLoggedIn: Function; 
    navigate: any;
  }

  interface Message {
    type: number;
    category: string;
    username: string;
    body: string;
    time: string;
    MessageID: string;
    votes: string;
    user_vote: string | null;
  }
  

class ChatPage extends React.Component<ChatPageProps, ChatPageState> {
    constructor(props:any) {
        super(props);
        this.state = {
          chatHistory: [],
          charsMessage: "256 characters remaining",
          chatMessageTooLong: false,
          messageScores: {} ,
          loadingHistory: false
        }
        console.log(styles)
    }

    componentWillReceiveProps = (nextProps: any) => {
      if (nextProps.loggedIn !== this.props.loggedIn) {
        this.props.navigate("/login")
      }
    }

    componentDidMount() {
        if (!this.props.loggedIn) {
          this.props.navigate("/login")
        }
        this.fetchHistory()
        const userDataString = localStorage.getItem("user");
        const userData = userDataString ? JSON.parse(userDataString) : null;
        const { username, token } = userData || {};
        const connection_status = connect(username, token, (msg: any) => {
          console.log("New Message")
          let msg_data = JSON.parse(msg.data)
          console.log(msg_data)
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
        });
        console.log("the connection status is ", connection_status)
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
    async fetchHistory() {
        this.setState({loadingHistory: true})
        const userDataString = localStorage.getItem("user");
        const userData = userDataString ? JSON.parse(userDataString) : null;
        const { username, token } = userData || {};
        try {
          const response = await fetch(`${endpoint_base}/history?username=${username}&token=${token}`);
          if (!response.ok) {
            throw new Error('Failed to fetch history');
          }
          const history: Message[] = await response.json();
          const historyMessageScores = history.reduce((acc: any, message) => {
            acc[message.MessageID] = message.votes;
            return acc;
        }, {});
          this.setState(prevState => ({
            chatHistory: [...history],
            messageScores: {...this.state.messageScores, ...historyMessageScores},
            loadingHistory: false
          }))
        } catch (error) {
          console.error('Error fetching history:', error);
        }
      }

    

    render() {
        return (
            <div className="chatContainer" style={{display: "flex", justifyContent: "center", alignItems:"center", flexDirection: "column"}}>
                <MainContainer style={{height: "90vh", width: "50vw", minWidth: "450px"}}> 
                    <ChatContainer style={{overflow: "auto"}}>
                     
                    <MessageList >
                      <MessageList.Content >
                        {this.state.loadingHistory && 
                        <div className="loaderContainer"><Loader>Loading Chat History...</Loader> </div>}
                          {
                            this.state.chatHistory.map((msg: any, i:number) => {
                                return (
                                    <MessageBox key={i} 
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

export default function(props: any) {
  const navigate = useNavigate();
  return <ChatPage {...props} navigate={navigate} />;
}