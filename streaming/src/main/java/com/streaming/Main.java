package com.streaming;

import java.net.URL;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;

import com.helper.LiveChatMessage;
import com.helper.MessageProcessor;

public class Main {

    static String lastChatId = "";
    static String chatUrl = "";
    static String[] ignoredAuthers = new String[]{"Nightbot", "YouTube"};
    // static final String ChannelLiveUrl = "https://www.youtube.com/channel/UCPFM_Ug62Ei3CUfvquG4KOg/live";
    static final String ChannelLiveUrl = "https://www.youtube.com/watch?v=FVgyYg368EU";

    static URL startStandaloneGrid() {
        int port = 8989;
        try {
            Main.main(
                    new String[]{
                        "standalone",
                        "--port",
                        String.valueOf(port),
                        "--selenium-manager",
                        "true",
                        "--enable-managed-downloads",
                        "true",
                        "--log-level",
                        "WARNING"
                    });
            return new URL("http://localhost:" + port);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public static void main(String[] args) {
        try {
            // Fetch the live broadcase from youtube API and get the live video link
            WebDriver driver;
            // Create ChromeOptions
            ChromeOptions options = new ChromeOptions();
            // Connect to the existing Chrome instance
            driver = new RemoteWebDriver(new URL("http://127.0.0.1:8989/wd/hub"), options);
            driver.get(ChannelLiveUrl);

            // Wait for the page to load
            try {
                Thread.sleep(4000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            List<WebElement> mainIframe = driver.findElements(By.className("ytd-live-chat-frame"));
            if (!mainIframe.isEmpty()) {
                chatUrl = mainIframe.get(0).getAttribute("src");
            } else {
                throw new RuntimeException("No chat found");
            }
            driver.get(chatUrl);
            try {
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            while (true) {
                fetchNewMessages(driver);

                // Adjust sleep time as needed to control polling frequency
                try {
                    Thread.sleep(1000); // Wait before fetching new messages again
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void fetchNewMessages(WebDriver driver) {
        try {
            List<WebElement> chatMessageElements = driver.findElements(By.cssSelector("yt-live-chat-text-message-renderer"));
            List<LiveChatMessage> newMessages = new ArrayList<>();

            boolean foundLastChatId = false;

            for (WebElement chatMessageElement : chatMessageElements) {
                String messageId = chatMessageElement.getAttribute("id");

                if (lastChatId.isEmpty()) {
                    LiveChatMessage chatContent = new LiveChatMessage(chatMessageElement);
                    newMessages.add(chatContent);
                    lastChatId = messageId;
                } else {
                    if (messageId.equals(lastChatId)) {
                        foundLastChatId = true;
                    } else if (foundLastChatId) {
                        LiveChatMessage chatContent = new LiveChatMessage(chatMessageElement);
                        if (!Arrays.asList(ignoredAuthers).contains(chatContent.getAuthorName())) {
                            newMessages.add(chatContent);
                        }
                    }
                }
            }

            if (!chatMessageElements.isEmpty()) {
                lastChatId = chatMessageElements.get(chatMessageElements.size() - 1).getAttribute("id");
            } else {
                System.out.println("No new messages found");
            }
            if (!newMessages.isEmpty()) {
                MessageProcessor msgProcessor = new MessageProcessor();
                msgProcessor.sendSpeakMsgs(newMessages);
                // processEachMessage(newMessages);

            }

        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
