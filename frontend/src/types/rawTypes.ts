// All types here represent the exact structure of the API responses
export interface RawUser {
    ID: string;
    Username: string;
    AvatarURL: string;
    IsAdmin: boolean;
}

export interface RawTopic {
    ID: string;
    Name: string;
    Slug: string;
}

export interface RawPost {
    ID: string;
    Title: string;
    Content: string;
    IsPinned: boolean;
    CreatedAt: string;
    UpdatedAt: string;
    Author: RawUser;
    Topic: RawTopic;
    Likes: number;
    Dislikes: number;
    MyVote: "like" | "dislike" | null;
}


export interface RawComment {
    ID: string
    Content: string
    CreatedAt: string
    Likes: number
    Dislikes: number
    Author: RawUser
}

