document.getElementById("uploadBtn").addEventListener("click", () => {
    const fileInput = document.getElementById("fileInput");
    const selectBox = document.getElementById("modelSelect");
    const uploadBtn = document.getElementById("uploadBtn");
    const loader = document.getElementById("loader");

    const file = fileInput.files[0];
    if (!file) {
        alert("Please select a file first");
        return;
    }

    const selectedModel = selectBox.value;
    if (!selectedModel) {
        alert("Please select a model");
        return;
    }

    // Disable the button and show loader
    uploadBtn.disabled = true;
    loader.style.display = "block";

    // Prepare multipart/form-data
    const formData = new FormData();
    formData.append("file", file);

    fetch(`http://localhost:8099/api/app/upload/${selectedModel}`, {
        method: "POST",
        body: formData,
    })
        .then((response) => response.json())
        .then((data) => {
            if (data.error) {
                addMessage("system", data.error);
            } else {
                addMessage("user", `File uploaded successfully: ${data.filename}`);
            }
        })
        .catch((err) => {
            console.error(err);
            addMessage("system", "An error occurred while uploading the file.");
        })
        .finally(() => {
            // Re-enable the button and hide loader
            uploadBtn.disabled = false;
            loader.style.display = "none";
        });
});

// Helper function to add a new message to the chat window
function addMessage(sender, text) {
    const chatContent = document.getElementById("chat-content");
    const newMessage = document.createElement("div");

    newMessage.classList.add("chat-message");
    if (sender === "system") {
        newMessage.classList.add("system-message");
    } else {
        newMessage.classList.add("user-message");
    }

    const p = document.createElement("p");
    p.textContent = text;
    newMessage.appendChild(p);
    chatContent.appendChild(newMessage);

    // Auto-scroll to bottom
    chatContent.scrollTop = chatContent.scrollHeight;
}
