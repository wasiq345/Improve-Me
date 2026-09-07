import React from "react"

export default function Profile() {
    const fetchProfileData = async (e) => {
        const username = "xyz"
        const token = localStorage.getItem("accessToken")
        e.preventDefault();
            const response = await fetch(`http://localhost:8080/Profile?username={username}`, {
                method: 'GET',
                headers: {
                     'Authorization': 'Bearer ${token}',
                     'Content-Type': 'application/json'
                }
            });
            if(response.ok) {
                const data = await response.json();
                // 
            } else console.error("failed to fetch")
        }

    return (
        <div>
        <h1>Welcome to Your Profile</h1>
        </div>
    )
};