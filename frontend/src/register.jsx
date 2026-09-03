import React, {useState} from "react";

export default function Register() {
    const[email, setEmail] = useState('');
    const[password, setPassword] = useState('');
    const[message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('Registering ...')


        try {
            const response = await fetch('http://localhost:8080/Dashboard/RegisterUser', {
                method: 'POST',
                headers: {
                    'Content-type': 'application/json',
                },
                body: JSON.stringify({email: email, password: password})
            });

            const data = await response.json();
            if(response.ok) {
                setMessage('Registration Successful');
            } else {
                setMessage('Invalid email or password');
            }
        } catch (error) {
            setMessage('Can not connect to server');
        }
    };

return (
    <div className="register-container">
        <h1> Register Account </h1>
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="email">Email: </label>
                <input
                type="email"
                id="email"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required></input>
            </div>
            <br />
                <div>
                    <label htmlFor="password"> Password: </label>
                    <input 
                    type="password"
                    id="password"
                    placeholder="Enter your password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required/>
                </div>
            <br />
            <button type="submit">Register</button>
        </form>
        {message && <p id="message">{message}</p>}
    </div> 
)
};