// Global state
let currentTrip = null;
let allProperties = [];
let allItems = [];

// API Base URL
const API_BASE = '/api';

// Utility Functions
function showLoading() {
    document.getElementById('loading').classList.add('active');
}

function hideLoading() {
    document.getElementById('loading').classList.remove('active');
}

function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
}

// API Functions
async function apiCall(endpoint, method = 'GET', body = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE}${endpoint}`, options);
    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || 'Request failed');
    }

    return data;
}

async function loadTrips() {
    try {
        showLoading();
        const data = await apiCall('/trips');
        renderTripList(data.trips || []);
    } catch (error) {
        showToast('Failed to load trips: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function createTrip(tripData) {
    try {
        showLoading();
        await apiCall('/trips', 'POST', tripData);
        showToast('Trip created successfully!');
        showPage('main-menu');
        loadTrips();
    } catch (error) {
        showToast('Failed to create trip: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function loadTrip(tripName) {
    try {
        showLoading();
        currentTrip = tripName;
        const data = await apiCall(`/trips/${encodeURIComponent(tripName)}`);
        renderTripDetails(data);
        showPage('trip-details');
    } catch (error) {
        showToast('Failed to load trip: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function updateTrip(tripName, updates) {
    try {
        showLoading();
        await apiCall(`/trips/${encodeURIComponent(tripName)}/update`, 'PUT', updates);
        showToast('Trip updated successfully!');
        loadTrip(tripName);
    } catch (error) {
        showToast('Failed to update trip: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function loadProperties(tripName) {
    try {
        showLoading();
        const data = await apiCall(`/trips/${encodeURIComponent(tripName)}/properties`);
        allProperties = data.properties || [];
        renderProperties(allProperties);
        showPage('properties-page');
    } catch (error) {
        showToast('Failed to load properties: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function toggleProperty(tripName, propertyName) {
    try {
        await apiCall(`/trips/${encodeURIComponent(tripName)}/properties/${encodeURIComponent(propertyName)}/toggle`, 'POST');
        // Update UI optimistically
        const property = allProperties.find(p => p.name === propertyName);
        if (property) {
            property.active = !property.active;
            renderProperties(allProperties);
        }
    } catch (error) {
        showToast('Failed to toggle property: ' + error.message, 'error');
        // Reload to ensure consistency
        loadProperties(tripName);
    }
}

async function loadItems(tripName) {
    try {
        showLoading();
        const data = await apiCall(`/trips/${encodeURIComponent(tripName)}/items`);
        allItems = data.categories || [];
        renderPackingList(allItems);
        showPage('packing-page');
    } catch (error) {
        showToast('Failed to load items: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function toggleItem(tripName, itemCode) {
    try {
        await apiCall(`/trips/${encodeURIComponent(tripName)}/items/${encodeURIComponent(itemCode)}/toggle`, 'POST');
        // Update UI optimistically
        for (const category of allItems) {
            const item = category.items.find(i => i.code === itemCode);
            if (item) {
                item.packed = !item.packed;
                break;
            }
        }
        renderPackingList(allItems);
    } catch (error) {
        showToast('Failed to toggle item: ' + error.message, 'error');
        // Reload to ensure consistency
        loadItems(tripName);
    }
}

// Render Functions
function renderTripList(trips) {
    const tripList = document.getElementById('trip-list');

    if (trips.length === 0) {
        tripList.innerHTML = `
            <div class="empty-state">
                <p>📭 No trips yet</p>
                <p style="font-size: 0.9rem; opacity: 0.8;">Create your first packing list!</p>
            </div>
        `;
        return;
    }

    tripList.innerHTML = trips.map(trip => `
        <div class="trip-card" onclick="loadTrip('${trip.replace(/\.[^.]+$/, '')}')">
            <h3>📋 ${trip.replace(/\.[^.]+$/, '')}</h3>
            <p>Click to view and pack</p>
        </div>
    `).join('');
}

function renderTripDetails(trip) {
    document.getElementById('trip-title').textContent = trip.name;
    document.getElementById('info-nights').textContent = trip.nights;
    document.getElementById('info-temp').textContent = `${trip.temperatureMin}°F - ${trip.temperatureMax}°F`;

    const properties = trip.properties.length > 0
        ? trip.properties.join(', ')
        : 'None';
    document.getElementById('info-properties').textContent = properties;
}

function renderProperties(properties) {
    const container = document.getElementById('properties-list');

    container.innerHTML = properties.map(prop => `
        <div class="property-item ${prop.active ? 'active' : ''}"
             onclick="toggleProperty('${currentTrip}', '${prop.name}')">
            <div class="checkbox"></div>
            <div class="info">
                <div class="name">${prop.name}</div>
                ${prop.description ? `<div class="description">${prop.description}</div>` : ''}
            </div>
        </div>
    `).join('');
}

function renderPackingList(categories) {
    const container = document.getElementById('packing-list');
    const hidePacked = document.getElementById('hide-packed').checked;

    let totalItems = 0;
    let packedItems = 0;

    container.innerHTML = categories.map(category => {
        const visibleItems = category.items.filter(item => {
            totalItems++;
            if (item.packed) packedItems++;
            return !hidePacked || !item.packed;
        });

        if (visibleItems.length === 0) return '';

        const categoryPacked = category.items.filter(i => i.packed).length;

        return `
            <div class="category">
                <div class="category-header">
                    <span>${category.name}</span>
                    <span class="category-count">${categoryPacked}/${category.items.length}</span>
                </div>
                <div class="category-items">
                    ${visibleItems.map(item => `
                        <div class="item ${item.packed ? 'packed' : ''}"
                             onclick="toggleItem('${currentTrip}', '${item.code}')">
                            <div class="checkbox"></div>
                            <div class="info">
                                <div class="name">${item.name}</div>
                                ${item.count > 1 ? `<div class="count">Quantity: ${item.count}</div>` : ''}
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }).join('');

    document.getElementById('pack-stats').textContent = `${packedItems}/${totalItems} packed`;
}

function filterProperties() {
    const searchTerm = document.getElementById('property-search').value.toLowerCase();
    const filtered = allProperties.filter(prop =>
        prop.name.toLowerCase().includes(searchTerm) ||
        (prop.description && prop.description.toLowerCase().includes(searchTerm))
    );
    renderProperties(filtered);
}

// Page Navigation Functions
function showPropertiesPage() {
    loadProperties(currentTrip);
}

function showPackingPage() {
    loadItems(currentTrip);
}

function showEditTripPage() {
    // Load current values into the edit form
    apiCall(`/trips/${encodeURIComponent(currentTrip)}`).then(data => {
        document.getElementById('edit-nights').value = data.nights;
        document.getElementById('edit-temp-min').value = data.temperatureMin;
        document.getElementById('edit-temp-max').value = data.temperatureMax;
        showPage('edit-trip-page');
    }).catch(error => {
        showToast('Failed to load trip data: ' + error.message, 'error');
    });
}

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    // Load initial trips
    loadTrips();

    // New trip button
    document.getElementById('new-trip-btn').addEventListener('click', () => {
        showPage('new-trip-page');
        document.getElementById('new-trip-form').reset();
    });

    // New trip form submission
    document.getElementById('new-trip-form').addEventListener('submit', (e) => {
        e.preventDefault();

        const tripData = {
            name: document.getElementById('trip-name').value,
            nights: parseInt(document.getElementById('trip-nights').value),
            temperatureMin: parseInt(document.getElementById('trip-temp-min').value),
            temperatureMax: parseInt(document.getElementById('trip-temp-max').value),
            properties: []
        };

        createTrip(tripData);
    });

    // Edit trip form submission
    document.getElementById('edit-trip-form').addEventListener('submit', (e) => {
        e.preventDefault();

        const updates = {
            nights: parseInt(document.getElementById('edit-nights').value),
            temperatureMin: parseInt(document.getElementById('edit-temp-min').value),
            temperatureMax: parseInt(document.getElementById('edit-temp-max').value)
        };

        updateTrip(currentTrip, updates);
    });

    // Property search
    document.getElementById('property-search').addEventListener('input', filterProperties);

    // Hide packed toggle
    document.getElementById('hide-packed').addEventListener('change', () => {
        renderPackingList(allItems);
    });
});
