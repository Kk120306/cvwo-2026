import { useState } from "react"
import { Box, Button, CircularProgress, TextField } from "@mui/material"
import { updatePost } from "../../api/handlePost"
import RichTextEditor from "../RichTextEditor"
import type { Post } from "../../types/globalTypes"
import { normalizePost } from "../../helpers/normalizer"

// Props for the update post component 
interface UpdatePostProps {
    postId: string
    initialTitle: string
    initialContent: string
    onCancel?: () => void
    newPost: (updatedPost: Post) => void
}

// Componenet pop up for updating post info
const UpdatePost = ({ postId, initialTitle, initialContent, onCancel, newPost }: UpdatePostProps) => {
    const [content, setContent] = useState(initialContent)
    const [title, setTitle] = useState(initialTitle)
    const [isUpdating, setIsUpdating] = useState(false)
    const [error, setError] = useState<string | null>(null)

    // Function to update the post 
    const handleUpdate = async () => {
        if (!content.trim()) {
            setError("Post content cannot be empty.")
            return
        }

        if (!title.trim()) {
            setError("Post title cannot be empty.")
            return
        }

        try {
            setIsUpdating(true)
            setError(null)
            const res = await updatePost(postId, title, content)
            const updatedPost = normalizePost(res) 
            // using the prop passed down that updates the state of the post in the parent component 
            newPost(updatedPost)
            // Closes the window for the form 
            onCancel?.()
        } catch (err) {
            console.error("Failed to update post:", err)
            setError("Failed to update post. Try again.")
        } finally {
            setIsUpdating(false)
        }
    }

    return (
        <Box mt={2} display="flex" flexDirection="column" gap={1}>
            <TextField
                label="Title"
                value={title}
                onChange={(e: any) => setTitle(e.target.value)}
                fullWidth
            />
            <RichTextEditor content={content} onChange={setContent} />

            {error && <Box color="error.main">{error}</Box>}

            <Box display="flex" gap={1} justifyContent="flex-end">
                {onCancel && (
                    <Button
                        variant="outlined"
                        color="secondary"
                        onClick={onCancel}
                        disabled={isUpdating}
                    >
                        Cancel
                    </Button>
                )}
                <Button
                    variant="contained"
                    color="primary"
                    onClick={handleUpdate}
                    disabled={isUpdating}
                >
                    {isUpdating ? <CircularProgress size={20} /> : "Update"}
                </Button>
            </Box>
        </Box>
    )
}

export default UpdatePost
