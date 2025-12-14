import { normalizePosts } from "../helpers/normalizePost";

// function that fetches posts based on topic from api 
export async function fetchPostByTopic(topicSlug: string) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST;

    // If all we call a different endpoint 
    const endpoint =
        topicSlug === "all"
            ? `${baseUrl}/posts/all`
            : `${baseUrl}/posts/topic/${topicSlug}`;

    const res = await fetch(endpoint, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!res.ok) {
        throw new Error("Failed to fetch posts");
    }
    // Parse the JSON response
    const data = await res.json();
    // Ensures Post type is met 
    const posts = normalizePosts(data.posts || []);

    return posts;
}
