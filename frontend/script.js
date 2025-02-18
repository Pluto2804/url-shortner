document.addEventListener("DOMContentLoaded",function(){
    const form = document.querySelector(".url-form");
    const input = document.querySelector(".url-input");

    form.addEventListener("submit",async function(event){
        event.preventDefault();
        const longURL = input.value.trim();
        if(!longURL){
            alert("Please enter a valid URL.")
            return;
        }
        try{
            const response = await fetch("http://localhost:8080/api/shorten",{
                method:"POST",
                headers:{
                    "Content-Type":"application/json",
                },
                body:JSON.stringify({url:longURL}),
            });
            if(!response.ok){
                throw new Error("Failed to shorten the URL");
            }
            const data = await response.json();
            const shortUrl = data.shortUrl;


            //Display the shortened URL
            displayShortenedUrl(shortUrl);
        }catch(error){
            console.error("Error:",error);
            alert("Something went wrong.Please try again. ");
        }
        
        
    });
    function displayShortenedUrl(shortUrl){
        {
            const resultDiv = document.createElement("div");
            resultDiv.innerHTML=`
            <p>Shortened URL: <a href="${shortUrl}" target="_blank">${shortUrl}</a></p>`;
            form.appendChild(resultDiv);
        }
    }
})
function openLoginModal(){
    document.getElementById("loginModal").style.display="block";
}
function openSignupModal(){
    document.getElementById("signupModal").style.display="block";
}
function closeLoginModal(){
    document.getElementById("loginModal").style.display="none";
}
function closeSignupModal(){
    document.getElementById("signupModal").style.display="none";
}
async function login(){
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password })
    });
    const data = await response.json();

    if(response.ok){
        alert("Login Successful!");
        localStorage.setItem("token",data.token);
        closeLoginModal();

    }else{
        alert("Login Failed: "+ data.error);
    }
}
async function signup() {
    const username = document.getElementById("signup-username").value;
    const email = document.getElementById("signup-email").value;
    const password = document.getElementById("signup-password").value;

    try {
        const response = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, email, password }),
        });

        const data = await response.json();

        if (response.status === 409) {
            alert("Email already registered! Try logging in.");
        } else if (response.ok) {
            alert("Signed Up Successfully!");
            localStorage.setItem("token", data.token);
            closeSignupModal();
        } else {
            alert("Signup Failed: " + data.error);
        }
    } catch (error) {
        console.error("Error:", error);
        alert("Something went wrong. Please try again.");
    }
}
document.querySelector(".auth-buttons .submit-btn").addEventListener("click", openLoginModal);
document.getElementById("s-2").addEventListener("click",openSignupModal);