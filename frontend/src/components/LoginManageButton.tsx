import React from "react"
import LogoutIcon from '@mui/icons-material/Logout';
import LoginIcon from '@mui/icons-material/Login';
import { Button } from "@mui/material";
import { close } from "../api/index.ts";

const host = "192.168.0.180";
const port = 8080;
const endpoint_base = `${host}:${port}`;



interface LoginManageButtonProps {
    isLoggedIn: boolean
    setIsLoggedIn: Function
  }

const LoginManageButton: React.FC<LoginManageButtonProps> = (({isLoggedIn, setIsLoggedIn}) => {

    const handleLoginLogout = async () => {
        if (isLoggedIn) {
            setIsLoggedIn(false);
            
            try {
                const userDataString = localStorage.getItem("user");
                const userData = userDataString ? JSON.parse(userDataString) : null;
                const { username, token } = userData || {};
                const response = await fetch(`http://${endpoint_base}/logout?username=${username}&token=${token}`, {
                    method: "POST",
                    headers: {
                        'Content-Type': 'application/json'
                    },
                });
                // Handle response if needed
            } catch (error) {
                console.error('Error during logout:', error);
            }

            close()
            localStorage.removeItem("user")
        }
    };

    return (
        <Button
            sx={{ borderRadius: 10 , 
                  color:"cornflowerblue", 
                  bgcolor:"white",
                  ":hover": {
                    color:"white"
                  }
                }}
            variant="contained"
            startIcon={isLoggedIn ? <LogoutIcon/> : <LoginIcon/>}
            onClick={handleLoginLogout}
        >
            {isLoggedIn ? "Log Out" : "Log In"}
        </Button>
    )
});

export default LoginManageButton;