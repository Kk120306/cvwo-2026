import { normalizeComments } from "../helpers/normalizer"
import { toast } from "react-hot-toast"

const baseUrl = import.meta.env.VITE_BACKEND_HOST

// function to get all comments under a post with postId
export async function getPostComment(postId: string) {
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
    console.log(data);
    return normalizeComments(data.comments || [])
}


// function that creates a comment under a postID
export async function createComment(commentData: { postId: string, newComment: string}) {
    const endpoint = `${baseUrl}/comments/create/${commentData.postId}`

    const res = await fetch(endpoint, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ commentData }),
    })
    console.log(res);

    if (!res.ok) {
        toast.error("Failed to create comment")
        throw new Error("Failed to create comment")
    }

    toast.success("Comment created successfully")
    const data = await res.json()
    return data;
}


// function that deletes a comment by commentId
export async function deleteComment(commentId: string) {
    const endpoint = `${baseUrl}/comments/delete/${commentId}`

    const res = await fetch(endpoint, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    })

    if (!res.ok) {
        toast.error("Failed to delete comment")
        throw new Error("Failed to delete comment")
    }

    toast.success("Comment deleted successfully")
}



// function that updates a comment by commentId
export async function updateComment(commentData: { commentId: string, content: string}) {
    const endpoint = `${baseUrl}/comments/update/${commentData.commentId}`

    const res = await fetch(endpoint, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ commentData }),
    })

    if (!res.ok) {
        toast.error("Failed to update comment")
        throw new Error("Failed to update comment")
    }

    toast.success("Comment updated successfully")
    const data = await res.json()
    return data

}