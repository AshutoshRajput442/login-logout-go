import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./dashboard.css";


const Dashboard = () => {
  const [email, setEmail] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/dashboard");
    } else {
      // Decode the JWT (simplified for example purposes)
      const payload = JSON.parse(atob(token.split(".")[1]));
      setEmail(payload.email);
    }
  }, [navigate]);

  return (
    <div>
      <h2>Welcome, {email}</h2>
      <button onClick={() => {
        localStorage.removeItem("token");
        navigate("/");
      }}>
        Logout
      </button>
    </div>
  );
};

export default Dashboard;
