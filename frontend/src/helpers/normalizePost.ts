// src/helpers/normalizePost.ts
import type { Post } from "../types/globalTypes";
import type { RawPost } from "../types/rawTypes";

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
    };
}

// Call for an array of posts
export function normalizePosts(dataArr: RawPost[]): Post[] {
    return dataArr.map(normalizePost);
}
