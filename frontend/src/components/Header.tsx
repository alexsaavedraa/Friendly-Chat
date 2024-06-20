import React from "react";
import LogoutButton from "./LogoutButton.tsx";


const Header = (props) => {
    const { loggedIn, setLoggedIn } = props
    return <div className="header">
            <h2>Friendly Chat</h2>
            <LogoutButton 
                isLoggedIn={loggedIn}
                setIsLoggedIn={setLoggedIn}
            />
        </div>
};

export default Header;