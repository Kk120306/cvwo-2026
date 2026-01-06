import { toast } from "react-hot-toast";


const baseUrl = import.meta.env.VITE_BACKEND_HOST;

// function that fetches posts based on topic from api 
export async function fetchPostByTopic(topicSlug: string) {

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
    return data.posts;
}


// Function that fetches a post by a specific post Id 
export async function fetchPostById(postId: string) {
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
    return data.post;
}

// funciton that creates a post under a topic which is identified by topicSlug
export async function createPost(postData: { title: string; content: string, topicSlug: string, imageUrl?: string | null }) {
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

// function that deletes a post by post Id
export async function deletePost(postId: string) {
    const endpoint = `${baseUrl}/posts/delete/${postId}`;

    const res = await fetch(endpoint, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (!res.ok) {
        toast.error("Failed to delete post");
        throw new Error("Failed to delete post");
    }

    toast.success("Post deleted successfully");
    return true;
}

// function that updates a post by post Id
export async function updatePost(postData: { title: string; content: string, postId: string, imageUrl?: string | null }) {
    const endpoint = `${baseUrl}/posts/update/${postData.postId}`;

    const res = await fetch(endpoint, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(postData),
    });

    if (!res.ok) {
        toast.error("Failed to update post");
        throw new Error("Failed to update post");
    }

    const data = await res.json();
    toast.success("Post updated successfully");
    return data.post;
}

// function that changes the pin status of a post 
export async function managePostPin(postId: string, isPinned: boolean) {
    const endpoint = `${baseUrl}/posts/pin/${postId}`;
    console.log(endpoint);

    const res = await fetch(endpoint, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ isPinned }),
    });

    if (!res.ok) {
        toast.error(`Failed to ${isPinned ? "pin" : "unpin"} post`);
        throw new Error(`Failed to ${isPinned ? "pin" : "unpin"} post`);
    }

    toast.success(`Post ${isPinned ? "pinned" : "unpinned"} successfully`);
    return true;
}