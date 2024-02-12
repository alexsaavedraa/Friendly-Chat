import React from "react"
import { useNavigate } from "react-router-dom";

const HomePage = (props) => {
    const { username, loggedIn } = props
    const navigate = useNavigate();
    console.log(username)
    
    const onButtonClick = () => {
        if (loggedIn) {navigate("/chat")}
        else { navigate("/login") }
    }

    return <div className="mainContainer">
        <div className={"titleContainer"}>
            { loggedIn ?
                <div>Welcome, {username}!</div> :
                <div>Welcome!</div> 
            }
        </div>
        <div>
            This is the home page.
        </div>
        <div className={"buttonContainer"}>
            <input
                className={"inputButton"}
                type="button"
                onClick={onButtonClick}
                value={loggedIn ? "Chat" : "Log in"} />
        </div>
    </div>
}

export default HomePage