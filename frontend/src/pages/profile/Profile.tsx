import { Container, Typography, Box, Avatar, Card, CardContent, CircularProgress, Alert, Chip, Stack, Divider } from '@mui/material';
import { useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { getUserProfile } from '../../api/handleUser';
import type { UserProfile } from '../../types/globalTypes';
import PostCard from '../../components/post/PostCard';
import CommentList from '../../components/comments/CommentList';

const ProfilePage = () => {
    const { username } = useParams<{ username: string }>();

    const [user, setUser] = useState<UserProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        const loadUserData = async () => {
            if (!username) {
                setError("Username is required");
                setLoading(false);
                return;
            }

            try {
                setLoading(true);
                setError("");

                const data = await getUserProfile(username, true, true);
                setUser(data);
                console.log(data.avatarUrl);
            } catch (err: unknown) {
                if (err instanceof Error) {
                    setError(err.message); 
                } else {
                    setError("Failed to load user profile"); 
                }
            } finally {
                setLoading(false);
            }
        };

        loadUserData();
    }, [username]);

    if (loading) {
        return (
            <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
            </Container>
        );
    }

    if (error) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{error}</Alert>
            </Container>
        );
    }

    if (!user) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="info">User not found</Alert>
            </Container>
        );
    }

    return (
        <Container sx={{ mt: 4, mb: 4 }}>
            <Card>
                <CardContent>
                    {/* Header Section */}
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, mb: 3 }}>
                        <Avatar
                            src={user.avatarUrl}
                            alt={user.username}
                            sx={{ width: 80, height: 80 }}
                        >
                        </Avatar>

                        <Box sx={{ flex: 1 }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                                <Typography variant="h4" component="h1">
                                    {user.username}
                                </Typography>
                                {user.isAdmin && (
                                    <Chip label="Admin" color="primary" size="small" />
                                )}
                            </Box>
                        </Box>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    {/* Stats Section */}
                    <Stack direction="row" spacing={4} sx={{ mt: 2 }}>
                        <Box>
                            <Typography variant="h6" color="primary">
                                {user.postCount || 0}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Posts
                            </Typography>
                        </Box>
                        <Box>
                            <Typography variant="h6" color="primary">
                                {user.commentCount || 0}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Comments
                            </Typography>
                        </Box>
                        <Box>
                            <Typography variant="body2" color="text.secondary">
                                Joined {new Date(user.createdAt).toLocaleDateString()}
                            </Typography>
                        </Box>
                    </Stack>
                </CardContent>
            </Card>
            <Box sx={{ mt: 4 }}>
                <Typography variant="h5" gutterBottom>
                    Posts by {user.username}
                </Typography>
                {user.posts && user.posts.length > 0 ? (
                    console.log(user.posts),
                    user.posts.map((post) => (
                        <PostCard
                            key={post.id}
                            post={post}
                            isAdmin={false} // Cant have piknning feature on profile page, dosent make sense 
                        />
                    ))
                ) : (
                    <Typography>No posts to display.</Typography>
                )}

                <Box sx={{ mt: 4 }}>
                    <Typography variant="h5" gutterBottom>
                        Comments by {user.username}
                    </Typography>
                    {user.comments && user.comments.length > 0 ? (
                        console.log(user.comments),
                        <CommentList
                            comments={user.comments}
                        />
                    ) : (
                        <Typography>No Comments to display.</Typography>
                    )}
                </Box>
            </Box>
        </Container>
    );
};

export default ProfilePage;
