import { toast } from "react-hot-toast";

const baseUrl = import.meta.env.VITE_BACKEND_HOST

// Function that modifies a users vote on a post
export async function votePost(
    postId: string,
    voteType: "like" | "dislike"
) {
    const endPoint = `${baseUrl}/vote/`

    const res = await fetch(endPoint, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
            votableId: postId,
            votableType: "post",
            voteType,
        }),
    })

    // Handle errors
    if (!res.ok) {
        toast.error("Couldnt Vote on Post")
        throw new Error("Failed to vote")
    }

    const data = await res.json()
    return data
}

// function that modifies a users vote on a comment 
export async function voteComment(
    commentId: string,
    voteType: "like" | "dislike"
) {
    const endPoint = `${baseUrl}/vote/`

    const res = await fetch(endPoint, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
            votableId: commentId,
            votableType: "comment",
            voteType,
        }),
    })

    // Handle errors
    if (!res.ok) {
        toast.error("Couldnt Vote on Comment")
        throw new Error("Failed to vote")
    }

    const data = await res.json()
    return data
}
