import { Box, Chip } from "@mui/material";
import { Link } from "react-router-dom";

// Props for topic bar. 
interface Props {
    topics: string[];
    selectedTopic: string;
    onSelect: (topic: string) => void;
}

// Topic Bar component 
const TopicBar = ({ topics, selectedTopic, onSelect }: Props) => {
    return (
        <Box display="flex" gap={1} flexWrap="wrap">
            {topics.map((topic) => (
                <Chip
                    key={topic}
                    label={topic}
                    sx={{ minWidth: 60 }}
                    clickable
                    color={selectedTopic.toLowerCase() === topic.toLowerCase() ? "primary" : "default"}
                    onClick={() => onSelect(topic)}
                />
            ))}
            {/* Allows user to create a topic */}
            <Chip
                label="+"
                clickable
                sx={{
                    '&:hover': {
                        cursor: 'pointer',
                        transform: 'scale(1.1)'
                    },
                    transition: 'transform 0.2s'
                }}
                color="secondary"
                component={Link}
                to="/topics/create"
            />
        </Box>
    );
};

export default TopicBar;