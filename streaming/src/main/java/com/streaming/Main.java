package com.streaming;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeDriver;

import com.helper.LiveChatMessage;
import com.helper.PostRequestExample;

import io.github.bonigarcia.wdm.WebDriverManager;

public class Main {


    static String lastChatId = "";
    static String chatUrl = "";
    static String[] ignoredAuthers = new String[]{"Nightbot", "YouTube"};
    static final String ChannelLiveUrl = "https://www.youtube.com/channel/UCPFM_Ug62Ei3CUfvquG4KOg/live";

    public static void main(String[] args) {
        // Fetch the live broadcase from youtube API and get the live video link
        WebDriver driver;
        WebDriverManager.chromedriver().setup();
        driver = new ChromeDriver();

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
                PostRequestExample postRequestExample = new PostRequestExample();
                postRequestExample.makePostRequest(newMessages);
            }
            
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
