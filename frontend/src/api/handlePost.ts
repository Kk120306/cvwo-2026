import { normalizePosts } from "../helpers/normalizer";
import { toast } from "react-hot-toast";

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
    console.log(data);
    // Ensures Post type is met 
    const posts = normalizePosts(data.posts || []);

    return posts;
}


// Function that fetches a post by a specific post Id 
export async function fetchPostById(postId: string) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST;

    const endpoint = `${baseUrl}/posts/id/${postId}`;

    const res = await fetch(endpoint, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!res.ok) {
        throw new Error("Failed to fetch post");
    }

    const data = await res.json();
    // Ensures Post type is met 
    const posts = normalizePosts([data.post || {}]);

    return { post: posts[0] };
}

// funciton that creates a post under a topic which is identified by topicSlug
export async function createPost(postData: { title: string; content: string, topicSlug: string }) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST;

    const endpoint = `${baseUrl}/posts/create/${postData.topicSlug}`;

    const res = await fetch(endpoint, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(postData),
    });

    if (!res.ok) {
        toast.error("Failed to create post");
        throw new Error("Failed to create post");
    }

    toast.success("Post created successfully");
    return true;
}