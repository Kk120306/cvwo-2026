import { toast } from "react-hot-toast";
import { normalizeTopics } from "../helpers/normalizer";

// Function to fetch all topics from the API 
export async function fetchAllTopics() {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/topics/`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
    });

    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Failed to fetch topics');
    }

    // Parse the JSON response
    const data = await res.json();
    const normalized = normalizeTopics(data.topics || []);

    return normalized;
}

// Function to create a new topic via the API
export const createTopic = async (name: string) => {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/topics/create`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ name }),
    });

    if (!res.ok) {
        const err = await res.json();
        toast.error(err.message || "Failed to create topic");
        throw new Error(err.message || "Failed to create topic");
    }

    // Parse the JSON response
    const data = await res.json();
    toast.success("Topic created successfully");
    return data.topic;
};