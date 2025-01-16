import { chromium, firefox } from "playwright";
import { LikeResponseObj } from "./utils/likesApi.modal";
import { CommentObj } from "./utils/commentApi.modal";
import { API_URLS, urls } from "./utils/constants";
import { msgsPayload, urlPayload } from "./utils/interfaces";
import { sendPostRequest } from "./utils/utils";

(async () => {
  let globalMaxLikeCount = 0;
  let previousLikeCount = 0;
  let videoID: undefined | string = undefined
  // Launch the browser
  let isVideoURLSent: boolean = false
  const browser = await firefox.launch({ headless: true });
  // const defaultContext = browser.contexts()[0];
  // const withVideoPage = defaultContext.pages()[0];
  const withVideoPage = await browser.newPage();
  // const chatOnlyPage = await defaultContext.newPage();

  let congoMap = new Map<number, boolean>();

  withVideoPage.on("response", async (response) => {
    const url = response.url();
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
          videoID = (json.continuationContents.liveChatContinuation.continuations[0].invalidationContinuationData.invalidationId.topic as string).split('~')[1]

          const actions =
            json.continuationContents.liveChatContinuation.actions;
          const msgsPayload: msgsPayload[] = [];
          if (actions) {
            for (const action of actions) {
              try {
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
                  videoId:videoID
                };
                console.log(payload.authorName);
                msgsPayload.push(payload);
              }
              catch (error) {
                console.log(error)
              }
            }
            sendPostRequest(API_URLS.msgs, { messages: msgsPayload });
          } else {
            console.log("noactions");
          }
          if (!isVideoURLSent) {
            const linkPlayload: urlPayload = {
              url: videoID
            };
            sendPostRequest(API_URLS.videoLink, linkPlayload)
            isVideoURLSent = true
          }
        } catch (error) {
          console.log("Response body could not be parsed as JSON.", error);
        }
      }
    }
  });

  await withVideoPage.goto(
    "https://www.youtube.com/@blazingbane5565/live"
  );
  // setTimeout(async () => {
  //   await chatOnlyPage.goto(
  //     "https://www.youtube.com/live_chat?v=" + videoID
  //   );
  // }, 20000);


  console.log("Listening for API calls...");
  while (true) {
    await withVideoPage.waitForTimeout(10000);
  }
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
  function findActionByKey(arr: any, key: string): any | null {
    return arr.find((action: any) => action[key] !== undefined) || null;
  }
})();