import type { User } from "../types/globalTypes";
// Helper function to normlize user data comming in from the api
// Removes any unnecessary fields and maps this to our frontend 
export function normalizeUser(user: User) {
    return {
        id: user.id,
        username: user.username,
        avatarURL: user.avatarURL,
        isAdmin: user.isAdmin,
    };
}
