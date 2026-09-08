import React, { useEffect, useState } from "react"

export default function Profile() {
    const [profile, setProfile] = useState(null);
    const [error, setError] = useState(null);
    const username = localStorage.getItem("username");
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
        <div style={containerStyle}>
            <h1>Welcome, {profile.user_name}!</h1>
            <p style={{ color: '#666' }}>Logged in as: {profile.email}</p>

            <div style={statsGridStyle}>
                <div style={statBoxStyle}>
                    <h3>🔥 Current Streak</h3>
                    <p style={numberStyle}>{profile.current_streak} Days</p>
                </div>
                <div style={statBoxStyle}>
                    <h3>🏆 Max Streak</h3>
                    <p style={numberStyle}>{profile.max_streak} Days</p>
                </div>
                <div style={statBoxStyle}>
                    <h3>📊 Total Notes</h3>
                    <p style={numberStyle}>{profile.total_notes}</p>
                </div>
                <div style={statBoxStyle}>
                    <h3>📝 Today's Notes</h3>
                    <p style={numberStyle}>{profile.today_notes}</p>
                </div>
            </div>

            <div style={{ marginTop: '30px' }}>
                <h2>Your Latest Notes</h2>
                {profile.Notes && profile.Notes.length > 0 ? (
                    <ul style={{ listStyleType: 'none', padding: 0 }}>
                        {profile.Notes.map((note, index) => (
                            <li key={note.note_id} style={noteItemStyle}>
                                <p>{note.daily_note}</p> 
                                <small style={{ color: '#999' }}>{note.created_at || 'Just now'}</small>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p style={{ color: '#888' }}>You haven't written any notes yet!</p>
                )}
            </div>
        </div>
    );
}

const containerStyle = { maxWidth: '800px', margin: '40px auto', padding: '20px', fontFamily: 'Arial, sans-serif' };
const statsGridStyle = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: '20px', marginTop: '20px' };
const statBoxStyle = { background: '#f4f4f9', padding: '15px', borderRadius: '8px', textAlign: 'center', boxShadow: '0 2px 4px rgba(0,0,0,0.1)' };
const numberStyle = { fontSize: '24px', fontWeight: 'bold', margin: '10px 0 0 0', color: '#007BFF' };
const noteItemStyle = { background: '#fff', border: '1px solid #ddd', padding: '15px', borderRadius: '6px', marginBottom: '10px', boxShadow: '0 1px 3px rgba(0,0,0,0.05)' };