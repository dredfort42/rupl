import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
// import logo from './logo.svg';
import './App.css';
import Index from "./Index";
// import Login from "./Login";
// import CreateUser from "./CreateUser";

function App() {
  return (
    // <div className="App">
    //   <header className="App-header">
    //     <img src={logo} className="App-logo" alt="logo" />
    //     <p>
    //       Edit <code>src/App.js</code> and save to reload.
    //     </p>
    //     <a
    //       className="App-link"
    //       href="https://reactjs.org"
    //       target="_blank"
    //       rel="noopener noreferrer"
    //     >
    //       Learn React
    //     </a>
    //   </header>
    // </div>
    <Router>
      <Routes>
        <Route path="/" element={<Index />} />
        {/* <Route path="/login" element={<Login />} /> */}
        {/* <Route path="/create-user" element={<CreateUser />} /> */}
        {/* <Route path="/chat" element={<MainChat />} /> */}
        {/* <Route path="/chat/:channelId" element={<MainChat />} /> */}
      </Routes>
    </Router>
  );
}

export default App;
