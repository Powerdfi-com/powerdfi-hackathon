export const uploadFile = (files: FileList) =>
  new Promise<string[]>(async (resolve, reject) => {
    try {
      const urls = [];
      const signResponse = await fetch("/api/signature");
      const signData = await signResponse.json();
      const url =
        "https://api.cloudinary.com/v1_1/" +
        signData.cloudname +
        "/auto/upload";
      for (let i = 0; i < files.length; i++) {
        const file = files.item(i);
        const formData = new FormData();
        formData.append("file", file!);
        formData.append("api_key", signData.apiKey);
        formData.append("timestamp", signData.timestamp);
        formData.append("signature", signData.signature);

        const res = await fetch(url, {
          method: "POST",
          body: formData,
        });
        if (!res.ok) {
          reject();
        }
        const uploadUrl = await res.json();
        urls.push(uploadUrl.url);
      }
      resolve(urls);
    } catch (e) {
      reject();
    }
  });
