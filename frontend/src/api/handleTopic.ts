import { toast } from "react-hot-toast";

const baseUrl = '/api';

// Function to fetch all topics from the API 
export async function fetchAllTopics() {
    const endPoint = `${baseUrl}/topics/`;

    const res = await fetch(endPoint, {
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
    return data.topics;
}

// Function to create a new topic via the API
export const createTopic = async (name: string) => {
    const endPoint = `${baseUrl}/topics/create`;

    const res = await fetch(endPoint, {
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