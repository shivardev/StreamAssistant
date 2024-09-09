var defaultVoiceName = 'Microsoft Neerja Online (Natural) - English (India)';
var ws = new WebSocket("ws://localhost:3000/ws/123");
var messageQueue = [];
var synth = speechSynthesis, isSpeaking = false;
function playNextMessage() {
    if (messageQueue.length > 0) {
        var nextMessage = messageQueue.shift();
        var imgElement = document.getElementById("alert-user-icon");
        var userNameElement = document.getElementById("alert-user-name");
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
    var msg = JSON.parse(event.data);
    console.log(msg);
    messageQueue.push(msg);
    if (!isSpeaking) {
        playNextMessage();
    }
};
ws.onopen = function (event) {
    console.log("WebSocket connection opened:", event);
    var mainElement = document.getElementById("alert-text-wrap");
    mainElement.style.opacity = "0";
};
// Handle WebSocket errors
ws.onerror = function (event) {
    console.error("WebSocket error observed:", event);
};
function playAudio(msg) {
    var text = msg.authorName + " says, " + msg.messageContent;
    console.log(text, typeof text);
    var utterance = new SpeechSynthesisUtterance(msg.authorName + " says, " + msg.messageContent);
    var voices = speechSynthesis.getVoices();
    // Find and set the default voice based on name
    var defaultVoice = voices.filter(function (voice) { return voice.name === defaultVoiceName; });
    utterance.voice = defaultVoice[0];
    utterance.onstart = function () {
        var mainElement = document.getElementById("alert-text-wrap");
        mainElement.style.opacity = "1";
        isSpeaking = true; // Set the flag when speech starts
        console.log("Speech started");
    };
    // Event fired when the speech ends
    utterance.onend = function () {
        var mainElement = document.getElementById("alert-text-wrap");
        mainElement.style.opacity = "0";
        isSpeaking = false; // Clear the flag when speech ends
        console.log("Speech finished");
        playNextMessage();
    };
    synth.speak(utterance);
}
