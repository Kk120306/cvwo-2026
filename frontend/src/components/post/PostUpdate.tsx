import { useState } from "react"
import { Box, Button, CircularProgress, TextField } from "@mui/material"
import { updatePost } from "../../api/handlePost"
import RichTextEditor from "../RichTextEditor"
import type { Post } from "../../types/globalTypes"
import { normalizePost } from "../../helpers/normalizer"
import ImageForm from "../image/ImageForm"
import { getS3Url, uploadFileToS3, deleteImage } from "../../api/handleImage"

interface UpdatePostProps {
    postId: string
    initialTitle: string
    initialContent: string
    initialImage: string | null
    onCancel?: () => void
    newPost: (updatedPost: Post) => void
}

const UpdatePost = ({ postId, initialTitle, initialContent, initialImage, onCancel, newPost }: UpdatePostProps) => {
    const [content, setContent] = useState(initialContent)
    const [title, setTitle] = useState(initialTitle)
    const [isUpdating, setIsUpdating] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [imageFile, setImageFile] = useState<File | null>(null)
    const [imagePreview, setImagePreview] = useState<string | null>(initialImage)

    // Function to change the post to new values 
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

            let imageUrl: string | null = initialImage

            // User uploaded a new image
            if (imageFile) {
                const signedUrl = await getS3Url()
                const uploadedUrl = await uploadFileToS3(signedUrl, imageFile)
                if (uploadedUrl) {
                    imageUrl = uploadedUrl

                    // Delete old image since we have a new image replacing it 
                    if (initialImage) {
                        try {
                            await deleteImage(initialImage)
                        } catch (err) {
                            console.error("Failed to delete previous image:", err)
                        }
                    }
                }
            }
            // User removed an existing image and did not add a new one 
            else if (!imagePreview && initialImage) {
                imageUrl = null
                try {
                    await deleteImage(initialImage)
                } catch (err) {
                    console.error("Failed to delete image:", err)
                }
            }
            // No changes to image - keep initialImage
            const res = await updatePost({ postId, title, content, imageUrl })
            const updatedPost = normalizePost(res)
            newPost(updatedPost)
            onCancel?.()
        } catch (err) {
            console.error("Failed to update post:", err)
            setError("Failed to update post. Try again.")
        } finally {
            setIsUpdating(false)
        }
    }

    // Function to handle cancel button click
    const handleCancel = () => {
        // If we created a in browser preview URL, no longer need it so revoke it 
        // This is done in the Image Form so that we dont end up uploading to bucket before user sends final submit 
        if (imagePreview && imagePreview.startsWith('blob:')) {
            URL.revokeObjectURL(imagePreview)
        }
        onCancel?.()
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

            {/* Componenet form for user to update new images or add image to existing post  */}
            <ImageForm
                imageFile={imageFile}
                setImageFile={setImageFile}
                imagePreview={imagePreview}
                setImagePreview={setImagePreview}
            />

            <Box display="flex" gap={1} justifyContent="flex-end">
                {onCancel && (
                    <Button
                        variant="outlined"
                        color="secondary"
                        onClick={handleCancel}
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