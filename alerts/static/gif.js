var defaultVoiceName = 'Microsoft Neerja Online (Natural) - English (India)';
const gifWS = new WebSocket("ws://localhost:3000/ws/123");
let gifQueue = [];
function handleGif() {
    if (gifQueue.length > 0) {
        const nextMessage = gifQueue.shift();
        console.log(nextMessage);
        const imgElement = document.getElementById("alert-user-icon");
        imgElement.src = nextMessage.url;
        console.log(imgElement.src);
        setTimeout(() => {
            imgElement.src = "";
            handleGif();
        }, 2000);
    }
    else {
        console.log("No messages to play");
    }
}
// Handle incoming messages
gifWS.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    if ('url' in msg)
        gifQueue.push(msg);
};
gifWS.onopen = function (event) {
    console.log("WebSocket connection opened:", event);
};
// Handle WebSocket errors
gifWS.onerror = function (event) {
    console.error("WebSocket error observed:", event);
};
