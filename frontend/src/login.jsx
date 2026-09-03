import React, {useState} from "react";

export default function Login() {
    const [email, setEmail] = useState('');
    const[password, setPassword] = useState('')
    const[message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault();
        //setMessage('Loggin in...');

        try {
            const response = await fetch('http://localhost:8080/Dashboard/LoginUser', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({email: email, password: password})
            });

            const data = await response.json();

            if(response.ok) {
                setMessage('Login Successful');
            } else {
                setMessage('Invalid email or password');
            }
        } catch (error) {
            setMessage('Can not connect to the Server')
        }
       // alert(`Sending Email: ${email} Password: ${password}`);
    };

    return (
        <div className="login-container">
            <h1> Improve-Me </h1>
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="email"> Email: </label>
                    <input
                    type="email"
                    id="email"
                    placeholder="Enter your email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required/>
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
                <button type="submit">Login</button>
            </form>
            {message && <p id="message">{message}</p>}
        </div>
    )
}