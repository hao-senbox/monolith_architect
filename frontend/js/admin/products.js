$(document).ready(function() {
    console.log('Document ready - jQuery version:', $.fn.jquery);
    const BASE_URL = 'https://monolith-architect.onrender.com';
    // const BASE_URL = 'http://localhost:8003';
    const API_ENDPOINTS = {
        products: `${BASE_URL}/api/v1/product`,
        product: (id) => `${BASE_URL}/api/v1/product/${id}`,
        categories: `${BASE_URL}/api/v1/category`
    };

    let products = [];
    let categories = [];
    let currentProductId = null;
    let sizeCount = 0;

    // Check if user is logged in
    if (!getToken()) {
        window.location.href = '/frontend/pages/auth/login.html';
        return;
    }

    // Load layout
    $('#layout').load('/frontend/components/layout.html', function() {
        setActiveNavItem('/admin/products.html');
        setPageTitle('Product Management');
        loadCategories();
        loadProducts();
    });

    // Load categories
    function loadCategories() {
        $.ajax({
            url: API_ENDPOINTS.categories,
            method: 'GET',
            success: function(response) {
                categories = response.data || [];
                populateCategorySelects();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to load categories', 'error');
            }
        });
    }

    // Load products
    function loadProducts() {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.products,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            success: function(response) {
                products = response.data || [];
                renderProductTable();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to load products', 'error');
                $('#productTableBody').html('<tr><td colspan="6" class="px-6 py-4 text-center error-state">Failed to load products</td></tr>');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Render product table
    function renderProductTable(filteredProducts = null) {
        const productsToRender = filteredProducts || products;
        
        if (productsToRender.length === 0) {
            $('#productTableBody').html('<tr><td colspan="6" class="px-6 py-4 text-center empty-state">No products found</td></tr>');
        }

        let html = '';
        productsToRender.forEach(product => {
            const mainImage = product.main_image || '';
            const category = categories.find(cat => cat.id === product.category_id);
            const priceRange = getPriceRange(product.sizes);
            const totalStock = getTotalStock(product.sizes);
            const stockStatus = getStockStatus(totalStock);

            html += `
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <img src="${mainImage}" alt="${product.product_name}" class="product-image">
                    </td>
                    <td class="px-6 py-4">
                        <div class="product-name">${product.product_name}</div>
                        <div class="text-sm text-gray-500">${product.product_description}</div>
                    </td>
                    <td class="px-6 py-4">
                        <div class="product-category">${category ? category.category_name : 'N/A'}</div>
                    </td>
                    <td class="px-6 py-4">
                        <div class="product-price">${priceRange}</div>
                    </td>
                    <td class="px-6 py-4">
                        <div class="product-stock ${stockStatus.class}">${stockStatus.text}</div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <button class="action-btn edit mr-2" data-id="${product.id}">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="action-btn delete" data-id="${product.id}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
        });

        $('#productTableBody').html(html);
        setupEventListeners();
    }

    // Helper functions for product table
    function getPriceRange(sizes) {
        if (!sizes || sizes.length === 0) return 'N/A';
        
        let minPrice = Infinity;
        let maxPrice = -Infinity;
        let currency = '';

        sizes.forEach(size => {
            const price = size.price * (1 - size.discount / 100);
            if (price < minPrice) minPrice = price;
            if (price > maxPrice) maxPrice = price;
            currency = size.currency;
        });

        if (minPrice === Infinity || maxPrice === -Infinity) return 'N/A';
        if (minPrice === maxPrice) return `${minPrice.toFixed(2)} ${currency}`;
        return `${minPrice.toFixed(2)} - ${maxPrice.toFixed(2)} ${currency}`;
    }

    function getTotalStock(sizes) {
        if (!sizes || sizes.length === 0) return 0;
        return sizes.reduce((total, size) => total + (size.stock || 0), 0);
    }

    function getStockStatus(stock) {
        if (stock === 0) return { text: 'Out of Stock', class: 'out-of-stock' };
        if (stock < 10) return { text: 'Low Stock', class: 'low-stock' };
        return { text: 'In Stock', class: 'in-stock' };
    }

    // Populate category selects
    function populateCategorySelects() {
        const select = $('#categoryId');
        const filter = $('#categoryFilter');
        
        select.empty().append('<option value="">Select Category</option>');
        filter.empty().append('<option value="">All Categories</option>');
        
        categories.forEach(category => {
            select.append(`<option value="${category.id}">${category.category_name}</option>`);
            filter.append(`<option value="${category.id}">${category.category_name}</option>`);
        });
    }

    // Setup event listeners
    function setupEventListeners() {
        // Function to show modal
        function showProductModal() {
            const modal = document.getElementById('productModal');
            if (!modal) {
                console.error('Modal element not found!');
                return;
            }
            modal.classList.remove('hidden');
            modal.classList.add('flex');
            console.log('Modal classes after show:', modal.className);
        }

        // Add product button
        $('#addProductBtn').off('click').on('click', function() {
            currentProductId = null;
            sizeCount = 0;
            $('#modalTitle').text('Add Product');
            $('#productForm')[0].reset();
            $('#sizesList').empty();
            addSize();
            showProductModal();
        });

        // Close modal
        $('#closeModal, #cancelBtn').off('click').on('click', function() {
            $('#productModal').removeClass('flex').addClass('hidden');
        });

        // Add size button
        $('#addSizeBtn').off('click').on('click', addSize);

        // Remove size
        $(document).on('click', '.remove-size', function() {
            if ($('.size-item').length > 1) {
                $(this).closest('.size-item').remove();
                updateSizeNumbers();
            } else {
                showToast('At least one size option is required', 'error');
            }
        });

        // Edit product
        $(document).on('click', '.action-btn.edit', function() {
            const productId = $(this).data('id');
            const product = products.find(p => p.id === productId);
            
            if (product) {
                currentProductId = productId;
                sizeCount = product.sizes ? product.sizes.length : 0;
                $('#modalTitle').text('Edit Product');
                $('#productName').val(product.product_name);
                $('#productDescription').val(product.product_description);
                $('#categoryId').val(product.category_id);
                $('#color').val(product.color);
                
                $('#sizesList').empty();
                if (product.sizes && product.sizes.length > 0) {
                    product.sizes.forEach((size, index) => {
                        addSize(size, index);
                    });
                } else {
                    addSize();
                }
                
                $('#productModal').removeClass('hidden').addClass('flex');
            }
        });

        // Delete product
        $(document).on('click', '.action-btn.delete', function() {
            const productId = $(this).data('id');
            currentProductId = productId;
            $('#deleteModal').removeClass('hidden').addClass('flex');
        });

        // Cancel delete
        $('#cancelDelete').off('click').on('click', function() {
            $('#deleteModal').removeClass('flex').addClass('hidden');
        });

        // Confirm delete
        $('#confirmDelete').off('click').on('click', function() {
            if (currentProductId) {
                deleteProduct(currentProductId);
            }
        });

        // Search products
        $('#searchProduct').off('input').on('input', function() {
            const searchTerm = $(this).val().toLowerCase();
            const categoryFilter = $('#categoryFilter').val();
            
            let filtered = products;
            
            if (searchTerm) {
                filtered = filtered.filter(product => 
                    product.product_name.toLowerCase().includes(searchTerm)
                );
            }
            
            if (categoryFilter) {
                filtered = filtered.filter(product => 
                    product.category_id === categoryFilter
                );
            }
            
            renderProductTable(filtered);
        });

        // Category filter
        $('#categoryFilter').off('change').on('change', function() {
            const searchTerm = $('#searchProduct').val().toLowerCase();
            const categoryFilter = $(this).val();
            
            let filtered = products;
            
            if (searchTerm) {
                filtered = filtered.filter(product => 
                    product.product_name.toLowerCase().includes(searchTerm)
                );
            }
            
            if (categoryFilter) {
                filtered = filtered.filter(product => 
                    product.category_id === categoryFilter
                );
            }
            
            renderProductTable(filtered);
        });
    }

    // Add size
    function addSize(size = null, index = null) {
        const sizeIndex = index !== null ? index : $('.size-item').length;
        const template = $('#sizeTemplate').html()
            .replace(/sizes\[0\]/g, `sizes[${sizeIndex}]`);
        
        const sizeElement = $(template);
        sizeElement.find('.size-number').text(sizeIndex + 1);
        
        if (size) {
            sizeElement.find('input[name$="[size]"]').val(size.size);
            sizeElement.find('input[name$="[sku]"]').val(size.sku);
            sizeElement.find('input[name$="[price]"]').val(size.price);
            sizeElement.find('select[name$="[currency]"]').val(size.currency);
            sizeElement.find('input[name$="[discount]"]').val(size.discount);
            sizeElement.find('input[name$="[stock]"]').val(size.stock);
        }
        
        $('#sizesList').append(sizeElement);
    }

    // Update size numbers
    function updateSizeNumbers() {
        $('.size-item').each(function(index) {
            $(this).find('.size-number').text(index + 1);
            $(this).find('input, select').each(function() {
                const name = $(this).attr('name');
                if (name) {
                    $(this).attr('name', name.replace(/sizes\[\d+\]/, `sizes[${index}]`));
                }
            });
        });
    }

    // Handle form submission
    $('#productForm').off('submit').on('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const sizeCount = $('.size-item').length;
        formData.append('size_count', sizeCount);

        // Log FormData entries
        for (const pair of formData.entries()) {
            console.log(pair[0] + ': ' + pair[1]);
        }

        if (currentProductId) {
            updateProduct(currentProductId, formData);
        } else {
            createProduct(formData);
        }
    });

    // Create product
    function createProduct(formData) {
        console.log('Creating product...', formData);
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.products,
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            xhrFields: {
                withCredentials: true
            },
            data: formData,
            processData: false,
            contentType: false,
            success: function(response) {
                showToast('Product created successfully', 'success');
                $('#productModal').removeClass('flex').addClass('hidden');
                loadProducts();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to create product', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Update product
    function updateProduct(id, formData) {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.product(id),
            method: 'PUT',
            headers: {
                'Accept': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            xhrFields: {
                withCredentials: true
            },
            data: formData,
            processData: false,
            contentType: false,
            success: function(response) {
                showToast('Product updated successfully', 'success');
                $('#productModal').removeClass('flex').addClass('hidden');
                loadProducts();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to update product', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }

    // Delete product
    function deleteProduct(id) {
        showLoading();
        $.ajax({
            url: API_ENDPOINTS.product(id),
            method: 'DELETE',
            headers: {
                'Accept': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            xhrFields: {
                withCredentials: true
            },
            success: function(response) {
                showToast('Product deleted successfully', 'success');
                $('#deleteModal').removeClass('flex').addClass('hidden');
                loadProducts();
            },
            error: function(xhr) {
                if (xhr.status === 401) {
                    window.location.href = '/frontend/pages/auth/login.html';
                    return;
                }
                showToast('Failed to delete product', 'error');
            },
            complete: function() {
                hideLoading();
            }
        });
    }
});