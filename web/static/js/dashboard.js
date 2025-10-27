document.addEventListener('DOMContentLoaded', async function() {
    const userInfo = document.getElementById('user-info');
    const logoutBtn = document.getElementById('logout-btn');
    
    // Fetch current user info
    try {
        const response = await fetch('/api/me');
        
        if (response.ok) {
            const user = await response.json();
            
            userInfo.innerHTML = `
                <h2>User Information</h2>
                <p><strong>Name:</strong> ${user.first_name} ${user.last_name}</p>
                <p><strong>Email:</strong> ${user.email}</p>
                <p><strong>Role:</strong> ${user.role}</p>
                <p><strong>Status:</strong> ${user.status}</p>
                <p><strong>Account Created:</strong> ${new Date(user.created_at).toLocaleDateString()}</p>
            `;
        } else {
            // Not authenticated, redirect to login
            window.location.href = '/login';
        }
    } catch (error) {
        // Error fetching user info, redirect to login
        window.location.href = '/login';
    }
    
    // Handle logout
    logoutBtn.addEventListener('click', async function() {
        try {
            const response = await fetch('/api/logout', {
                method: 'POST'
            });
            
            if (response.ok) {
                window.location.href = '/login';
            }
        } catch (error) {
            console.error('Logout failed:', error);
        }
    });
});
