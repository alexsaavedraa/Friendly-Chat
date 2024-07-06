import React from "react"
import { useState } from "react";
import ReactTimeAgo from 'react-time-ago'
import VoteButton from './VoteButton.tsx'
import { Divider } from "@mui/material";

const MessageBox = (props: any) => {
    const message = props.message.body;
    const messageID = props.message?.MessageID;
    const category = props.message.category;
    const username = props.message.username;
    const timestamp =  props.message.time ? new Date(props.message.time) : null;
    const messageScore = props.messageScore
    let initialVote:number;
    if (props.message.user_vote=="up") {initialVote=1}
    else if (props.message.user_vote=="down") {initialVote=-1}
    else  {initialVote=0}
    const [userVote, setUserVote] = useState(initialVote);

    enum eventTypes {
        msg  = "message",
        user   = "new_user",
        vote = "votes"
    }
    
    if (category==eventTypes.user) {
        return (
            <div>
            <div className="votingContainer">
                <p><b>{username}</b> has joined the chat...</p> 
            </div>
            <Divider/>
            </div>
        )
    } else if (category==eventTypes.msg) {
    return (
        <div>
            <div className="messageContainer">
                <div className="messageInfoContainer">
                    <h4>{username}</h4>
                    {timestamp ? (
                        <ReactTimeAgo date={(timestamp.getTime())} locale="en-US" timeStyle={"round-minute"} />
                    ) : (
                        <ReactTimeAgo date={Date.now()} locale="en-US" timeStyle={"round-minute"} />
                    )}                </div>
                <p>{message}</p>
                <div className="votingContainer">
                    <VoteButton voteType={"up"} 
                                scoreEffect={1} 
                                userVote={userVote}
                                setUserVote={setUserVote}
                                messageID={messageID}
                                />
                    <div>{messageScore}</div>
                    <VoteButton voteType={"down"} 
                                scoreEffect={-1}
                                userVote={userVote}
                                setUserVote={setUserVote}
                                messageID={messageID}
                                />

                </div>  
            </div>
            <Divider/>
        </div>
    )}
    return null
};

export default MessageBox