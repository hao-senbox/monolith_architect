$(document).ready(function() {
    const BASE_URL = 'https://monolith-architect.onrender.com';
    const API_ENDPOINTS = {
        login: `${BASE_URL}/api/v1/user/login`,
        register: `${BASE_URL}/api/v1/user/register`,
        refreshToken: `${BASE_URL}/api/v1/user/refresh`
    };

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

    // Refresh Token Function
    function refreshToken() {
        return new Promise((resolve, reject) => {
            const refreshToken = getRefreshToken();
            
            if (!refreshToken) {
                reject(new Error('No refresh token available'));
                return;
            }

            $.ajax({
                url: API_ENDPOINTS.refreshToken,
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${refreshToken}`
                },
                success: function(response) {
                    if (response.data && response.data.token && response.data.refreshToken) {
                        setTokens(response.data.token, response.data.refreshToken);
                        resolve(response.data.token);
                    } else {
                        reject(new Error('Invalid token response'));
                    }
                },
                error: function(xhr) {
                    clearTokens();
                    reject(xhr);
                }
            });
        });
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

    // Toggle password visibility
    $('.toggle-password').click(function() {
        const passwordInput = $(this).closest('.relative').find('input');
        const type = passwordInput.attr('type') === 'password' ? 'text' : 'password';
        passwordInput.attr('type', type);
        $(this).toggleClass('fa-eye fa-eye-slash');
    });

    // Login Form Handler
    if ($('#loginForm').length) {
        $('#loginForm').submit(function(e) {
            e.preventDefault();
            
            hideError('emailError');
            hideError('passwordError');
            
            const email = $('#email').val().trim();
            const password = $('#password').val().trim();
            
            if (!email) {
                showError('emailError', 'Email is required');
                return;
            }
            
            if (!validateEmail(email)) {
                showError('emailError', 'Please enter a valid email address');
                return;
            }
            
            if (!password) {
                showError('passwordError', 'Password is required');
                return;
            }
            
            showLoading();
            
            $.ajax({
                url: API_ENDPOINTS.login,
                method: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({
                    email: email,
                    password: password
                }),
                success: function(response) {
                    setTokens(response.data.token, response.data.refreshToken);
                    
                    showToast('Login successful! Redirecting...', 'success');
                    
                    if (response.data.user_type === 'admin') {
                        window.location.href = '/frontend/admin/categories.html';
                    } else {
                        window.location.href = '/frontend/admin/categories.html';
                    }

                },
                error: function(xhr) {
                    hideLoading();
                    
                    if (xhr.status === 401) {
                        showError('passwordError', 'Invalid email or password');
                    } else if (xhr.status === 403) {
                        showToast('Your account has been locked. Please contact support.', 'error');
                    } else {
                        showToast('An error occurred. Please try again later.', 'error');
                    }
                }
            });
        });

        // Clear error messages on input
        $('#email, #password').on('input', function() {
            const fieldId = $(this).attr('id');
            hideError(`${fieldId}Error`);
        });
    }

    // Registration Form Handler
    if ($('#registerForm').length) {
        $('#registerForm').submit(function(e) {
            e.preventDefault();
            
            // Reset all error messages
            $('[id$="Error"]').addClass('hidden');
            
            // Get form values
            const firstName = $('#firstName').val().trim();
            const lastName = $('#lastName').val().trim();
            const email = $('#email').val().trim();
            const phone = $('#phone').val().trim();
            const password = $('#password').val().trim();
            const confirmPassword = $('#confirmPassword').val().trim();
            const terms = $('#terms').is(':checked');
            
            // Validate all fields
            let isValid = true;
            
            if (!firstName) {
                showError('firstNameError', 'First name is required');
                isValid = false;
            }
            
            if (!lastName) {
                showError('lastNameError', 'Last name is required');
                isValid = false;
            }
            
            if (!email) {
                showError('emailError', 'Email is required');
                isValid = false;
            } else if (!validateEmail(email)) {
                showError('emailError', 'Please enter a valid email address');
                isValid = false;
            }
            
            if (!phone) {
                showError('phoneError', 'Phone number is required');
                isValid = false;
            } else if (!validatePhone(phone)) {
                showError('phoneError', 'Please enter a valid phone number');
                isValid = false;
            }
            
            if (!password) {
                showError('passwordError', 'Password is required');
                isValid = false;
            } else if (password.length < 6) {
                showError('passwordError', 'Password must be at least 6 characters');
                isValid = false;
            }
            
            if (!confirmPassword) {
                showError('confirmPasswordError', 'Please confirm your password');
                isValid = false;
            } else if (password !== confirmPassword) {
                showError('confirmPasswordError', 'Passwords do not match');
                isValid = false;
            }
            
            if (!terms) {
                showError('termsError', 'You must agree to the terms and conditions');
                isValid = false;
            }
            
            if (!isValid) return;
            
            showLoading();
            
            // Make API call
            $.ajax({
                url: API_ENDPOINTS.register,
                method: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({
                    fristName: firstName,
                    lastName: lastName,
                    email: email,
                    phone: phone,
                    password: password
                }),
                success: function(response) {
                    showToast('Registration successful! Redirecting to login...', 'success');
                    
                    setTimeout(() => {
                        window.location.href = '/frontend/pages/auth/login.html';
                    }, 2000);
                },
                error: function(xhr) {
                    hideLoading();
                    
                    if (xhr.status === 409) {
                        showError('emailError', 'Email already exists');
                    } else {
                        showToast('An error occurred. Please try again later.', 'error');
                    }
                }
            });
        });

        // Clear error messages on input
        $('#firstName, #lastName, #email, #phone, #password, #confirmPassword').on('input', function() {
            const fieldId = $(this).attr('id');
            hideError(`${fieldId}Error`);
        });

        $('#terms').on('change', function() {
            hideError('termsError');
        });
    }
}); 