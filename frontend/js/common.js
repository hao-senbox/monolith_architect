// Token Management
function getToken() {
    return localStorage.getItem('token');
}

function getRefreshToken() {
    return localStorage.getItem('refreshToken');
}

function setTokens(token, refreshToken) {
    localStorage.setItem('token', token);
    localStorage.setItem('refreshToken', refreshToken);
}

function clearTokens() {
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
}

// Common Functions
function showLoading() {
    $('#loadingSpinner').removeClass('hidden').addClass('flex');
}

function hideLoading() {
    $('#loadingSpinner').removeClass('flex').addClass('hidden');
}

function showToast(message, type = 'info') {
    const toast = $(`
        <div class="toast toast-${type}">
            <div class="flex items-center">
                <i class="fas ${type === 'success' ? 'fa-check-circle' : type === 'error' ? 'fa-exclamation-circle' : 'fa-info-circle'} mr-2"></i>
                <span>${message}</span>
            </div>
        </div>
    `);
    
    $('#toastContainer').append(toast);
    
    setTimeout(() => {
        toast.fadeOut(300, function() {
            $(this).remove();
        });
    }, 3000);
}

function showError(elementId, message) {
    $(`#${elementId}`).text(message).removeClass('hidden');
}

function hideError(elementId) {
    $(`#${elementId}`).addClass('hidden');
}

function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

function validatePhone(phone) {
    const re = /^\+?[\d\s-]{10,}$/;
    return re.test(phone);
}

// Setup AJAX Interceptor
$.ajaxSetup({
    beforeSend: function(xhr) {
        const token = getToken();
        if (token) {
            xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        }
    },
    error: function(xhr, status, error) {
        if (xhr.status === 401) {
            // Token expired, try to refresh
            refreshToken()
                .then(newToken => {
                    // Retry the original request with new token
                    const originalRequest = this;
                    originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
                    return $.ajax(originalRequest);
                })
                .catch(error => {
                    // Refresh failed, redirect to login
                    showToast('Session expired. Please login again.', 'error');
                    setTimeout(() => {
                        window.location.href = '/frontend/pages/auth/login.html';
                    }, 1500);
                });
        }
    }
}); 