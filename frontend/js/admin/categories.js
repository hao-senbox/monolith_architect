$(document).ready(function() {
    const BASE_URL = 'https://monolith-architect.onrender.com';
    const API_ENDPOINTS = {
        categories: `${BASE_URL}/api/v1/category`,
        category: (id) => `${BASE_URL}/api/v1/category/${id}`
    };

    let categories = [];
    let currentCategoryId = null;

    // Check if user is logged in
    if (!getToken()) {
        window.location.href = '/frontend/pages/auth/login.html';
        return;
    }

    // Load layout
    $('#layout').load('/frontend/components/layout.html', function() {
        setActiveNavItem('/admin/categories.html');
        setPageTitle('Category Management');
        loadCategories();
    });

    // Load categories
    function loadCategories() {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.categories,
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
            success: function(response) {
                categories = response.data || [];
                renderCategoryTree();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to load categories', 'error');
                $('#categoryTree').html('<div class="error-state">Failed to load categories</div>');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Render category tree
    function renderCategoryTree(filteredCategories = null) {
        const categoriesToRender = filteredCategories || categories || [];
        const rootCategories = categoriesToRender.filter(cat => !cat.parent_id);
        
        let html = '';
        rootCategories.forEach(category => {
            html += createCategoryElement(category, categoriesToRender, new Set());
        });

        if (html === '') {
            html = '<div class="empty-state">No categories found</div>';
        }

        $('#categoryTree').html(html);
        setupEventListeners();
    }

    // Create category element with circular reference protection
    function createCategoryElement(category, allCategories, visitedIds = new Set(), level = 0) {
        // Prevent circular references
        if (visitedIds.has(category.id)) {
            console.warn(`Circular reference detected for category: ${category.category_name} (ID: ${category.id})`);
            return `<div class="category-item error" data-id="${category.id}">
                <div class="category-content">
                    <i class="fas fa-exclamation-triangle category-icon text-red-500"></i>
                    <span class="category-name text-red-500">${category.category_name} (Circular Reference)</span>
                </div>
            </div>`;
        }

        // Add current category to visited set
        const newVisitedIds = new Set(visitedIds);
        newVisitedIds.add(category.id);

        // Get children, excluding any that would create circular references
        const children = allCategories.filter(cat => {
            return cat.parent_id === category.id && !visitedIds.has(cat.id);
        });
        
        const hasChildren = children.length > 0;
        
        let html = `
            <div class="category-item" data-id="${category.id}" style="margin-left: ${level * 20}px;">
                <div class="category-content">
                    <i class="fas ${hasChildren ? 'fa-folder' : 'fa-folder-open'} category-icon"></i>
                    <span class="category-name">${category.category_name}</span>
                    ${level > 0 ? `<span class="text-xs text-gray-500 ml-2">(Level ${level})</span>` : ''}
                </div>
                <div class="category-actions">
                    <button class="category-action-btn edit" title="Edit">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="category-action-btn delete" title="Delete">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `;

        // Recursively render children with protection against infinite loops
        if (hasChildren && level < 10) { // Limit nesting depth as additional protection
            html += '<div class="category-children">';
            children.forEach(child => {
                html += createCategoryElement(child, allCategories, newVisitedIds, level + 1);
            });
            html += '</div>';
        } else if (level >= 10) {
            html += '<div class="category-children"><div class="text-xs text-yellow-600 ml-4">Maximum nesting depth reached</div></div>';
        }

        return html;
    }

    // Validate category hierarchy before saving
    function validateCategoryHierarchy(categoryId, parentId, allCategories) {
        if (!parentId) return true; // No parent is always valid
        
        const visited = new Set();
        let currentId = parentId;
        
        while (currentId) {
            if (currentId === categoryId) {
                return false; // Would create circular reference
            }
            
            if (visited.has(currentId)) {
                return false; // Circular reference exists in parent chain
            }
            
            visited.add(currentId);
            const parent = allCategories.find(cat => cat.id === currentId);
            currentId = parent ? parent.parent_id : null;
        }
        
        return true;
    }

    // Setup event listeners
    function setupEventListeners() {
        // Add category button
        $('#addCategoryBtn').off('click').on('click', function() {
            currentCategoryId = null;
            $('#modalTitle').text('Add Category');
            $('#categoryForm')[0].reset();
            populateParentSelect();
            $('#categoryModal').removeClass('hidden').addClass('flex');
        });

        // Close modal
        $('#closeModal, #cancelBtn').off('click').on('click', function() {
            $('#categoryModal').removeClass('flex').addClass('hidden');
        });

        // Edit category
        $('.category-action-btn.edit').off('click').on('click', function() {
            const categoryId = $(this).closest('.category-item').data('id');
            const category = categories.find(cat => cat.id === categoryId);
            
            if (category) {
                currentCategoryId = categoryId;
                $('#modalTitle').text('Edit Category');
                $('#categoryName').val(category.category_name);
                populateParentSelect(category.parent_id, categoryId);
                $('#categoryModal').removeClass('hidden').addClass('flex');
            }
        });

        // Delete category
        $('.category-action-btn.delete').off('click').on('click', function() {
            const categoryId = $(this).closest('.category-item').data('id');
            currentCategoryId = categoryId;
            $('#deleteModal').removeClass('hidden').addClass('flex');
        });

        // Cancel delete
        $('#cancelDelete').off('click').on('click', function() {
            $('#deleteModal').removeClass('flex').addClass('hidden');
        });

        // Confirm delete
        $('#confirmDelete').off('click').on('click', function() {
            if (currentCategoryId) {
                deleteCategory(currentCategoryId);
            }
        });

        // Search categories
        $('#searchCategory').off('input').on('input', function() {
            const searchTerm = $(this).val().toLowerCase();
            if (searchTerm) {
                const filtered = categories.filter(cat => 
                    cat.category_name.toLowerCase().includes(searchTerm)
                );
                renderCategoryTree(filtered);
            } else {
                renderCategoryTree();
            }
        });
    }

    // Populate parent category select with circular reference prevention
    function populateParentSelect(selectedParentId = null, excludeId = null) {
        const select = $('#parentCategory');
        select.empty().append('<option value="">None</option>');
        
        categories.forEach(category => {
            // Exclude current category and its descendants to prevent circular references
            if (category.id !== excludeId && !isDescendant(category.id, excludeId)) {
                const option = $(`<option value="${category.id}">${category.category_name}</option>`);
                if (category.id === selectedParentId) {
                    option.prop('selected', true);
                }
                select.append(option);
            }
        });
    }

    // Check if a category is a descendant of another category
    function isDescendant(categoryId, ancestorId) {
        if (!ancestorId) return false;
        
        const visited = new Set();
        const children = categories.filter(cat => cat.parent_id === ancestorId);
        
        for (const child of children) {
            if (child.id === categoryId) return true;
            if (!visited.has(child.id)) {
                visited.add(child.id);
                if (isDescendant(categoryId, child.id)) return true;
            }
        }
        
        return false;
    }

    // Handle form submission with validation
    $('#categoryForm').off('submit').on('submit', function(e) {
        e.preventDefault();
        
        const categoryData = {
            category_name: $('#categoryName').val().trim(),
            parent_id: $('#parentCategory').val() || null
        };

        // Validate hierarchy before submitting
        if (!validateCategoryHierarchy(currentCategoryId, categoryData.parent_id, categories)) {
            showToast('Invalid parent selection: This would create a circular reference', 'error');
            return;
        }

        if (currentCategoryId) {
            updateCategory(currentCategoryId, categoryData);
        } else {
            createCategory(categoryData);
        }
    });

    // Create category
    function createCategory(data) {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.categories,
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(response) {
                showToast('Category created successfully', 'success');
                $('#categoryModal').removeClass('flex').addClass('hidden');
                loadCategories();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to create category', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Update category
    function updateCategory(id, data) {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.category(id),
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(response) {
                showToast('Category updated successfully', 'success');
                $('#categoryModal').removeClass('flex').addClass('hidden');
                loadCategories();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to update category', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Delete category
    function deleteCategory(id) {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.category(id),
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
            success: function(response) {
                showToast('Category deleted successfully', 'success');
                $('#deleteModal').removeClass('flex').addClass('hidden');
                loadCategories();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to delete category', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Helper functions
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
});