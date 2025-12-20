import { Box, Button, IconButton } from "@mui/material"
import CloseIcon from '@mui/icons-material/Close'

// Prop for Image form 
interface ImageFormProps {
    imageFile: File | null
    setImageFile: (file: File | null) => void
    imagePreview: string | null
    setImagePreview: (url: string | null) => void
}

// Component form that handles image upload and previw
const ImageForm = ({
    imageFile,
    setImageFile,
    imagePreview,
    setImagePreview,
}: ImageFormProps) => {

    // Function that handles when a user clicks a new file 
    const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0]

        if (!file) {
            return
        }

        // Clean up old blob preview if exists - ensures that we dont leak memory 
        if (imagePreview && imagePreview.startsWith('blob:')) {
            URL.revokeObjectURL(imagePreview)
        }

        // Store the File object (not uploaded yet)
        setImageFile(file)

        // Create local preview URL - this makes sure that we dont have to upload to s3 bucket first 
        const previewUrl = URL.createObjectURL(file)
        setImagePreview(previewUrl)
    }

    // Function that handles removal of selected image 
    const handleRemoveImage = () => {
        // Clean up the preview URL to avoid memory leaks - if the image that was shown was saved to browser memory 
        if (imagePreview && imagePreview.startsWith('blob:')) {
            URL.revokeObjectURL(imagePreview)
        }
        // Clear both the file and preview since there is only one media that can be uploaded 
        setImageFile(null)
        setImagePreview(null)
    }

    return (
        <Box mt={2}>
            <Button variant="outlined" component="label">
                {imageFile || imagePreview ? "Change Image" : "Upload Image"}
                <input
                    type="file"
                    accept="image/*"
                    hidden
                    onChange={handleFileSelect}
                />
            </Button>

            {imagePreview && (
                <Box mt={2} position="relative" display="inline-block">
                    <img
                        src={imagePreview}
                        alt="Preview"
                        style={{ maxWidth: '200px', maxHeight: '200px', borderRadius: '8px' }}
                    />
                    <IconButton
                        onClick={handleRemoveImage}
                        size="small"
                        sx={{
                            position: 'absolute',
                            top: 4,
                            right: 4,
                            backgroundColor: 'rgba(0,0,0,0.6)',
                            '&:hover': { backgroundColor: 'rgba(0,0,0,0.8)' }
                        }}
                    >
                        <CloseIcon fontSize="small" sx={{ color: 'white' }} />
                    </IconButton>
                </Box>
            )}
        </Box>
    )
}

export default ImageForm