import React from "react"
import { useNavigate } from "react-router-dom";

const HomePage = (props: any) => {
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
            Welcome to the home page of my chatroom app. To create an account, log in with any credentials and click OK -Alex Saavedra.
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