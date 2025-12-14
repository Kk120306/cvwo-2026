import type { Comment, Post, Topic, User } from "../types/globalTypes"
import type { RawComment, RawPost, RawTopic } from "../types/rawTypes"

// These functions convert raw API data to frontend types

// Function that normalizes a single raw comment to Comment type
export function normalizeComment(data: RawComment): Comment {
    return {
        id: data.ID,
        content: data.Content,
        createdAt: data.CreatedAt,
        likes: data.Likes,
        dislikes: data.Dislikes,
        author: {
            id: data.Author.ID,
            username: data.Author.Username,
            avatarURL: data.Author.AvatarURL,
            isAdmin: data.Author.IsAdmin,
        },
    }
}

export function normalizeComments(dataArr: RawComment[]): Comment[] {
    return dataArr.map(normalizeComment)
}

// Function that converts a single raw post to Post type
export function normalizePost(data: RawPost): Post {
    return {
        id: data.ID,
        title: data.Title,
        content: data.Content,
        isPinned: data.IsPinned,
        createdAt: data.CreatedAt,
        updatedAt: data.UpdatedAt,
        author: {
            id: data.Author.ID,
            username: data.Author.Username,
            avatarURL: data.Author.AvatarURL,
            isAdmin: data.Author.IsAdmin,
        },
        topic: {
            id: data.Topic.ID,
            name: data.Topic.Name,
            slug: data.Topic.Slug,
        },
        likes: data.Likes,
        dislikes: data.Dislikes,
    };
}

// Call for an array of posts
export function normalizePosts(dataArr: RawPost[]): Post[] {
    return dataArr.map(normalizePost);
}

// Removes any unnecessary fields and maps this to our frontend 
export function normalizeUser(user: User) {
    return {
        id: user.id,
        username: user.username,
        avatarURL: user.avatarURL,
        isAdmin: user.isAdmin,
    };
}

// Function that converts a single raw topic to Topic type
export function normalizeTopic(data: RawTopic): Topic {
    return {
        id: data.ID,
        name: data.Name,
        slug: data.Slug,
    }
}

export function normalizeTopics(dataArr: RawTopic[]): Topic[] {
    return dataArr.map(normalizeTopic)
}