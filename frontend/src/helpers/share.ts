function sharePost(postId: string, title: string) {
    const url = `${window.location.origin}/posts/${postId}`;

    if (navigator.share) {
        navigator.share({
            title,
            text: "Check out this post",
            url,
        }).catch(() => { });
    } else if (navigator.clipboard?.writeText) {
        navigator.clipboard.writeText(url)
            .then(() => alert("Link copied to clipboard"))
            .catch(() => prompt("Copy this link:", url));
    } else {
        prompt("Copy this link:", url);
    }
}

export { sharePost };