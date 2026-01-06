import {
    Box,
    Typography,
    Card,
    CardContent,
    IconButton,
    Tooltip,
} from "@mui/material"
import { useNavigate } from "react-router-dom"
import ThumbUpAltOutlinedIcon from '@mui/icons-material/ThumbUpAltOutlined';
import ThumbDownAltOutlinedIcon from '@mui/icons-material/ThumbDownAltOutlined';
import ShareIcon from '@mui/icons-material/Share';
import PushPinIcon from '@mui/icons-material/PushPin';
import type { Post } from "../../types/globalTypes"
import { sharePost } from "../../helpers/share"

interface PostCardProps {
    post: Post
    isAdmin?: boolean
    onVote?: (postId: string, type: "like" | "dislike") => void
    onPin?: (postId: string, pin: boolean) => void
}

export default function PostCard({ post, isAdmin, onVote, onPin }: PostCardProps) {
    const navigate = useNavigate()

    return (
        <Card
            variant="outlined"
            sx={{ cursor: "pointer" }}
            onClick={() => navigate(`/posts/${post.id}`)}
        >
            <CardContent>
                {/* Pin Section - Visible to all, toggle only for admins */}
                {(post.isPinned || isAdmin) && (
                    <Box display="flex" alignItems="center" mb={1} gap={1}>
                        {isAdmin && onPin ? (
                            <IconButton
                                size="small"
                                color={post.isPinned ? "error" : "default"}
                                onClick={(e) => {
                                    e.stopPropagation()
                                    onPin(post.id, !post.isPinned)
                                }}
                            >
                                <PushPinIcon fontSize="small" />
                            </IconButton>
                        ) : (
                            post.isPinned && <PushPinIcon fontSize="small" color="error" />
                        )}
                        <Typography variant="caption">
                            {post.isPinned ? "Pinned by Admin" : "Click to Pin Post"}
                        </Typography>
                    </Box>
                )}


                {/* Image */}
                {post.imageUrl && (
                    <Box mb={2}>
                        <img
                            src={post.imageUrl}
                            alt="Post Image"
                            style={{ maxWidth: '100%', borderRadius: 8 }}
                        />
                    </Box>
                )}

                {/* Title */}
                <Typography variant="h6" gutterBottom>
                    {post.title}
                </Typography>

                {/* Metadata */}
                <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <Typography variant="caption" color="text.secondary">
                        {post.topic.name}
                    </Typography>

                    <Typography variant="caption" color="text.secondary">
                        • Posted by {post.author.username}
                    </Typography>

                    <Tooltip title="Share">
                        <IconButton
                            onClick={(e) => {
                                e.stopPropagation()
                                sharePost(post.id, post.title)
                            }}
                            size="small"
                        >
                            <ShareIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </Box>

                {/* Content */}
                <Typography
                    variant="body2"
                    dangerouslySetInnerHTML={{ __html: post.content }}
                />

                {/* Vote Section */}
                {onVote &&
                    (<Box mt={2} display="flex" gap={1} alignItems="center">
                        <IconButton
                            size="small"
                            color={post.myVote === "like" ? "primary" : "default"}
                            onClick={(e) => {
                                e.stopPropagation()
                                onVote(post.id, "like")
                            }}
                        >
                            <ThumbUpAltOutlinedIcon fontSize="small" />
                        </IconButton>
                        <Typography fontWeight={post.myVote === "like" ? 600 : 400}>
                            {post.likes}
                        </Typography>

                        <IconButton
                            size="small"
                            color={post.myVote === "dislike" ? "error" : "default"}
                            onClick={(e) => {
                                e.stopPropagation()
                                onVote(post.id, "dislike")
                            }}
                        >
                            <ThumbDownAltOutlinedIcon fontSize="small" />
                        </IconButton>
                        <Typography fontWeight={post.myVote === "dislike" ? 600 : 400}>
                            {post.dislikes}
                        </Typography>
                    </Box>
                    )}
            </CardContent>
        </Card>
    )
}