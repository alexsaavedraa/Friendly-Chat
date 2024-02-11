import React from "react"
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import ChatPage from './pages/ChatPage.tsx';
import './App.css';
import Header from './components/Header.tsx'
import LoginPage from './pages/LoginPage.tsx';
import HomePage from "./pages/HomePage.tsx";
import { useState } from 'react';
import TimeAgo from 'javascript-time-ago'

import en from 'javascript-time-ago/locale/en.json'

TimeAgo.addDefaultLocale(en)

function App() {
  const [loggedIn, setLoggedIn] = useState(false)
  const [userName, setUsername] = useState("")
  return (
    <div className="App">
      <Header loggedIn={loggedIn} setLoggedIn={setLoggedIn} />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage username={userName} loggedIn={loggedIn} setLoggedIn={setLoggedIn}/>}/>
          <Route path="/chat" element={<ChatPage loggedIn={loggedIn} setLoggedIn={setLoggedIn} />} />
          <Route path="/login" element={<LoginPage setLoggedIn={setLoggedIn} setUsername={setUsername} />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App

