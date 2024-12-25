import { chromium } from "playwright";
import { LikeResponseObj } from "./utils/likesApi.modal";
import { CommentObj } from "./utils/commentApi.modal";
import { API_URLS } from "./utils/constants";
import { msgsPayload } from "./utils/interfaces";

(async () => {
  const sendPostRequest = (url: string, payload: any) => {
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        console.log("Request sent successfully");
      })
      .catch((error) => {
        console.error("Error sending request:", error);
      });

  };

  let globalMaxLikeCount = 0;
  let previousLikeCount = 0;
  const urls: string[] = [
    "https://www.youtube.com/youtubei/v1/live_chat/get_live_chat?prettyPrint=false",
    "https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false",
  ];
  
  // Launch the browser
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();
  function findActionByKey(arr: any, key: string): any | null {
    return arr.find((action: any) => action[key] !== undefined) || null;
  }

  let congoMap = new Map<number, boolean>();
  function checkAndCongratulate(number: number): boolean {
    let base = Math.floor(number / 5) * 5;

    if (number >= base && number < base + 5) {
      if (!congoMap.has(base)) {
        congoMap.set(base, true);
        return true;
      }
    }
    return false;
  }
  // Intercept network responses

  page.on("response", async (response) => {
    const url = response.url();
    console.log(url)
    if (urls.includes(url)) {
      if (urls.indexOf(url) == 1) {
        try {
          const json: LikeResponseObj = await response.json();
          const actions = json.actions;
          // Viewer count
          const actionWithViewership = findActionByKey(
            json.actions,
            "updateViewershipAction"
          );
          const viewerCount =
            actionWithViewership.updateViewershipAction.viewCount
              .videoViewCountRenderer.originalViewCount;
          // Like count
          const likeCountNow = Number(
            json.frameworkUpdates.entityBatchUpdate.mutations[0].payload
              .likeCountEntity.likeCountIfDislikedNumber
          );

          const statsPayload = {
            likes: Number(likeCountNow),
            previousLikes: previousLikeCount,
            viewers: Number(viewerCount),
            maxLikes: globalMaxLikeCount,
            shouldCongratulate: checkAndCongratulate(globalMaxLikeCount),
          };
          console.log(statsPayload);
          sendPostRequest(API_URLS.stats, { stats: statsPayload });
          previousLikeCount = likeCountNow;
          globalMaxLikeCount = Math.max(globalMaxLikeCount, likeCountNow);
        } catch (error) {
          console.log("Response body could not be parsed as JSON.", error);
        }
      } else {
        console.log("mostly chat URL")
        try {
          const json: CommentObj = await response.json();
          const actions =
            json.continuationContents.liveChatContinuation.actions;
          const msgsPayload: msgsPayload[] = [];
          if (actions) {
            for (const action of actions) {
              const item = action.addChatItemAction.item;
              let finalMessage = "";
              let authorName = "";
              let authorId = "";
              let timestamp = "";
              let authorPhoto = "";
              authorName =
                item.liveChatTextMessageRenderer.authorName.simpleText;
              authorId =
                item.liveChatTextMessageRenderer.authorExternalChannelId;
              timestamp = item.liveChatTextMessageRenderer.timestampUsec;
              authorPhoto =
                item.liveChatTextMessageRenderer.authorPhoto.thumbnails[
                  item.liveChatTextMessageRenderer.authorPhoto.thumbnails
                    .length - 1
                ].url;
              for (const run of item.liveChatTextMessageRenderer.message.runs) {
                if (run.text) {
                  finalMessage += run.text;
                }
                // if(run.emoji){
                //   finalMessage += run.emoji.emojiId;
                // }
              }
              const payload: msgsPayload = {
                authorName: authorName,
                authorId: authorId,
                timestamp: timestamp,
                authorPhotoUrl: authorPhoto,
                messageContent: finalMessage,
              };
              console.log(payload.authorName);
              msgsPayload.push(payload);
            }
            sendPostRequest(API_URLS.msgs, { messages: msgsPayload });
          }else{
            console.log('noactions')
          }
        } catch (error) {
          console.log("Response body could not be parsed as JSON.");
        }
      }
    }
  });

  await page.goto(
    "https://www.youtube.com/@blazingbane5565/live"
  );

  console.log("Listening for API calls...");
  while (true) {
    await page.waitForTimeout(10000);
  }
})();
