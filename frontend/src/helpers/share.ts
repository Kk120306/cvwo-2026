// Share functionality for posts 
function sharePost(postId: string, title: string) {
    const url = `${window.location.origin}/posts/${postId}`;

    // If the web share api is avaliable in the browsser use it to send to different points
    if (navigator.share) {
        navigator.share({
            title,
            text: "Check out this post",
            url,
        }).catch(() => {
            // ignore
        });
    } else {
        // If not it will just copy the url to the clipboard 
        navigator.clipboard.writeText(url);
        alert("Link copied to clipboard");
    }
}


export { sharePost };