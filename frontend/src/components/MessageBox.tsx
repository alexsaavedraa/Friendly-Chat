import React from "react"
import { useState, ReactNode } from "react";

import ThumbDownAltIcon from '@mui/icons-material/ThumbDownAlt';
import ThumbDownOffAltIcon from '@mui/icons-material/ThumbDownOffAlt';
import ThumbUpAltIcon from '@mui/icons-material/ThumbUpAlt';
import ThumbUpOffAltIcon from '@mui/icons-material/ThumbUpOffAlt';
import { IconButton, Divider } from "@mui/material";

const MessageBox = (props) => {
    const message = props.message
    const [messageScore, setMessageScore] = useState(0)
    const [userVote, setUserVote] = useState(0)

    interface VoteButtonProps {
        scoreEffect: number
        selectedIcon: ReactNode
        unselectedIcon: ReactNode
        selectedColor: string
      }

    const VoteButton: React.FC<VoteButtonProps> = (({scoreEffect, selectedIcon, unselectedIcon, selectedColor}) => {
        return (
            <IconButton sx={userVote==scoreEffect && { color: selectedColor }} onClick={() => {
                if (userVote===scoreEffect) {
                    setUserVote(0);
                    setMessageScore(messageScore - scoreEffect)
                } else if (userVote===0) {
                    setUserVote(scoreEffect);
                    setMessageScore(messageScore + scoreEffect)
                } else if (userVote===-scoreEffect) {
                    setUserVote(scoreEffect);
                    setMessageScore(messageScore + (2*scoreEffect))
                }
                }}>
                {userVote==scoreEffect ? selectedIcon : unselectedIcon}
            </IconButton>
        )
    })
  
    return (
        <div>
        <div className="messageContainer">
            <p>{message}</p>
            <div className="votingContainer">
                <VoteButton scoreEffect={1} selectedIcon={<ThumbUpAltIcon/>} unselectedIcon={<ThumbUpOffAltIcon/>} selectedColor={"green"}/>
                <div>{messageScore}</div>
                <VoteButton scoreEffect={-1} selectedIcon={<ThumbDownAltIcon/>} unselectedIcon={<ThumbDownOffAltIcon/>} selectedColor={"red"}/>
            </div>   
        </div>
        <Divider/>
        </div>
    )
};

export default MessageBox