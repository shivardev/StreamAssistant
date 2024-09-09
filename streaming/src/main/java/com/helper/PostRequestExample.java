package com.helper;

import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;

import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;

public class PostRequestExample {

    private Gson gson = new Gson();

    public void makePostRequest(List<LiveChatMessage> messages) {
        try {
            // URL of the endpoint
            URL url = new URL("http://localhost:3000/takemsgs");

            // Open a connection to the URL
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json; utf-8");
            connection.setRequestProperty("Accept", "application/json");
            connection.setDoOutput(true);

            // Convert the list of messages to a JSON string
            JsonArray jsonArray = gson.toJsonTree(messages).getAsJsonArray();
            JsonObject jsonObject = new JsonObject();
            jsonObject.add("messages", jsonArray);
            String jsonPayload = gson.toJson(jsonObject);
            System.out.println("JSON Payload: " + jsonPayload);

            // Send the request
            try (OutputStream os = connection.getOutputStream()) {
                byte[] input = jsonPayload.getBytes("utf-8");
                os.write(input, 0, input.length);
            }

            // Get the response
            int responseCode = connection.getResponseCode();
            System.out.println("Response Code: " + responseCode);

            connection.disconnect();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
