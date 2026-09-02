const LoginForm = document.getElementById("LoginForm");
const message = document.getElementById("message");

LoginForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    try {
        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        const response = await fetch("http://localhost:8080/Dashboard/LoginUser", {
            method: "POST",
            headers: {
                "Content-Type":"application/json"
            },

            body: JSON.stringify({
                email: email,
                password: password
            })
        });

        const data = await response.json()

        if(!response.ok()) {
            message.textContent = data.error();
            return;
        }
    
        message.textContent = data;
    } catch (error) {
        message.textContent = "Unable to connect to server";
    }   
})