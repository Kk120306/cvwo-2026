const baseUrl = '/api';

// Function that returns a presigned s3 url for uploading an image to the s3 bucket 
export async function getS3Url() {

    const endpoint = `${baseUrl}/images/s3Url`

    // Calls the endpoint to retrive the s3 upload url 
    try {
        const response = await fetch(endpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });

        if (!response.ok) {
            throw new Error('Failed to get S3 URL');
        }

        const data = await response.json();
        return data.upload_url;
    } catch (error) {
        console.error('Error fetching S3 URL:', error);
        throw error;
    }
}

// Function that uploads the file to the s3 bucket 
export async function uploadFileToS3(url: string, file: File | null) {
    if (!file) {
        console.error('No file provided for upload');
        return;
    }

    // Call the s3 signed url and then sends the file 
    try {
        await fetch(url, {
            method: 'PUT',
            body: file,
            headers: {
                'Content-Type': "multipart/form-data"
            },
        });


        // extract the image name by taking the last part of the url 
        const imageName = url.split("/").pop();
        // Since image is supported with Cloudfront CDN, use CDN as the way to retrive image url thorughout the app
        // Ensures global low latency 
        const imageUrl = "https://d1nxlczpemry9k.cloudfront.net/" + imageName;
        return imageUrl;
    } catch (error) {
        console.error('Error uploading file to S3:', error);
        throw error;
    }

}


// Function handle deletion of an image by the url 
export async function deleteImage(imageUrl: string) {

    // Extract path from URL before any queries put in place 
    const urlWithoutParams = imageUrl.split("?")[0];
    // Get the image name 
    const imageName = urlWithoutParams.split("/").pop();

    // Validate image name
    if (!imageName || imageName === 'undefined') {
        throw new Error("Invalid image URL - could not extract image name");
    }

    const endpoint = `${baseUrl}/images/delete/${imageName}`;

    const res = await fetch(endpoint, {
        method: "DELETE",
        credentials: "include",
    });

    if (!res.ok) {
        throw new Error("Failed to delete image");
    }

    return true;
}