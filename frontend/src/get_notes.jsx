import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function GetNotes() {
    const [notes, setNotes] = useState([]);
    const [error, setError] = useState('');
    const token = localStorage.getItem("accessToken");
    const {username} = useParams();
    const fetchGetNotes = async (e) => {
        try {
            const response = await fetch(`http://localhost:8080/Profile/${username}/GetNotes`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });
            if (response.ok) {
                const data = await response.json();
                setNotes(data);
            }
            else {
                setError("Failed To Fetch Notes")
            }
        } catch (err) {
            setError("Failed to Connect to Server")
        }
    };
    useEffect(() => { fetchGetNotes(); }, []);
    return (
        <div className="notes-container">
            <h1>Your Notes</h1>
            {error && (<p className="error-message"> {error} </p>)}
            {notes.length > 0 ? (
                <div className="notes-list">
                    {notes.map((note) => (
                        <div key={note.note_id} className="note-card" >
                            <p className="note-text"> {note.daily_note} </p>
                            <small className="note-date"> {note.created_at} </small>
                        </div>))}
                </div>) :
                (!error && (
                    <p className="no-notes">
                        You haven't written any notes yet!
                    </p>
                )
                )}
        </div>);
}