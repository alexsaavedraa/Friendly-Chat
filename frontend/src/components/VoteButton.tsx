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
    setUserVote: Function
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
                                                 setUserVote,
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
                // remove vote
                setUserVote(0);
                payload.body = null
            } else if (userVote===0) {
                // add vote
                setUserVote(scoreEffect);
            } else if (userVote===-scoreEffect) {
                // switch vote
                setUserVote(scoreEffect);
            }
            sendMsg(JSON.stringify(payload))
            }}>
            {userVote==scoreEffect ? selectedIcon : unselectedIcon}
        </IconButton>
    )
})

export default VoteButton;