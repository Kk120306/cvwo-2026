import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import {
    Box,
    Typography,
    Card,
    CardContent,
    CircularProgress,
    Button,
} from "@mui/material"
import { fetchPostById } from "../../api/handlePost"
import { getPostComment, createComment } from "../../api/handleComment"
import RichTextEditor from "../../components/RichTextEditor"
import { useNavigate } from "react-router-dom"
import type { Post, Comment } from "../../types/globalTypes"

// Page that shows a specific post and its comments 
const PostPage = () => {
    const { id } = useParams<{ id: string }>() // gets post Id from url
    const [post, setPost] = useState<Post | null>(null)
    const [comments, setComments] = useState<Comment[]>([])
    const [newComment, setNewComment] = useState("")
    const [loading, setLoading] = useState(true)
    const [submitting, setSubmitting] = useState(false)
    const [error, setError] = useState("")
    const navigate = useNavigate()

    // On mount and changes in id, fetch post and comments of the post 
    useEffect(() => {
        const loadData = async () => {
            try {
                if (!id) throw new Error("Invalid post ID")
                setLoading(true)
                setError("")

                const postRes = await fetchPostById(id)
                const commentRes = await getPostComment(id)

                // since postRes returns an object get the post property 
                setPost(postRes.post)
                setComments(commentRes)
            } catch (err: any) {
                setError(err.message || "Failed to load post")
            } finally {
                setLoading(false)
            }
        }
        loadData()
    }, [id])

    // Function that handles adding a new comment 
    const handleAddComment = async () => {
        // if there is no id or comment is empty
        if (!id || !newComment.trim()) return

        try {
            setSubmitting(true)
            await createComment(id, newComment)
            setNewComment("")
            // Refresh to show new comment by mounting compoenent 
            navigate(0)
        } catch {
            alert("Failed to post comment")
        } finally {
            setSubmitting(false)
        }
    }

    if (loading) return <p>Loading...</p>
    if (error) return <Typography color="error">{error}</Typography>
    if (!post) return <Typography>No post found</Typography>

    return (
        <Box mt={4} display="flex" flexDirection="column" alignItems="center">
            {/* Post */}
            <Card sx={{ width: "100%", maxWidth: 700 }}>
                <CardContent>
                    <Typography variant="h4" fontWeight={600} gutterBottom>
                        {post.title}
                    </Typography>

                    <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                        {post.topic.name} • by {post.author.username}
                    </Typography>

                    <Typography
                        variant="body1"
                        dangerouslySetInnerHTML={{ __html: post.content }}
                    />

                    <Box mt={2} display="flex" gap={2}>
                        <Typography color="primary">👍 {post.likes}</Typography>
                        <Typography color="error">👎 {post.dislikes}</Typography>
                    </Box>
                </CardContent>
            </Card>

            {/* Add Comment */}
            <Box mt={4} width="100%" maxWidth={700}>
                <Typography variant="h6" gutterBottom>
                    Add a comment
                </Typography>

                {/* Use Rich text editor commponent from tiptap */}
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

            {/* Comments */}
            <Box mt={4} width="100%" maxWidth={700}>
                <Typography variant="h5" gutterBottom>
                    Comments
                </Typography>

                {comments.length === 0 ? (
                    <Typography>No comments yet.</Typography>
                ) : (
                    comments.map((comment) => (
                        <Card key={comment.id} sx={{ mb: 2 }}>
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

                                {/* TODO : add like and dislike logic conenction  */}
                                <Box mt={1} display="flex" gap={2}>
                                    <Typography color="primary">
                                        👍 {comment.likes}
                                    </Typography>
                                    <Typography color="error">
                                        👎 {comment.dislikes}
                                    </Typography>
                                </Box>
                            </CardContent>
                        </Card>
                    ))
                )}
            </Box>
        </Box>
    )
}

export default PostPage
