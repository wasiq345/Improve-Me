import React, {useState} from "react";
import { useNavigate, useParams } from "react-router-dom";

export default function CreateNote() {
    const[note, setNote] = useState("");
    const[message, setMessage] = useState("");
    const[error, setError] = useState("");
    const {username} = useParams(); 
    const navigate = useNavigate();
    const token = localStorage.getItem("accessToken");
    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch(`http://localhost:8080/Profile/${username}/CreateNote`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(
                    {daily_note: note}
                )
            });
            if(response.ok) {
                setMessage("Note Created Successfully")
                navigate(`/Profile/${username}`);
                setNote("");
            } else {
                 setError("Failed to fetch profile details");
            }
        } catch(err) {
            setError("Failed to Connect to Server")
        }
    };

    return (
        <div> className="create-note-container"
            <h1>Daily Note</h1>

            <form onSubmit={handleSubmit}>
                <textarea
                    value={note}
                    onChange ={(e) => setNote(e.target.value)}
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