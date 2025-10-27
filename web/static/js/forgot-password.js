document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('forgot-password-form');
    const errorMessage = document.getElementById('error-message');
    const successMessage = document.getElementById('success-message');

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Clear previous messages
        errorMessage.style.display = 'none';
        successMessage.style.display = 'none';
        
        const email = document.getElementById('email').value;
        
        try {
            const response = await fetch('/api/forgot-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                // Show success message
                successMessage.textContent = data.message || 'If an account exists with this email, password reset instructions have been sent';
                successMessage.style.display = 'block';
                form.reset();
            } else {
                // Show error message
                errorMessage.textContent = data.message || 'Failed to process request';
                errorMessage.style.display = 'block';
            }
        } catch (error) {
            errorMessage.textContent = 'An error occurred. Please try again.';
            errorMessage.style.display = 'block';
        }
    });
});
