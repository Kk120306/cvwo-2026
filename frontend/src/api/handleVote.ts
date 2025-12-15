import {toast} from "react-hot-toast";

// Function that modifies a users vote on a post
export async function votePost(
    postId: string,
    voteType: "like" | "dislike"
) {
    const baseUrl = import.meta.env.VITE_BACKEND_HOST

    const res = await fetch(`${baseUrl}/vote/`, {
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

    return res.json()
}
