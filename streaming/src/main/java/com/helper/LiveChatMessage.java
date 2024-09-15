package com.helper;

import org.openqa.selenium.By;
import org.openqa.selenium.WebElement;

public class LiveChatMessage {

    final String chatid;
    final String authorName;
    final String authorPhotoUrl;
    private String timestamp;
    final String messageContent;

    // Constructors
    public LiveChatMessage(WebElement chatElement) {
        this.chatid = chatElement.getAttribute("id");
        this.authorPhotoUrl = chatElement.findElement(By.cssSelector("yt-img-shadow img")).getAttribute("src");
        this.authorName = chatElement.findElement(By.cssSelector("yt-live-chat-author-chip #author-name")).getText();
        this.messageContent = chatElement.findElement(By.cssSelector("#message")).getText();
    }

    @Override
    public String toString() {
        return "LiveChatMessage{"
                + "chatid='" + chatid + '\''
                + ", authorName='" + authorName + '\''
                + ", authorPhotoUrl='" + authorPhotoUrl + '\''
                + ", timestamp='" + timestamp + '\''
                + ", messageContent='" + messageContent + '\''
                + '}';
    }
    public String getAuthorName() {
        return authorName;
    }
    public String getchatId() {
        return chatid;
    }
    public String getMessage(){
        return messageContent;
    }
}
