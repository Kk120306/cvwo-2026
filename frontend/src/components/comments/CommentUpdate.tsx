import { useState } from "react"
import { Box, Button, CircularProgress } from "@mui/material"
import { updateComment } from "../../api/handleComment"
import RichTextEditor from "../provider/RichTextEditor"
import { toast } from "react-hot-toast"

// Props for the udpate comment compoenent 
// On update is passed by the parent component so changes are made to the state in parent 
interface UpdateCommentProps {
    commentId: string
    initialContent: string
    onCancel: () => void
    onUpdate: (newContent: string) => void
}

// Update comment compoenent that renders as a pop up in the parent 
const UpdateComment = ({ commentId, initialContent, onCancel, onUpdate }: UpdateCommentProps) => {
    const [content, setContent] = useState(initialContent)
    const [isUpdating, setIsUpdating] = useState(false)

    // Function that handles any new content changes. 
    const handleUpdate = async () => {
        if (!content.trim()) {
            toast.error("Comment content cannot be empty.")
            return
        }

        try {
            setIsUpdating(true)
            await updateComment({commentId, content})
            toast.success("Comment updated successfully")
            onUpdate(content) // Update parent state so that it instantly reflects change and no refresh requied 
        } catch (err) {
            console.error("Failed to update comment:", err)
            toast.error("Failed to update comment")
        } finally {
            setIsUpdating(false)
        }
    }

    return (
        <Box mt={2} display="flex" flexDirection="column" gap={1}>
            <RichTextEditor content={content} onChange={setContent} />

            <Box display="flex" gap={1} justifyContent="flex-end">
                <Button variant="outlined" color="secondary" onClick={onCancel} disabled={isUpdating}>
                    Cancel
                </Button>
                <Button variant="contained" color="primary" onClick={handleUpdate} disabled={isUpdating}>
                    {isUpdating ? <CircularProgress size={20} /> : "Update"}
                </Button>
            </Box>
        </Box>
    )
}

export default UpdateComment
