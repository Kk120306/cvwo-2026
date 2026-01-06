// Normalised user type for frontend use - used once API data is cleaned

export interface User {
    id: string;
    username: string;
    avatarURL: string;
    isAdmin: boolean;
}


export interface Topic {
    id: string;
    name: string;
    slug: string;
}


export interface Post {
    id: string;
    title: string;
    content: string;
    isPinned: boolean;
    createdAt: string;
    updatedAt: string;
    author: User;
    topic: Topic;
    likes: number;
    dislikes: number;
    myVote?: "like" | "dislike" | null
    imageUrl: string | null
}


export interface Comment {
    id: string
    content: string
    createdAt: string
    likes: number
    dislikes: number
    author: User
    myVote?: "like" | "dislike" | null
}


export interface UserProfile {
    id: string;
    username: string;
    avatarUrl: string;
    isAdmin: boolean;
    postCount: number;
    commentCount: number;
    createdAt: string;
    posts?: Post[];
    comments?: Comment[];
}