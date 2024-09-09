package com.streaming;

import java.io.IOException;
import java.util.List;

import com.google.api.client.http.HttpTransport;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.JsonFactory;
import com.google.api.client.json.gson.GsonFactory;
import com.google.api.services.youtube.YouTube;
import com.google.api.services.youtube.model.LiveBroadcast;
import com.google.api.services.youtube.model.LiveBroadcastListResponse;

public class YouTubeLiveFetcher {

    private static final String APPLICATION_NAME = "streaming";
    private static final JsonFactory JSON_FACTORY = GsonFactory.getDefaultInstance();
    private static final HttpTransport HTTP_TRANSPORT = new NetHttpTransport();
    private static final String API_KEY = "AIzaSyBS1k9cg22QT7R6o2gtd_vrDWrBjsuYuQY";

    private static YouTube getService() {
        return new YouTube.Builder(HTTP_TRANSPORT, JSON_FACTORY, httpRequest -> {})
                .setApplicationName(APPLICATION_NAME)
                .build();
    }

    public static String fetchURL() {
        try {
            YouTube youtubeService = getService();
            YouTube.LiveBroadcasts.List request = youtubeService.liveBroadcasts()
                    .list("snippet,contentDetails,status");
            request.setBroadcastType("all");
            request.setMine(true); // Set this to true if you want to fetch broadcasts from the authenticated user
            request.setKey(API_KEY);

            LiveBroadcastListResponse response = request.execute();
            List<LiveBroadcast> broadcasts = response.getItems();

            for (LiveBroadcast broadcast : broadcasts) {
                if ("live".equals(broadcast.getStatus().getLifeCycleStatus())) {
                    String videoId = broadcast.getId();
                    String liveVideoLink = "https://www.youtube.com/watch?v=" + videoId;
                    System.out.println("Live Video Link: " + liveVideoLink);
                    return liveVideoLink;
                }
            }
            return "";
        } catch (IOException e) {
            e.printStackTrace();
            return "";
        }
    }
}
