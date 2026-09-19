import React, { useEffect, useState } from "react";
import { Link, Navigate, useNavigate, useParams } from "react-router-dom";
import { apiFetch } from "./api_fetch";

export default function GetNotes() {
    const [notes, setNotes] = useState([]);
    const [error, setError] = useState('');
    const { username } = useParams();
    const navigate = useNavigate();
    const fetchGetNotes = async (e) => {
        try {
            // const response = await fetch(`http://localhost:8080/Profile/${username}/GetNotes`, {
            //     method: 'GET',
            //     headers: {
            //         'Authorization': `Bearer ${token}`,
            //         'Content-Type': 'application/json'
            //     }
            // });
            // if (response.ok) {
            //     const data = await response.json();
            //     setNotes(data);
            // }
            // else {
            //     setError("Failed To Fetch Notes")
            // }
            const response = await apiFetch(`http://localhost:8080/Profile/${username}/GetNotes`, {
                method: 'GET',                
            });
            if (response.ok) {
                const data = await response.json();
                setNotes(data);
            } else if(response.status == 401) navigate("/LoginUser", {replace:true}); 
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
                        <Link to={`/Profile/${username}/ReadNote/${note.note_id}`} key={note.note_id} className="note-card-link">
                            <div key={note.note_id} className="note-card" >
                                <p className="note-text"> {note.daily_note} </p>
                                <small className="note-date"> {new Date(note.created_at).toLocaleString()} </small>
                            </div>
                        </Link>
                    ))}    
                </div>):
                (!error && (
                    <p className="no-notes">
                        You haven't written any notes yet!
                    </p>
                )
                )}
        </div>);
}