
import type { UserProfile } from "../types/globalTypes";

const baseUrl = '/api';

export async function getUserProfile(username: string, includePosts = false, includeComments = false) {
    const params = new URLSearchParams();
    if (includePosts) params.append('posts', 'true');
    if (includeComments) params.append('comments', 'true');

    const url = `${baseUrl}/user/profile/${username}${params.toString() ? '?' + params : ''}`;

    const response = await fetch(url, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    });

    if (!response.ok) throw new Error("Failed to load user");

    const data = await response.json();
    console.log(data);
    return data.user as UserProfile;
}