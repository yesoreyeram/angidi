document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('reset-password-form');
    const errorMessage = document.getElementById('error-message');
    const successMessage = document.getElementById('success-message');
    const tokenInput = document.getElementById('token');
    
    // Get token from URL query parameter
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    
    if (!token) {
        errorMessage.textContent = 'Invalid reset link. Please request a new password reset.';
        errorMessage.style.display = 'block';
        form.style.display = 'none';
        return;
    }
    
    tokenInput.value = token;

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Clear previous messages
        errorMessage.style.display = 'none';
        successMessage.style.display = 'none';
        
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        // Validate passwords match
        if (newPassword !== confirmPassword) {
            errorMessage.textContent = 'Passwords do not match';
            errorMessage.style.display = 'block';
            return;
        }
        
        try {
            const response = await fetch('/api/reset-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    token: token,
                    new_password: newPassword
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                // Password reset successful
                successMessage.textContent = 'Password reset successful! Redirecting to login...';
                successMessage.style.display = 'block';
                
                // Redirect to login after 2 seconds
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                // Show error message
                errorMessage.textContent = data.message || 'Password reset failed';
                errorMessage.style.display = 'block';
            }
        } catch (error) {
            errorMessage.textContent = 'An error occurred. Please try again.';
            errorMessage.style.display = 'block';
        }
    });
});
