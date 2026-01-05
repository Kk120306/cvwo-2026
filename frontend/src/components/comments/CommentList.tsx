import { useState, useMemo, type Dispatch, type SetStateAction } from "react"
import {
    Box,
    Typography,
    TextField,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
} from "@mui/material"
import CommentCard from "./CommentCard"
import type { Comment } from "../../types/globalTypes"

// Props for the List of Comment componenet 
// on vote update is a callback that is used to update the vote count of a certain comment in parent state.
// everytime a vote is made in child component which is each comment card, this changes the amount of likes/dislikes in parent state  
// ON update is a callback to update the vote counts in parent state 
// On delete is a callback to remove the comment from parent state
interface CommentListProps {
    comments: Comment[]
    postAuthorId: string
    onVoteUpdate: (
        commentId: string,
        likes: number,
        dislikes: number,
        myVote: "like" | "dislike" | null
    ) => void
    onDelete: (commentId: string) => void
    onUpdate: (commentId: string, newContent: string) => void
    setComments: Dispatch<SetStateAction<Comment[]>>
}

// Different options for sorting 
type SortOption = "recent" | "oldest" | "likes" | "dislikes"

// Compoenent for comment list to be put in each post page 
const CommentList = ({ comments, onVoteUpdate, onDelete, onUpdate, postAuthorId, setComments }: CommentListProps) => {
    const [search, setSearch] = useState("")
    const [sortBy, setSortBy] = useState<SortOption>("recent")


    // Ensures that everytime a new commment is passed, we make sure that there are no undefined values as .toLowerCase() cannot be called 
    const normalizedComments = useMemo(
        () =>
            comments.map((c) => ({
                ...c,
                content: c.content ?? "",
                author: c.author ?? {
                    id: "unknown",
                    username: "Unknown",
                },
            })),
        [comments]
    )

    // Client side search filtering and sorting
    const filteredAndSortedComments = useMemo(() => {
        const q = search.toLowerCase()

        const filtered = normalizedComments.filter((comment) => {
            if (!q) return true

            return (
                comment.content.toLowerCase().includes(q) ||
                comment.author.username.toLowerCase().includes(q)
            )
        })

        // Sort with pinned comments always at the top
        return [...filtered].sort((a, b) => {
            // Pinned comments always come first
            if (a.isPinned && !b.isPinned) return -1
            if (!a.isPinned && b.isPinned) return 1

            // If both pinned or both not pinned, use selected sort option
            switch (sortBy) {
                case "likes":
                    return b.likes - a.likes
                case "dislikes":
                    return b.dislikes - a.dislikes
                case "oldest":
                    return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
                case "recent":
                default:
                    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
            }
        })
    }, [normalizedComments, search, sortBy])

    return (
        <Box mt={4} width="100%" maxWidth={700}>
            <Typography variant="h5" gutterBottom>
                Comments ({comments.length})
            </Typography>

            {/* Search and Sort Controls */}
            {comments.length > 0 && (
                <Box display="flex" gap={2} mb={3}>
                    <TextField
                        placeholder="Search comments..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        size="small"
                        fullWidth
                    />
                    <FormControl sx={{ minWidth: 150 }} size="small">
                        <InputLabel>Sort By</InputLabel>
                        <Select
                            value={sortBy}
                            label="Sort By"
                            onChange={(e) => setSortBy(e.target.value as SortOption)}
                        >
                            <MenuItem value="recent">Most Recent</MenuItem>
                            <MenuItem value="oldest">Oldest First</MenuItem>
                            <MenuItem value="likes">Most Liked</MenuItem>
                            <MenuItem value="dislikes">Most Disliked</MenuItem>
                        </Select>
                    </FormControl>
                </Box>
            )}

            {/* Comments */}
            {filteredAndSortedComments.length === 0 ? (
                <Typography color="text.secondary">
                    {comments.length === 0
                        ? "No comments yet. Be the first to comment!"
                        : "No comments match your search."}
                </Typography>
            ) : (filteredAndSortedComments.map((comment) => (
                <CommentCard
                    key={comment.id}
                    postAuthorId={postAuthorId}
                    comment={comment}
                    onVoteUpdate={onVoteUpdate}
                    onDelete={onDelete}
                    onUpdate={onUpdate}
                    setComments={setComments}
                />
            )))}
        </Box>
    )
}

export default CommentList