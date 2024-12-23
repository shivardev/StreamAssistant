interface Msg {
    userName:string
    platform:string
    action:string
}
const instaWS = new WebSocket("ws://localhost:3000/ws/123");
let instaQ: Msg[] = [];

function handleInstaAlert() {
    if (instaQ.length > 0) {
        const nextMessage = instaQ.shift();
        console.log(nextMessage)
        const wrapper = document.getElementById('alert-text-wrap')
        wrapper.style.opacity = '1'
        document.getElementById('userName').innerText = `${nextMessage.userName} has ${nextMessage.platform}`
        setTimeout(() => {
            wrapper.style.opacity = '0'
            handleInstaAlert()
        }, 2000);
    }
    else {
        console.log("No messages to play");
    }
}
// Handle incoming messages
instaWS.onmessage = function (event) {
    const msg: Msg = JSON.parse(event.data);
    if('platform' in msg)
        instaQ.push(msg);
};
instaWS.onopen = function (event) {
    console.log("WebSocket connection opened:", event);
};
// Handle WebSocket errors
instaWS.onerror = function (event) {

    console.error("WebSocket error observed:", event);
};



