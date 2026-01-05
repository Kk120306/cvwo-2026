import { useState, useEffect } from 'react'
import { Box, Button, TextField, Typography, MenuItem, Select, FormControl, InputLabel } from '@mui/material'
import RichTextEditor from '../../components/provider/RichTextEditor'
import { createPost } from '../../api/handlePost'
import { useNavigate } from 'react-router-dom'
import type { Topic } from '../../types/globalTypes'
import { fetchAllTopics } from "../../api/handleTopic"
import { useAppSelector } from '../../hooks/reduxHooks'
import { toast } from "react-hot-toast";
import ImageForm from '../../components/image/ImageForm'
import { getS3Url, uploadFileToS3 } from '../../api/handleImage'

// Page for users to create a new post
const CreatePostPage = () => {
    const user = useAppSelector(state => state.auth.user);
    const [title, setTitle] = useState('')
    const [content, setContent] = useState('')
    const [topics, setTopics] = useState<Topic[]>([])
    const [selectedTopic, setSelectedTopic] = useState<string>('')
    const [imageFile, setImageFile] = useState<File | null>(null)
    const [imagePreview, setImagePreview] = useState<string | null>(null)

    const navigate = useNavigate()

    // If the user is not found in state, redirects them to signup
    useEffect(() => {
        if (!user) navigate('/signup');
    }, [user, navigate]);

    if (!user) return <p>Loading...</p>;

    // Loads all topics that exist currently
    useEffect(() => {
        const loadTopics = async () => {
            try {
                const data: Topic[] = await fetchAllTopics()
                setTopics([...data])
            } catch {
                console.error('Failed to load topics')
            }
        }

        loadTopics()
    }, [])

    // handles the submission of the post
    const handleSubmit = async () => {
        if (!title || !content || !selectedTopic) {
            toast.error('Please fill in all fields');
            return;
        }

        let imageUrl = null;
        // If there is an image we upload the image and then save the uploaded Url as the image Url 
        if (imageFile) {
            const signedUrl = await getS3Url()
            const uploadedUrl = await uploadFileToS3(signedUrl, imageFile)
            if (uploadedUrl) {
                imageUrl = uploadedUrl
            }
        }

        await createPost({ title, content, topicSlug: selectedTopic, imageUrl });
        navigate(-1)
    }

    return (
        <Box maxWidth={600} mx="auto" mt={4}>
            <Typography variant="h4" mb={2}>Create Post</Typography>

            <TextField
                label="Title"
                fullWidth
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                margin="normal"
            />

            {/* Topic dropdown */}
            <FormControl fullWidth margin="normal">
                <InputLabel id="topic-label">Topic</InputLabel>
                <Select
                    labelId="topic-label"
                    value={selectedTopic}
                    label="Topic"
                    onChange={(e) => setSelectedTopic(e.target.value)}
                >
                    {topics.map((topic) => (
                        <MenuItem key={topic.id} value={topic.slug}>
                            {topic.name}
                        </MenuItem>
                    ))}
                    {/* Allows users to create a new topic if not there */}
                    <MenuItem
                        onClick={() => navigate('/topics/create')}
                        sx={{ fontStyle: 'italic', backgroundColor: '#f0f0f5' }}
                    >
                        + Create New Topic
                    </MenuItem>
                </Select>
            </FormControl>

            <RichTextEditor content={content} onChange={setContent} />

            {/* Component that handles image uploads and preview */}
            <ImageForm
                imageFile={imageFile}
                setImageFile={setImageFile}
                imagePreview={imagePreview}
                setImagePreview={setImagePreview}
            />

            <Button variant="contained" color="primary" onClick={handleSubmit} sx={{ mt: 2 }}>
                Submit
            </Button>
        </Box>
    )
}

export default CreatePostPage
