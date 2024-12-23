var defaultVoiceName = 'Microsoft Neerja Online (Natural) - English (India)';
// const ws = new WebSocket("ws://localhost:3000/ws/123");
const ws = new WebSocket("ws://10.0.0.236:3000/ws/123");
let messageQueue = [];
let synth = speechSynthesis, isSpeaking = false;
function playNextMessage() {
    if (messageQueue.length > 0) {
        const nextMessage = messageQueue.shift();
        console.log(nextMessage);
        const imgElement = document.getElementById("alert-user-icon");
        const userNameElement = document.getElementById("alert-user-name");
        imgElement.src = nextMessage.authorPhotoUrl;
        userNameElement.innerHTML = nextMessage.authorName;
        if (nextMessage) {
            playAudio(nextMessage); // Play the next message
        }
    }
    else {
        console.log("No messages to play");
    }
}
// Handle incoming messages
ws.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    console.log(msg);
    if ('authorName' in msg)
        messageQueue.push(msg);
    if (!isSpeaking) {
        playNextMessage();
    }
};
ws.onopen = function (event) {
    console.log("WebSocket connection opened:", event);
    const mainElement = document.getElementById("alert-text-wrap");
    mainElement.style.opacity = "0";
};
// Handle WebSocket errors
ws.onerror = function (event) {
    console.error("WebSocket error observed:", event);
};
function playAudio(msg) {
    const text = `${msg.authorName} says, ${msg.messageContent}`;
    console.log(text, typeof text);
    let utterance = new SpeechSynthesisUtterance(`${msg.authorName} says, ${msg.messageContent.replace("!speak", "")}`);
    const voices = speechSynthesis.getVoices();
    // Find and set the default voice based on name
    const defaultVoice = voices.filter(voice => voice.name === defaultVoiceName);
    utterance.voice = defaultVoice[0];
    utterance.onstart = function () {
        const mainElement = document.getElementById("alert-text-wrap");
        mainElement.style.opacity = "1";
        isSpeaking = true; // Set the flag when speech starts
        console.log("Speech started");
    };
    // Event fired when the speech ends
    utterance.onend = function () {
        const mainElement = document.getElementById("alert-text-wrap");
        mainElement.style.opacity = "0";
        isSpeaking = false; // Clear the flag when speech ends
        console.log("Speech finished");
        playNextMessage();
    };
    synth.speak(utterance);
}
