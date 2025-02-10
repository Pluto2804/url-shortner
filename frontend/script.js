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