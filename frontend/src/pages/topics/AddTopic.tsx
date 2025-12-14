// pages/AddTopic.tsx
import { useState } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import { createTopic } from "../../api/handleTopic";

const AddTopic = () => {
    const [name, setName] = useState("");
    const [message, setMessage] = useState("");

    const handleSubmit = async () => {
        if (!name.trim()) {
            setMessage("Topic name cannot be empty");
            return;
        }

        try {
            const topic = await createTopic(name);
            setMessage(`Topic "${topic.Name}" created successfully!`);
            setName(""); // reset input
        } catch (err: any) {
            setMessage(err.message || "Failed to create topic");
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
            {message && (
                <Typography color={message.includes("successfully") ? "primary" : "error"}>
                    {message}
                </Typography>
            )}
        </Box>
    );
};

export default AddTopic;
