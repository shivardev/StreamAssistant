"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const playwright_1 = require("playwright");
const constants_1 = require("./utils/constants");
(() => __awaiter(void 0, void 0, void 0, function* () {
    const sendPostRequest = (url, payload) => {
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
    const urls = [
        "https://www.youtube.com/youtubei/v1/live_chat/get_live_chat?prettyPrint=false",
        "https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false",
    ];
    // Launch the browser
    const browser = yield playwright_1.chromium.connectOverCDP("http://127.0.0.1:8989");
    const defaultContext = browser.contexts()[0];
    const page = defaultContext.pages()[0];
    function findActionByKey(arr, key) {
        return arr.find((action) => action[key] !== undefined) || null;
    }
    let congoMap = new Map();
    function checkAndCongratulate(number) {
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
    page.on("response", (response) => __awaiter(void 0, void 0, void 0, function* () {
        const url = response.url();
        console.log(url);
        if (urls.includes(url)) {
            if (urls.indexOf(url) == 1) {
                try {
                    const json = yield response.json();
                    const actions = json.actions;
                    // Viewer count
                    const actionWithViewership = findActionByKey(json.actions, "updateViewershipAction");
                    const viewerCount = actionWithViewership.updateViewershipAction.viewCount
                        .videoViewCountRenderer.originalViewCount;
                    // Like count
                    const likeCountNow = Number(json.frameworkUpdates.entityBatchUpdate.mutations[0].payload
                        .likeCountEntity.likeCountIfDislikedNumber);
                    const statsPayload = {
                        likes: Number(likeCountNow),
                        previousLikes: previousLikeCount,
                        viewers: Number(viewerCount),
                        maxLikes: globalMaxLikeCount,
                        shouldCongratulate: checkAndCongratulate(globalMaxLikeCount),
                    };
                    console.log(statsPayload);
                    sendPostRequest(constants_1.API_URLS.stats, { stats: statsPayload });
                    previousLikeCount = likeCountNow;
                    globalMaxLikeCount = Math.max(globalMaxLikeCount, likeCountNow);
                }
                catch (error) {
                    console.log("Response body could not be parsed as JSON.", error);
                }
            }
            else {
                console.log("mostly chat URL");
                try {
                    const json = yield response.json();
                    const actions = json.continuationContents.liveChatContinuation.actions;
                    const msgsPayload = [];
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
                                item.liveChatTextMessageRenderer.authorPhoto.thumbnails[item.liveChatTextMessageRenderer.authorPhoto.thumbnails
                                    .length - 1].url;
                            for (const run of item.liveChatTextMessageRenderer.message.runs) {
                                if (run.text) {
                                    finalMessage += run.text;
                                }
                                // if(run.emoji){
                                //   finalMessage += run.emoji.emojiId;
                                // }
                            }
                            const payload = {
                                authorName: authorName,
                                authorId: authorId,
                                timestamp: timestamp,
                                authorPhotoUrl: authorPhoto,
                                messageContent: finalMessage,
                            };
                            console.log(payload.authorName);
                            msgsPayload.push(payload);
                        }
                        sendPostRequest(constants_1.API_URLS.msgs, { messages: msgsPayload });
                    }
                    else {
                        console.log("noactions");
                    }
                }
                catch (error) {
                    console.log("Response body could not be parsed as JSON.");
                }
            }
        }
    }));
    yield page.goto("https://www.youtube.com/@blazingbane5565/live");
    console.log("Listening for API calls...");
    while (true) {
        yield page.waitForTimeout(10000);
    }
}))();
