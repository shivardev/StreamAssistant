let serverIP = "10.0.0.236";

const args = process.argv.slice(2);
args.forEach((arg, index) => {
  if (arg === "-env" && args[index + 1]) {
    const env = args[index + 1].toLowerCase();
    if (env === "dev") {
      serverIP = "10.0.0.128"; 
    } else {
      serverIP = "10.0.0.236";
    } 
  }
});
console.log(serverIP)
export const API_URLS = {
  msgs: `http://${serverIP}:3000/takemsgs`,
  stats: `http://${serverIP}:3000/takestats`,
  videoLink: `http://${serverIP}:3000/streamurl`
};

export const urls: string[] = [
  "https://www.youtube.com/youtubei/v1/live_chat/get_live_chat?prettyPrint=false",
  "https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false",
];