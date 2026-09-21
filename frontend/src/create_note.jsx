import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { apiFetch } from "./api_fetch";

export default function CreateNote() {
    const [note, setNote] = useState("");
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const { username } = useParams();
    const navigate = useNavigate();
    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
        
            const response = await apiFetch(`http://localhost:8080/Profile/${username}/CreateNote`, {
                method: 'POST',
                body: JSON.stringify(
                    { daily_note: note }
                )
            });
            if (response.ok) {
                setMessage("Note Created Successfully")
                navigate(`/Profile/${username}`);
                setNote("");
            } else if (response.status == 401) navigate("/LoginUser", { replace: true });
            else {
                setError("Failed To Fetch")
            }
        } catch (err) {
            setError("Failed to Connect to Server")
        }
    };

    return (
        <div className="create-note-container">
            <h1>Daily Note</h1>

            <form onSubmit={handleSubmit}>
                <textarea
                    value={note}
                    onChange={(e) => setNote(e.target.value)}
                    placeholder="Write your note ..."
                    rows="8"
                />
                <button type="submit">Create Note</button>
            </form>

            {message && (<p className="success-message"> {message}</p>)}
            {error && (<p className="error-message">{error}</p>)}
        </div>
    );
}