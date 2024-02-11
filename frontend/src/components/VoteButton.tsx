import React from "react"
import ThumbDownAltIcon from '@mui/icons-material/ThumbDownAlt';
import ThumbDownOffAltIcon from '@mui/icons-material/ThumbDownOffAlt';
import ThumbUpAltIcon from '@mui/icons-material/ThumbUpAlt';
import ThumbUpOffAltIcon from '@mui/icons-material/ThumbUpOffAlt';
import { IconButton } from "@mui/material";
import { sendMsg } from "../api/index.ts";

enum voteTypes {
    up  = "up",
    down   = "down",
}

interface VoteButtonProps {
    voteType: string
    scoreEffect: number
    userVote: number
    messageScore: number
    setUserVote: Function
    setMessageScore: Function
    messageID: string
  }

interface votePayload {
    category: "vote"
    MessageID: string, 
    body: null | voteTypes
}

const VoteButton: React.FC<VoteButtonProps> = (({voteType, 
                                                 scoreEffect,
                                                 userVote,
                                                 messageScore,
                                                 setUserVote,
                                                 setMessageScore,
                                                 messageID
                                                }) => {
    let selectedIcon;
    let unselectedIcon;
    let selectedColor;
    if (voteType === voteTypes.up) {
        selectedIcon=<ThumbUpAltIcon fontSize={"small"}/>
        unselectedIcon=<ThumbUpOffAltIcon fontSize={"small"}/>
        selectedColor="green"
    } else if (voteType === voteTypes.down) {
        selectedIcon=<ThumbDownAltIcon fontSize={"small"}/>
        unselectedIcon=<ThumbDownOffAltIcon fontSize={"small"}/>
        selectedColor="red"
    } else {
        return
    }
    
    return (
        <IconButton 
        sx={{ color: (userVote==scoreEffect && selectedColor), '&:hover': {color: selectedColor}}} 
        onClick={() => {
            let payload:votePayload = {category:"vote", MessageID: messageID, body: voteType}
            if (userVote===scoreEffect) {
                // current vote removed
                setUserVote(0);
                setMessageScore(messageScore - scoreEffect)
                payload.body = null
            } else if (userVote===0) {
                // new vote selected
                setUserVote(scoreEffect);
                setMessageScore(messageScore + scoreEffect)
            } else if (userVote===-scoreEffect) {
                // vote switch
                setUserVote(scoreEffect);
                setMessageScore(messageScore + (2*scoreEffect))
            }
            sendMsg(JSON.stringify(payload))
            }}>
            {userVote==scoreEffect ? selectedIcon : unselectedIcon}
        </IconButton>
    )
})

export default VoteButton;