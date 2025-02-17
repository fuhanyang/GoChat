
function getHeaders() {
    const jwt = localStorage.getItem('jwt');
    if (jwt) {
        return {
            "Content-Type": "application/json",
            'Authorization': `Bearer ${jwt}`
        };
    }else {
        alert('You are not logged in. Please log in to access this page.');
        return null;
    }
}