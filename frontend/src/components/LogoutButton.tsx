import React from "react"
import LogoutIcon from '@mui/icons-material/Logout';
import { Button } from "@mui/material";
import { close } from "../api/index";
import { endpoint_base, protocol_base } from "../config";




interface LogoutButtonProps {
    isLoggedIn: boolean
    setIsLoggedIn: Function
  }

const LogoutButton: React.FC<LogoutButtonProps> = (({isLoggedIn, setIsLoggedIn}) => {

    const handleLoginLogout = async () => {
        if (isLoggedIn) {
            setIsLoggedIn(false);
            
            try {
                const userDataString = localStorage.getItem("user");
                const userData = userDataString ? JSON.parse(userDataString) : null;
                const { username, token } = userData || {};
                const response = await fetch(`${protocol_base}//${endpoint_base}/logout?username=${username}&token=${token}`, {
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
    if (isLoggedIn) {
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
                startIcon={<LogoutIcon/> }
                onClick={handleLoginLogout}
            >
                {isLoggedIn ? "Log Out" : "Log In"}
            </Button>
        )
    } else {
        return null 
    }
});

export default LogoutButton;