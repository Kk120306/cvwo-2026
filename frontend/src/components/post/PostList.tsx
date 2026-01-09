import { useEffect, useState, useMemo } from "react"
import {
    Box,
    Typography,
    TextField,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
} from "@mui/material"
import { fetchPostByTopic, managePostPin } from "../../api/handlePost"
import type { Post } from "../../types/globalTypes"
import { votePost } from "../../api/handleVote"
import { useAppSelector } from "../../hooks/reduxHooks";
import PostCard from "../post/PostCard"
import { useNavigate } from "react-router-dom";

// Props for PostList component
interface PostListProps {
    topic: string
}

// The ways the post can be sorted by
type SortOption = "recent" | "likes" | "dislikes";

// componenet that renders the posts that exist on the forum
export default function PostList({ topic }: PostListProps) {
    const [posts, setPosts] = useState<Post[]>([])
    const [search, setSearch] = useState("")
    const [sortBy, setSortBy] = useState<SortOption>("recent")
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")

    const user = useAppSelector(state => state.auth.user)
    const navigate = useNavigate();

    // On topic change, goes and retrives post under topic 
    useEffect(() => {
        const loadPosts = async () => {
            try {
                setLoading(true)
                setError("")
                const data = await fetchPostByTopic(topic)
                console.log(data)
                setPosts(data || [])
            } catch {
                setError("Failed to load posts")
                setPosts([])
            } finally {
                setLoading(false)
            }
        }

        loadPosts()
    }, [topic])

    // Client side search filtering and sorting
    // Only computed when dependencies change not on every render
    const filteredAndSortedPosts = useMemo(() => {
        const q = search.toLowerCase() // query

        // Filter posts
        const filtered = posts.filter((post) =>
            post.title.toLowerCase().includes(q) ||
            post.content.toLowerCase().includes(q) ||
            post.author.username.toLowerCase().includes(q)
        )

        // Sort posts - use spread because we are using useMemo and don't want to mutate original array
        const sorted = [...filtered].sort((a, b) => {
            // Pinned posts always come first
            if (a.isPinned && !b.isPinned) return -1
            if (!a.isPinned && b.isPinned) return 1

            // Then sort by selected option
            switch (sortBy) {
                case "likes":
                    return b.likes - a.likes
                case "dislikes":
                    return b.dislikes - a.dislikes
                case "recent":
                default:
                    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
            }
        })

        return sorted
    }, [posts, search, sortBy])

    // Hnadle voting on a post 
    const handleVote = async (postId: string, type: "like" | "dislike") => {
        if (!user) {
            navigate("/login")
        }

        const res = await votePost(postId, type) // Call API to vote 


        // Takes most recent post state and updates it to new post state 
        setPosts(prev =>
            // take a post obj and if the modified postID matches the post, update its likes/dislikes/myVote
            prev.map(post =>
                post.id === postId
                    ? {
                        ...post,
                        likes: res.likes,
                        dislikes: res.dislikes,
                        myVote: res.myVote,
                    }
                    : post
            )
        )
    }

    // Function that handles the pinning and unpinning of a post 
    const handlePin = async (postId: string, pin: boolean) => {
        try {
            await managePostPin(postId, pin)
            setPosts(prev =>
                prev.map(post =>
                    post.id === postId
                        ? { ...post, isPinned: pin }
                        : post
                )
            )
        } catch {
            alert("Failed to update pin status")
            return
        }
    }


    if (loading) return <Typography>Loading posts...</Typography>
    if (error) return <Typography color="error">{error}</Typography>

    return (
        <Box display="flex" flexDirection="column" gap={2}>
            {/* Search and Sort Controls */}
            <Box display="flex" gap={2}>
                <TextField
                    placeholder="Search posts..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    fullWidth
                />
                <FormControl sx={{ minWidth: 150 }}>
                    <InputLabel>Sort By</InputLabel>
                    <Select
                        value={sortBy}
                        label="Sort By"
                        onChange={(e) => setSortBy(e.target.value as SortOption)}
                    >
                        <MenuItem value="recent">Most Recent</MenuItem>
                        <MenuItem value="likes">Most Liked</MenuItem>
                        <MenuItem value="dislikes">Most Disliked</MenuItem>
                    </Select>
                </FormControl>
            </Box>

            {filteredAndSortedPosts.length === 0 ? (
                <Typography>No posts found</Typography>
            ) : (
                filteredAndSortedPosts.map((post) => (
                    <PostCard
                        key={post.id}
                        post={post}
                        isAdmin={user?.isAdmin}
                        onVote={handleVote}
                        onPin={user?.isAdmin ? handlePin : undefined}
                    />
                ))
            )}
        </Box>
    )
}