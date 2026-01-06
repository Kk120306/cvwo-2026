import { useState } from "react"
import { useNavigate } from "react-router-dom"
import {
    Card,
    CardContent,
    Typography,
    Box,
    IconButton,
} from "@mui/material"
import ThumbUpAltOutlinedIcon from "@mui/icons-material/ThumbUpAltOutlined"
import ThumbDownAltOutlinedIcon from "@mui/icons-material/ThumbDownAltOutlined"
import { voteComment } from "../../api/handleVote"
import { useAppSelector } from "../../hooks/reduxHooks"
import type { Comment } from "../../types/globalTypes"
import { deleteComment } from "../../api/handleComment"
import UpdateComment from "./CommentUpdate"


// Componenet prop for each comment card. 
// On Vote Update callback that updates the vote counts in parent state 
// On Delete callback that removes the comment from parent state
// On update callback that updates the comment content in parent state
interface CommentCardProps {
    comment: Comment
    onVoteUpdate?: (commentId: string, likes: number, dislikes: number, myVote: "like" | "dislike" | null) => void
    onDelete?: (commentId: string) => void
    onUpdate?: (commentId: string, newContent: string) => void
}

// Comment card componenet 
const CommentCard = ({ comment, onVoteUpdate, onDelete, onUpdate }: CommentCardProps) => {
    const [isVoting, setIsVoting] = useState(false)
    const [isEditing, setIsEditing] = useState(false)

    const navigate = useNavigate()
    const user = useAppSelector(state => state.auth.user)

    // Function to handle the vote of each comment 
    const handleVote = async (type: "like" | "dislike") => {
        if (onVoteUpdate === undefined) return
        if (!user) {
            navigate("/login")
            return
        }

        if (isVoting) return // ensures that user cant spam click 

        try {
            setIsVoting(true)
            const res = await voteComment(comment.id, type)
            // uses the results to pass the new values to store comment data 
            onVoteUpdate(comment.id, res.likes, res.dislikes, res.myVote)
        } catch (err) {
            console.error("Failed to vote:", err)
        } finally {
            setIsVoting(false)
        }
    }

    const handleDelete = async (commentId: string) => {
        if (onDelete === undefined) return
        try {
            const confirmed = window.confirm("Are you sure you want to delete this Comment?")
            if (!confirmed) return
            await deleteComment(commentId);
            onDelete(commentId)
        }
        catch (err) {
            console.error("Failed to delete comment:", err)
        }
    }


    return (
        <Card sx={{ mb: 2 }}>
            <CardContent>
                <Typography
                    variant="subtitle2"
                    color="text.secondary"
                    gutterBottom
                >
                    by {comment.author.username}
                </Typography>

                <Typography
                    variant="body1"
                    dangerouslySetInnerHTML={{
                        __html: comment.content,
                    }}
                />

                {/* Voting */}
                {onVoteUpdate && (
                    <Box mt={2} display="flex" gap={1} alignItems="center">
                        <IconButton
                            size="small"
                            color={comment.myVote === "like" ? "primary" : "default"}
                            onClick={() => handleVote("like")}
                            disabled={isVoting}
                        >
                            <ThumbUpAltOutlinedIcon fontSize="small" />
                        </IconButton>

                        <Typography fontWeight={comment.myVote === "like" ? 600 : 400}>
                            {comment.likes}
                        </Typography>

                        <IconButton
                            size="small"
                            color={comment.myVote === "dislike" ? "error" : "default"}
                            onClick={() => handleVote("dislike")}
                            disabled={isVoting}
                        >
                            <ThumbDownAltOutlinedIcon fontSize="small" />
                        </IconButton>

                        <Typography fontWeight={comment.myVote === "dislike" ? 600 : 400}>
                            {comment.dislikes}
                        </Typography>
                    </Box>
                )}
                <Box>
                    {((user?.isAdmin || comment.author.id === user?.id) && onDelete && onUpdate) && (
                        <Box>
                            <Typography
                                variant="caption"
                                color="error"
                                onClick={(e) => {
                                    e.stopPropagation()
                                    handleDelete(comment.id)
                                }}
                            >
                                Delete Comment
                            </Typography>
                            {(isEditing && onUpdate) ? (
                                <UpdateComment
                                    commentId={comment.id}
                                    initialContent={comment.content}
                                    onCancel={() => setIsEditing(false)}
                                    onUpdate={(newContent) => {
                                        onUpdate(comment.id, newContent)
                                        setIsEditing(false)
                                    }}
                                />
                            ) : (
                                <Typography
                                    variant="caption"
                                    color="primary"
                                    onClick={(e) => {
                                        e.stopPropagation()
                                        setIsEditing(true)
                                    }}
                                    sx={{ cursor: 'pointer' }}
                                >
                                    Update Comment
                                </Typography>
                            )}
                        </Box>
                    )}
                </Box>
            </CardContent>
        </Card>
    )
}

export default CommentCard