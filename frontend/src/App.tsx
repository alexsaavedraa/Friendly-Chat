import React from "react"
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Chat from './pages/ChatPage.tsx';
import './App.css';
import LoginPage from './pages/LoginPage.tsx';
import HomePage from "./pages/HomePage.tsx";
import { useState } from 'react';

function App() {
  const [loggedIn, setLoggedIn] = useState(false)
  const [userName, setUsername] = useState("")

  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage username={userName} loggedIn={loggedIn} setLoggedIn={setLoggedIn}/>}/>
          <Route path="/chat" element={<Chat username={userName} />} />
          <Route path="/login" element={<LoginPage setLoggedIn={setLoggedIn} setUsername={setUsername} />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App

