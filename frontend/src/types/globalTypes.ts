// Normalised user type for frontend use
export interface User {
    id: string;
    username: string;
    avatarURL: string;
    isAdmin: boolean;
}

// Topic type
export interface Topic {
    id: string;
    name: string;
    slug: string;
}

// Post type
export interface Post {
    id: string;
    title: string;
    content: string;
    isPinned: boolean;
    createdAt: string;
    updatedAt: string;
    author: User;
    topic: Topic;
}

