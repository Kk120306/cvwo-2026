import { useEffect, useState } from "react"
import { Box, Container, Typography } from "@mui/material"
import TopicBar from "../components/dashboard/TopicBar"
import PostList from "../components/dashboard/PostList"
import { fetchAllTopics } from "../api/handleTopic"
import type { Topic } from "../types/globalTypes"


// TODO : abstract Load topics  - possibly add redux for topics instead
// The main page - displays topics, posts 
const DiscoverPage = () => {
    const [topics, setTopics] = useState<Topic[]>([]);
    const [selectedTopic, setSelectedTopic] = useState<string>("all");
    const [error, setError] = useState("");

    // On mount, fetch all topics and also adds a "All" topic
    useEffect(() => {
        const loadTopics = async () => {
            try {
                const data: Topic[] = await fetchAllTopics();
                setTopics([{ id: "all", name: "All", slug: "all" }, ...data]);
            } catch {
                setError("Failed to load topics");
            }
        };
        loadTopics();
    }, []);

    return (
        <Container maxWidth="md">
            <Typography variant="h4" fontWeight={600} mt={4} mb={2}>
                Discover
            </Typography>

            {/* If there is a error fetching topics */}
            {error && (
                <Typography color="error" mb={2}>
                    {error}
                </Typography>
            )}

            {/* Componenet that maps the topics */}
            <TopicBar
                topics={topics.map(t => t.name)}
                selectedTopic={selectedTopic}
                onSelect={(name) => {
                    const topic = topics.find(t => t.name === name);
                    if (topic) setSelectedTopic(topic.slug);
                }} // Logic to set selected topic based on name in state
            />

            {/* Compoenent that renders posts based on selected topic  */}
            <Box mt={4}>
                <PostList topic={selectedTopic} />
            </Box>
        </Container>
    );
};

export default DiscoverPage;