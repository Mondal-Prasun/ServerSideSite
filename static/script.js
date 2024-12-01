// Counter logic
document.addEventListener("DOMContentLoaded", () => {
    const counterDisplay = document.getElementById("counter-display");
    const incrementButton = document.getElementById("increment-btn");
    const submitButton = document.getElementById("submit-btn");

    let counter = 0;

    incrementButton.addEventListener("click", () => {
        counter++;
        console.log(`counter: ${counter}`)
        counterDisplay.textContent = counter;
    });
    submitButton.addEventListener("click", () => {
        fetch('/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ counter }),
        }).then((res) => {
            console.log(`subbmited: ${counter}`)
        })
        // .then(response => response.json())
        // .then(data => {
        //     alert(`Server received counter value: ${data.counter}`);
        // })
        // .catch(error => {
        //     console.error('Error:', error);
        // });
    });
});
