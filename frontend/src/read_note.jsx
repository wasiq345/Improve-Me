import React, { useEffect, useState } from "react"
import { useParams, useSearchParams } from "react-router-dom"

export default function ReadNote() {
    const [error, setError] = useState('');
    const token = localStorage.getItem("accessToken")
    const [note, setNote] = useState(null);
    const { username, noteId } = useParams();

    const fetchReadNote = async (e) => {
        try {
            const response = await fetch(`http://localhost:8080/Profile/${username}/ReadNote/${noteId}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });
            if (response.ok) {
                const data = await response.json();
                setNote(data);
            } else {
                setError("Failed to fetch the Note")
            }
        } catch (err) {
            setError("Failed to connect to Server")
        }
    };

    useEffect(() => { fetchReadNote(); }, []);
    if (!note) {
        return <div className="loading">Loading note...</div>;
    }
    return (
        <div className="read-note">
            <h1>Daily Note</h1>

            <div className="note-card">
                <p className="note-date">
                    {new Date(note.created_at).toLocaleString()}
                </p>

                <div className="note-content">
                    {note.daily_note}
                </div>
            </div>
        </div>
    );
}