// Alpine.js app
function packingApp() {
    return {
        // State
        currentPage: 'main-menu',
        loading: false,
        trips: [],
        currentTrip: null,
        currentTripName: null,
        properties: [],
        propertySearch: '',
        categories: [],
        hidePacked: false,
        collapsedCategories: new Set(),
        refreshInterval: null,
        toast: {
            show: false,
            message: '',
            type: 'success'
        },
        newTrip: {
            name: '',
            nights: 2,
            temperatureMin: 60,
            temperatureMax: 80
        },
        editForm: {
            name: '',
            nights: 0,
            temperatureMin: 0,
            temperatureMax: 0
        },

        // Computed properties
        get filteredProperties() {
            if (!this.propertySearch) return this.properties;
            const search = this.propertySearch.toLowerCase();
            return this.properties.filter(p =>
                p.name.toLowerCase().includes(search) ||
                (p.description ?? '').toLowerCase().includes(search)
            );
        },

        get packedCount() {
            return this.categories.reduce((total, cat) =>
                total + cat.items.filter(i => i.packed).length, 0
            );
        },

        get totalItems() {
            return this.categories.reduce((total, cat) =>
                total + cat.items.length, 0
            );
        },

        // Lifecycle
        init() {
            this.loadTrips();
        },

        // Helper methods
        visibleItemsInCategory(category) {
            if (!this.hidePacked) return category.items.length;
            return category.items.filter(i => !i.packed).length;
        },

        packedInCategory(category) {
            return category.items.filter(i => i.packed).length;
        },

        toggleCategoryCollapse(categoryName) {
            if (this.collapsedCategories.has(categoryName)) {
                this.collapsedCategories.delete(categoryName);
            } else {
                this.collapsedCategories.add(categoryName);
            }
        },

        isCategoryCollapsed(categoryName) {
            return this.collapsedCategories.has(categoryName);
        },

        showToast(message, type = 'success') {
            this.toast = { show: true, message, type };
            setTimeout(() => {
                this.toast.show = false;
            }, 3000);
        },

        // API methods
        async apiCall(endpoint, method = 'GET', body = null) {
            const options = {
                method,
                headers: {
                    'Content-Type': 'application/json',
                },
            };

            if (body) {
                options.body = JSON.stringify(body);
            }

            const response = await fetch(`/api${endpoint}`, options);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        },

        // Trip operations
        async loadTrips() {
            try {
                this.loading = true;
                const data = await this.apiCall('/trips');
                this.trips = data.trips || [];
            } catch (error) {
                this.showToast('Failed to load trips: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        async createTrip() {
            try {
                this.loading = true;
                await this.apiCall('/trips', 'POST', {
                    name: this.newTrip.name,
                    nights: this.newTrip.nights,
                    temperatureMin: this.newTrip.temperatureMin,
                    temperatureMax: this.newTrip.temperatureMax,
                    properties: []
                });
                this.showToast('Trip created successfully!');
                this.newTrip = { name: '', nights: 2, temperatureMin: 60, temperatureMax: 80 };
                await this.loadTrips();
                this.currentPage = 'main-menu';
            } catch (error) {
                this.showToast('Failed to create trip: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        async loadTrip(tripName) {
            try {
                this.loading = true;
                this.currentTripName = tripName;
                const data = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}`);
                this.currentTrip = data;
                this.currentPage = 'trip-details';
            } catch (error) {
                this.showToast('Failed to load trip: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        showEditTrip() {
            this.editForm = {
                name: this.currentTrip.name,
                nights: this.currentTrip.nights,
                temperatureMin: this.currentTrip.temperatureMin,
                temperatureMax: this.currentTrip.temperatureMax
            };
            this.currentPage = 'edit-trip';
        },

        async updateTrip() {
            try {
                this.loading = true;
                const result = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/update`, 'PUT', this.editForm);
                // If the name changed, the backend returns the new name
                if (result.name) {
                    this.currentTripName = result.name;
                }
                this.showToast('Trip updated successfully!');
                await this.loadTrip(this.currentTripName);
                this.currentPage = 'trip-details';
            } catch (error) {
                this.showToast('Failed to update trip: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        backToMainMenu() {
            this.currentTrip = null;
            this.currentTripName = null;
            this.loadTrips();
            this.currentPage = 'main-menu';
        },

        // Properties operations
        async loadProperties() {
            try {
                this.loading = true;
                const data = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/properties`);
                this.properties = data.properties || [];
                this.propertySearch = '';
                this.currentPage = 'properties';
            } catch (error) {
                this.showToast('Failed to load properties: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        async toggleProperty(propertyName) {
            try {
                await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/properties/${encodeURIComponent(propertyName)}/toggle`, 'POST');
                // Reload all properties to catch any automatic changes
                const data = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/properties`);
                this.properties = data.properties || [];
            } catch (error) {
                this.showToast('Failed to toggle property: ' + error.message, 'error');
                this.loadProperties();
            }
        },

        // Items operations
        async loadItems() {
            try {
                this.loading = true;
                const data = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/items`);
                this.categories = data.categories || [];
                this.hidePacked = false;
                this.collapsedCategories.clear();
                this.currentPage = 'packing';
                this.startAutoRefresh();
            } catch (error) {
                this.showToast('Failed to load items: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        async refreshItems() {
            // Silently refresh items without changing UI state
            try {
                const data = await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/items`);
                this.categories = data.categories || [];
            } catch (error) {
                // Silent failure - don't interrupt user
                console.error('Auto-refresh failed:', error);
            }
        },

        startAutoRefresh() {
            this.stopAutoRefresh();
            // Refresh every 10 seconds
            this.refreshInterval = setInterval(() => {
                if (this.currentPage === 'packing') {
                    this.refreshItems();
                }
            }, 10000);
        },

        stopAutoRefresh() {
            if (this.refreshInterval) {
                clearInterval(this.refreshInterval);
                this.refreshInterval = null;
            }
        },

        async toggleItem(itemCode) {
            try {
                await this.apiCall(`/trips/${encodeURIComponent(this.currentTripName)}/items/${encodeURIComponent(itemCode)}/toggle`, 'POST');
                // Update local state optimistically
                for (const category of this.categories) {
                    const item = category.items.find(i => i.code === itemCode);
                    if (item) {
                        item.packed = !item.packed;
                        // Check if all items are now packed
                        if (this.packedCount === this.totalItems && this.totalItems > 0) {
                            this.triggerFireworks();
                        }
                        break;
                    }
                }
            } catch (error) {
                this.showToast('Failed to toggle item: ' + error.message, 'error');
                this.loadItems();
            }
        },

        triggerFireworks() {
            var duration = 10 * 1000;
            var end = Date.now() + duration;

            (function frame() {
              // launch a few confetti from the left edge
              confetti({
                particleCount: 7,
                angle: 60,
                spread: 180,
                origin: { x: 0 }
              });
              // and launch a few from the right edge
              confetti({
                particleCount: 7,
                angle: 120,
                spread: 180,
                origin: { x: 1 }
              });

              // keep going until we are out of time
              if (Date.now() < end) {
                requestAnimationFrame(frame);
              }
            }());
        }
    };
}
