$(document).ready(function() {
    const BASE_URL = 'https://monolith-architect.onrender.com';
    const API_ENDPOINTS = {
        products: `${BASE_URL}/api/v1/product`,
        product: (id) => `${BASE_URL}/api/v1/product/${id}`,
        categories: `${BASE_URL}/api/v1/category`
    };

    let products = [];
    let categories = [];
    let currentProductId = null;
    let variantCount = 0;

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
            xhrFields: {
                withCredentials: true
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
            return;
        }

        let html = '';
        productsToRender.forEach(product => {
            const mainImage = product.variants[0]?.main_image || '';
            const category = categories.find(cat => cat.id === product.category_id);
            const priceRange = getPriceRange(product.variants);
            const totalStock = getTotalStock(product.variants);
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
    function getPriceRange(variants) {
        let minPrice = Infinity;
        let maxPrice = -Infinity;
        let currency = '';

        variants.forEach(variant => {
            variant.sizes.forEach(size => {
                const price = size.price * (1 - size.discount / 100);
                if (price < minPrice) minPrice = price;
                if (price > maxPrice) maxPrice = price;
                currency = size.currency;
            });
        });

        if (minPrice === Infinity || maxPrice === -Infinity) return 'N/A';
        if (minPrice === maxPrice) return `${minPrice.toFixed(2)} ${currency}`;
        return `${minPrice.toFixed(2)} - ${maxPrice.toFixed(2)} ${currency}`;
    }

    function getTotalStock(variants) {
        return variants.reduce((total, variant) => {
            return total + variant.sizes.reduce((variantTotal, size) => variantTotal + size.stock, 0);
        }, 0);
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
        // Add product button
        $('#addProductBtn').off('click').on('click', function() {
            currentProductId = null;
            variantCount = 0;
            $('#modalTitle').text('Add Product');
            $('#productForm')[0].reset();
            $('#variantsList').empty();
            addVariant();
            $('#productModal').removeClass('hidden').addClass('flex');
        });

        // Close modal
        $('#closeModal, #cancelBtn').off('click').on('click', function() {
            $('#productModal').removeClass('flex').addClass('hidden');
        });

        // Add variant button
        $('#addVariantBtn').off('click').on('click', addVariant);

        // Remove variant
        $(document).on('click', '.remove-variant', function() {
            $(this).closest('.variant-item').remove();
            updateVariantNumbers();
        });

        // Add size
        $(document).on('click', '.add-size', function() {
            const variantItem = $(this).closest('.variant-item');
            const sizesList = variantItem.find('.sizes-list');
            const variantIndex = variantItem.index();
            addSize(sizesList, variantIndex);
        });

        // Remove size
        $(document).on('click', '.remove-size', function() {
            $(this).closest('.size-item').remove();
        });

        // Edit product
        $(document).on('click', '.action-btn.edit', function() {
            const productId = $(this).data('id');
            const product = products.find(p => p.id === productId);
            
            if (product) {
                currentProductId = productId;
                variantCount = product.variants.length;
                $('#modalTitle').text('Edit Product');
                $('#productName').val(product.product_name);
                $('#productDescription').val(product.product_description);
                $('#categoryId').val(product.category_id);
                
                $('#variantsList').empty();
                product.variants.forEach((variant, index) => {
                    addVariant(variant, index);
                });
                
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

    // Add variant
    function addVariant(variant = null, index = null) {
        const variantIndex = index !== null ? index : variantCount;
        const template = $('#variantTemplate').html()
            .replace(/variants\[0\]/g, `variants[${variantIndex}]`);
        
        const variantElement = $(template);
        variantElement.find('.variant-number').text(variantIndex + 1);
        
        if (variant) {
            variantElement.find('input[name$="[color]"]').val(variant.color);
            
            const sizesList = variantElement.find('.sizes-list');
            variant.sizes.forEach((size, sizeIndex) => {
                addSize(sizesList, variantIndex, size, sizeIndex);
            });
        } else {
            const sizesList = variantElement.find('.sizes-list');
            addSize(sizesList, variantIndex);
        }
        
        $('#variantsList').append(variantElement);
        variantCount++;
    }

    // Add size
    function addSize(container, variantIndex, size = null, sizeIndex = null) {
        const sizeIndex2 = sizeIndex !== null ? sizeIndex : container.children().length;
        const template = $('#sizeTemplate').html()
            .replace(/variants\[0\]/g, `variants[${variantIndex}]`)
            .replace(/sizes\[0\]/g, `sizes[${sizeIndex2}]`);
        
        const sizeElement = $(template);
        
        if (size) {
            sizeElement.find('input[name$="[size]"]').val(size.size);
            sizeElement.find('input[name$="[sku]"]').val(size.sku);
            sizeElement.find('input[name$="[price]"]').val(size.price);
            sizeElement.find('select[name$="[currency]"]').val(size.currency);
            sizeElement.find('input[name$="[discount]"]').val(size.discount);
            sizeElement.find('input[name$="[stock]"]').val(size.stock);
        }
        
        container.append(sizeElement);
    }

    // Update variant numbers
    function updateVariantNumbers() {
        $('.variant-item').each(function(index) {
            $(this).find('.variant-number').text(index + 1);
            $(this).find('input, select').each(function() {
                const name = $(this).attr('name');
                if (name) {
                    $(this).attr('name', name.replace(/variants\[\d+\]/, `variants[${index}]`));
                }
            });
        });
        variantCount = $('.variant-item').length;
    }

    // Handle form submission
    $('#productForm').off('submit').on('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        formData.append('variant_count', variantCount);
        
        // Get size count for each variant
        $('.variant-item').each(function(index) {
            const sizeCount = $(this).find('.size-item').length;
            formData.append('size_options', sizeCount);
        });

        if (currentProductId) {
            updateProduct(currentProductId, formData);
        } else {
            createProduct(formData);
        }
    });

    // Create product
    function createProduct(formData) {
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