import React from "react";
import LoginManageButton from "./LoginManageButton.tsx";


const Header = (props) => {
    const { loggedIn, setLoggedIn } = props
    return <div className="header">
            <h2>Nimble Challenge Chat</h2>
            <LoginManageButton 
                isLoggedIn={loggedIn}
                setIsLoggedIn={setLoggedIn}
            />
        </div>
};

export default Header;