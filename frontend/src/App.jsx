
import './App.css'

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
// import LoginForm from "./LoginForm";
import Dashboard from "./Dashboard";
import LoginSignup from "./LoginSignup";
// import Hello from "./Hello";


function App() {
  return (
    <Router>
      <Routes>
        {/* <Route path="/" element={<LoginForm />} /> */}
        <Route path="/" element={<LoginSignup />} />
        <Route path="/dashboard" element={<Dashboard />} />
        {/* <Route path="/hello" element={<Hello />} /> */}

      </Routes>
    </Router>
  );
}

export default App;
