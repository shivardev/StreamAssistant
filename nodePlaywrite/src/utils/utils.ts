export const sendPostRequest = async (url: string, payload: any) => {
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