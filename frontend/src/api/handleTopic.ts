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

    return data.topics;
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
        throw new Error(err.message || "Failed to create topic");
    }

    // Parse the JSON response
    const data = await res.json();
    return data.topic;
};