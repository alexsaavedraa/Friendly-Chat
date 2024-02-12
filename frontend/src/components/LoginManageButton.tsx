import React from "react"
import LogoutIcon from '@mui/icons-material/Logout';
import LoginIcon from '@mui/icons-material/Login';
import { Button } from "@mui/material";
import { close } from "../api/index.ts";


interface LoginManageButtonProps {
    isLoggedIn: boolean
    setIsLoggedIn: Function
  }

const LoginManageButton: React.FC<LoginManageButtonProps> = (({isLoggedIn, setIsLoggedIn}) => {

    const handleLoginLogout = (() => {
        if (isLoggedIn) {
            setIsLoggedIn(false);
            close()
            localStorage.removeItem("user")
        }
        
    });

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