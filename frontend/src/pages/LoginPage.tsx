import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { endpoint_base } from "../config";

const Login = (props: any) => {

    const { setLoggedIn, setGlobalUsername } = props
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

        checkAccountExists((accountExists: boolean) => {
            if (accountExists)
                logIn()
            else
                if (window.confirm("An account does not exist with this username: " + username + ". Do you want to create a new account?")) {
                    createAccount()
                    .then(() => {
                        // This code will only run after createAccount completes successfully
                        logIn();
                    })
                    .catch(error => {
                        // Handle any errors that occurred during createAccount
                        console.error("Error creating account:", error);
                    });
                }
        })        


    }

    const checkAccountExists =(callback: any)  => {
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
        .catch(e => setLoginError("Error finding account. Please try again."))
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
                setLoggedIn(true)
                setGlobalUsername(username)
                navigate("/chat")
            } else {
                window.alert("Wrong username or password")
            }
        })
        .catch(e => setLoginError("Error logging in. Please try again."))
    }

    const createAccount = async () => {
        try {
            const response = await fetch(`http://${endpoint_base}/signup`, {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });
            
            const data = await response.json();
            
            if ('success' === data.message) {
                localStorage.setItem("user", JSON.stringify({ username, token: data.token }));
                setLoggedIn(true);
                setGlobalUsername(username);
                navigate("/chat");
            } else {
                window.alert("Signup Failed. Please Try a different username");
            }
        } catch (error) {
            setLoginError("Error logging in. Please try a different username.");
        }
    };




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
                type="password"
                onChange={ev => setPassword(ev.target.value)}
                className={"inputBox"} />
                <br/>
            <label className="errorLabel">{loginError}</label>
        </div>
        <br />
        <div className={"inputContainer"}>
            <input
                className={"inputButton"}
                type="button"
                onClick={onButtonClick}
                value={"Log in"} />
            
        </div>
    </div>
}

export default Login