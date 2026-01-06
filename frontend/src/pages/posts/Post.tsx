import { useEffect, useState } from "react"
import { useParams, useNavigate, Link } from "react-router-dom"
import {
    Box,
    Typography,
    Card,
    CardContent,
    Button,
    IconButton,
    Tooltip,
    Container,
} from "@mui/material"
import ThumbUpAltOutlinedIcon from "@mui/icons-material/ThumbUpAltOutlined"
import ThumbDownAltOutlinedIcon from "@mui/icons-material/ThumbDownAltOutlined"
import { deletePost, fetchPostById } from "../../api/handlePost"
import { getPostComment, createComment } from "../../api/handleComment"
import { votePost } from "../../api/handleVote"
import RichTextEditor from "../../components/provider/RichTextEditor"
import CommentList from "../../components/comments/CommentList"
import { useAppSelector } from "../../hooks/reduxHooks"
import type { Post, Comment } from "../../types/globalTypes"
import UpdatePost from "../../components/post/PostUpdate"
import ShareIcon from '@mui/icons-material/Share';
import { sharePost } from "../../helpers/share"
import { deleteImage } from "../../api/handleImage"


// Page that shows a specific post and its comments 
const PostPage = () => {
    const { id: postId } = useParams<{ id: string }>() // gets post Id from url

    const [post, setPost] = useState<Post | null>(null)
    const [comments, setComments] = useState<Comment[]>([])
    const [newComment, setNewComment] = useState("")
    const [loading, setLoading] = useState(true)
    const [submitting, setSubmitting] = useState(false)
    const [error, setError] = useState("")
    const [isEditing, setIsEditing] = useState(false)
    const navigate = useNavigate()
    const user = useAppSelector(state => state.auth.user)

    // Changes in id, fetch post and comments of the post 
    useEffect(() => {
        const loadData = async () => {
            try {
                if (!postId) throw new Error("Invalid post ID")
                setLoading(true)
                setError("")

                const postRes = await fetchPostById(postId)
                const commentRes = await getPostComment(postId)

                // since postRes returns an object get the post property 
                setPost(postRes)
                setComments(commentRes)
            } catch (err: unknown) {
                if (err instanceof Error) {
                    setError(err.message); 
                } else {
                    setError("Failed to load post"); 
                }
            } finally {
                setLoading(false)
            }
        }
        loadData()
    }, [postId])

    // Function that handles post voting
    const handleVote = async (type: "like" | "dislike") => {
        if (!post) return

        const res = await votePost(post.id, type)

        // Takes most recent post state and updates vote info with likes, dislikes, myVote
        // Vote post API already ensures post is under the postId
        setPost(prev =>
            prev
                ? {
                    ...prev,
                    likes: res.likes,
                    dislikes: res.dislikes,
                    myVote: res.myVote,
                }
                : prev
        )
    }

    // Function that handles deleting a post using its Id
    const handleDelete = async (postId: string) => {
        try {
            const confirmed = window.confirm("Are you sure you want to delete this post? This will delete all comments under it.")
            if (!confirmed) return
            setLoading(true)

            await deletePost(postId)
            // If the post has an image, delete from bucket
            if (post?.imageUrl) {
                await deleteImage(post.imageUrl);
            }
            navigate(-1)
        } catch {
            console.error("Failed to delete post")
        } finally {
            setLoading(false)
        }
    }

    // Function that handles adding a new comment 
    const handleAddComment = async () => {
        // if there is no id or comment is empty
        if (!postId || !newComment.trim()) return

        try {
            setSubmitting(true)
            const created = await createComment({ postId, newComment })
            setNewComment("<p></p>")
            setComments(prev => [created.comment, ...prev])
        } catch {
            alert("Failed to post comment")
        } finally {
            setSubmitting(false)
        }
    }

    // Function to update comment votes
    const handleCommentVoteUpdate = (
        commentId: string,
        likes: number,
        dislikes: number,
        myVote: "like" | "dislike" | null
    ) => {
        // Take the most recent comments array and update comment that matches comment id with 
        // new likes, dislikes, myVote
        setComments(prev =>
            prev.map((comment): Comment =>
                comment.id === commentId
                    ? { ...comment, likes, dislikes, myVote }
                    : comment
            )
        )
    }


    // Function that is passed to CommentUpdate to delete comments from state 
    const handleCommentDelete = (commentId: string) => {
        setComments(prev => prev.filter(c => c.id !== commentId))
    }

    // function that handles updating comment content to state 
    const handleCommentUpdate = (commentId: string, newContent: string) => {
        setComments(prev =>
            prev.map(c => c.id === commentId ? { ...c, content: newContent } : c)
        )
    }

    // Different rendering conditions 
    if (loading) return <Typography>Loading...</Typography>
    if (error) return <Typography color="error">{error}</Typography>
    if (!post) return <Typography>No post found</Typography>

    return (
        <Container sx={{ mt: 4, mb: 4 }}>
            <Box mt={4} display="flex" flexDirection="column" alignItems="center" >
                <Card sx={{ width: "100%" }}>
                    <CardContent>
                        <Typography variant="h4" fontWeight={600} gutterBottom>
                            {post.title}
                        </Typography>

                        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                            {post.topic?.name || 'Unknown Topic'} • by{' '}
                            {post.author ? (
                                <Link
                                    to={`/profile/${post.author.username}`}
                                    color="yellow"
                                    style={{ textDecoration: 'none' }}
                                >
                                    {post.author.username}
                                </Link>
                            ) : (
                                'Unknown User'
                            )}
                        </Typography>


                        {/* Share button */}
                        <Tooltip title="Share">
                            <IconButton
                                onClick={() => sharePost(post.id, post.title)}
                                size="small"
                            >
                                <ShareIcon fontSize="small" />
                            </IconButton>
                        </Tooltip>

                        {/* Will display the image directly from CDN if post has a image */}
                        <Box>
                            {post.imageUrl && (
                                <Box
                                    component="img"
                                    src={post.imageUrl}
                                    alt="Post Image"
                                    sx={{
                                        width: "100%",
                                        maxHeight: 400,
                                        objectFit: "cover",
                                        borderRadius: 2,
                                        my: 2,
                                    }}
                                />
                            )}
                        </Box>

                        <Typography
                            variant="body1"
                            dangerouslySetInnerHTML={{ __html: post.content }} // TipTap rich text editor content
                        />

                        {/* Voting */}
                        <Box mt={2} display="flex" gap={1} alignItems="center">
                            <IconButton
                                size="small"
                                color={post.myVote === "like" ? "primary" : "default"}
                                onClick={() => {
                                    if (!user) {
                                        navigate("/login")
                                        return
                                    }
                                    handleVote("like")
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
                                onClick={() => {
                                    if (!user) {
                                        navigate("/login")
                                        return
                                    }
                                    handleVote("dislike")
                                }}
                            >
                                <ThumbDownAltOutlinedIcon fontSize="small" />
                            </IconButton>

                            <Typography fontWeight={post.myVote === "dislike" ? 600 : 400}>
                                {post.dislikes}
                            </Typography>
                        </Box>
                        <Box>
                            {(user?.isAdmin || post.author?.id === user?.id) && (
                                <Box>
                                    <Typography
                                        variant="caption"
                                        color="error"
                                        onClick={(e) => {
                                            e.stopPropagation()
                                            handleDelete(post.id)
                                        }}
                                        sx={{ cursor: 'pointer', mr: 2 }}
                                    >
                                        Delete Post
                                    </Typography>
                                    {/* If user has clicked edit, shows a pop up component to edit post, not a new page */}
                                    {isEditing ? (
                                        <UpdatePost
                                            postId={post.id}
                                            initialTitle={post.title}
                                            initialContent={post.content}
                                            initialImage={post.imageUrl}
                                            onCancel={() => setIsEditing(false)}
                                            newPost={(updatedPost: Post) => setPost(updatedPost)}
                                        />
                                    ) : (
                                        <Typography
                                            variant="caption"
                                            color="primary"
                                            onClick={() => setIsEditing(true)}
                                            sx={{ cursor: 'pointer' }}
                                        >
                                            Update Post
                                        </Typography>
                                    )}
                                </Box >
                            )}
                        </Box>
                    </CardContent>
                </Card>

                {/* Add Comment */}
                <Box mt={4} width="100%">
                    <Typography variant="h6" gutterBottom>
                        Add a comment
                    </Typography>
                    {/* Use Rich text editor component from tiptap */}
                    <RichTextEditor
                        content={newComment}
                        onChange={setNewComment}
                    />

                    <Button
                        variant="contained"
                        sx={{ mt: 2 }}
                        onClick={handleAddComment}
                        disabled={submitting || !newComment.trim()}
                    >
                        {submitting ? "Posting..." : "Post Comment"}
                    </Button>
                </Box>

                {/* Comments are shown through a component */}
                <CommentList
                    comments={comments}
                    onVoteUpdate={handleCommentVoteUpdate}
                    onDelete={handleCommentDelete}
                    onUpdate={handleCommentUpdate}
                />
            </Box>
        </Container>
    )
}

export default PostPage