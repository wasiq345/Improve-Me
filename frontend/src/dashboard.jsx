import React from "react";
import { Link } from "react-router-dom";

export default function Dashboard() {

    return (
        <div>
            <h1> Improve-me </h1>
            <Link to="/LoginUser"><button>Login</button></Link>
            <Link to="/RegisterUser"><button>Register</button></Link>
        </div>
    )
}