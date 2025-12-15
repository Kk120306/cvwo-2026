import { useState, useEffect } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import { createTopic } from "../../api/handleTopic";
import { formatTopicName } from "../../helpers/formatter";
import { useNavigate } from "react-router-dom";
import { useAppSelector } from "../../hooks/reduxHooks";
import { toast } from "react-hot-toast";

// Page to create a new topic
const AddTopic = () => {
    const user = useAppSelector(state => state.auth.user);
    const [name, setName] = useState("");
    const navigate = useNavigate();

    // If the user is not found in state, redirects them to signup
    useEffect(() => {
        if (!user) navigate('/signup');
    }, [user, navigate]);

    if (!user) return <p>Loading...</p>;

    // function handles the submission of a new topic 
    const handleSubmit = async () => {
        // ensures name is not empty
        if (!name.trim()) {
            toast.error("Topic name cannot be empty");
            return;
        }

        // Creates a topic via backend 
        try {
            await createTopic(formatTopicName(name));
            setName("");
            navigate(-1);
        } catch (err: any) {
            console.log(err.message || "Failed to create topic");
        }
    };

    return (
        <Box maxWidth="400px" mx="auto" mt={6} display="flex" flexDirection="column" gap={2}>
            <Typography variant="h5">Add New Topic</Typography>
            <TextField
                label="Topic Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                fullWidth
            />
            <Button variant="contained" onClick={handleSubmit}>
                Create Topic
            </Button>
        </Box>
    );
};

export default AddTopic;
