import { normalizeComments } from "../helpers/normalizer"
import { toast } from "react-hot-toast"

// function to get all comments under a post with postId
export async function getPostComment(postId: string) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST
    const endpoint = `${baseUrl}/comments/post/${postId}`

    const res = await fetch(endpoint, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    })

    if (!res.ok) {
        throw new Error("Failed to fetch comments")
    }

    const data = await res.json()
    return normalizeComments(data.comments || [])
}


// function that creates a comment under a postID
export async function createComment(postId: string, content: string) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST
    const endpoint = `${baseUrl}/comments/create/${postId}`

    const res = await fetch(endpoint, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ content }),
    })

    if (!res.ok) {
        toast.error("Failed to create comment")
        throw new Error("Failed to create comment")
    }

    toast.success("Comment created successfully")
    return await res.json()
}