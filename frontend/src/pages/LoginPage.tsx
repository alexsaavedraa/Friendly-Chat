import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const Login = (props) => {
    let host = window.location.hostname;
    const port = 8080;
    const endpoint_base = `${host}:${port}`

    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [loginError, setLoginError] = useState("")
    
    const navigate = useNavigate();
        
    const onButtonClick = () => {
        // Check if the user has entered both fields correctly
        if ("" === username) {
            setLoginError("Please enter your username")
            return
        }

        if ("" === password) {
            setLoginError("Please enter a password")
            return
        }

        checkAccountExists(accountExists => {
            if (accountExists)
                logIn()
            else
                if (window.confirm("An account does not exist with this email address: " + username + ". Do you want to create a new account?")) {
                    logIn()
                }
        })        


    }

    const checkAccountExists = (callback) => {
        fetch(`http://${endpoint_base}/check-account`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
              },
            body: JSON.stringify({username})
        })
        .then(r => r.json())
        .then(r => {
            callback(r?.userExists)
        })
    }

    // Log in a user using username and password
    const logIn = () => {
        fetch(`http://${endpoint_base}/auth`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
              },
            body: JSON.stringify({username, password})
        })
        .then(r => r.json())
        .then(r => {
            if ('success' === r.message) {
                localStorage.setItem("user", JSON.stringify({username, token: r.token}))
                props.setLoggedIn(true)
                props.setUsername(username)
                navigate("/chat")
            } else {
                window.alert("Wrong username or password")
            }
        })
    }


    return <div className={"mainContainer"}>
        <div className={"titleContainer"}>
            <div>Login</div>
        </div>
        <br />
        <div className={"inputContainer"}>
            <input
                value={username}
                placeholder="Enter your username here"
                onChange={ev => setUsername(ev.target.value)}
                className={"inputBox"} />
        </div>
        <br />
        <div className={"inputContainer"}>
            <input
                value={password}
                placeholder="Enter your password here"
                onChange={ev => setPassword(ev.target.value)}
                className={"inputBox"} />
        </div>
        <br />
        <div className={"inputContainer"}>
            <input
                className={"inputButton"}
                type="button"
                onClick={onButtonClick}
                value={"Log in"} />
            <br/>
            <label className="errorLabel">{loginError}</label>
        </div>
    </div>
}

export default Login