import React, { useEffect, useState } from "react"
import { useParams } from "react-router-dom";

export default function Profile() {
    const [profile, setProfile] = useState(null);
    const [error, setError] = useState(null);
    const {username} = useParams(); 
     const token = localStorage.getItem("accessToken");
    useEffect(() => {
        const fetchProfileData = async (e) => {
            try {
                const response = await fetch(`http://localhost:8080/Profile/${username}`, {
                    method: 'GET',
                    headers: {
                        'Authorization':`Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                });

                if(response.ok) {
                    const data = await response.json();
                    setProfile(data);
                } else {
                    setError("Failed to fetch profile details");
                    console.error("Failed to fetch")
                }
            } catch (err) {
                setError("Can't Connect to Server")
                console.log(err);
            }
        };
        fetchProfileData();
    }, [username]);

    if (error) return <div style={{ padding: '20px', color: 'red' }}>{error}</div>;
    if (!profile) return <div style={{ padding: '20px' }}>Loading profile...</div>;

   return (
    <div className="container">
        <h1>Welcome, {profile.user_name}!</h1>
        <p className="logged-in">
            Logged in as: {profile.email}
        </p>

        <div className="stats-grid">
            <div className="stat-box">
                <h3>🔥 Current Streak</h3>
                <p className="number">{profile.current_streak} Days</p>
            </div>

            <div className="stat-box">
                <h3>🏆 Max Streak</h3>
                <p className="number">{profile.max_streak} Days</p>
            </div>

            <div className="stat-box">
                <h3>📊 Total Notes</h3>
                <p className="number">{profile.total_notes}</p>
            </div>

            <div className="stat-box">
                <h3>📝 Today's Notes</h3>
                <p className="number">{profile.today_notes}</p>
            </div>
        </div>


        <div className="latest-notes">
            <h2>Your Latest Notes</h2>

            {profile.notes && profile.notes.length > 0 ? (
                <ul className="notes-list">
                    {profile.notes.map((note) => (
                        <li key={note.note_id} className="note-item">
                            <p>{note.daily_note}</p>

                            <small className="note-date">
                                {note.created_at || 'Just now'}
                            </small>
                        </li>
                    ))}
                </ul>
            ) : (
                <p className="no-notes">
                    You haven't written any notes yet!
                </p>
            )}
        </div>

        <div className="All-notes">
             <button
                className="all-notes-btn"
                onClick={() => window.location.href = `/Profile/${username}/GetNotes`}
            >
             View All Notes
            </button>
</div>
    </div>
);
}``