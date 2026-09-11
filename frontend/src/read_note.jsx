import React, { useEffect, useState } from "react"
import { useParams, useSearchParams } from "react-router-dom"

export default function ReadNote() {
    const [error, setError] = useState('');
    const token = localStorage.getItem("accessToken")
    const [note, setNote] = useState('');
    const { username, noteId } = useParams();
    const [isEditing, setIsEditing] = useState(false);
    const [editText, setEditText] = useState('');

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
                setEditText(data);
            } else {
                setError("Failed to fetch the Note")
            }
        } catch (err) {
            setError("Failed to connect to Server")
        }
    };

    const fetchUpdateNote = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`http://localhost:8080/Profile/${username}/UpdateNote/${noteId}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ daily_note: editText })
            });
            if (response.ok) {
                const data = await response.json();
                setNote({ ...note, daily_note: editText });
                setEditText(data);
                setIsEditing(false);
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

                {isEditing ? (
                    <form onSubmit={fetchUpdateNote} className="edit-note-form">
                        <textarea
                            value={editText}
                            onChange={(e) => setEditText(e.target.value)}
                            rows="6"
                            className="edit-note-textarea"
                            required
                        />
                        <div className="edit-form-actions">
                            <button type="submit" className="btn btn-save">Save Changes</button>
                            <button type="button" className="btn btn-cancel" onClick={() => setIsEditing(false)}>Cancel</button>
                        </div>
                    </form>
                ) : (
                    <div>
                        <div className="note-content">
                            {note.daily_note}
                        </div>
                        <button className="btn btn-edit" onClick={()=>{setEditText(note.daily_note); setIsEditing(true);}}>Edit Note</button>
                    </div>
                )}
            </div>
        </div>
    );
}
