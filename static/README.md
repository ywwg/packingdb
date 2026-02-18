# PackingDB Web Frontend

A mobile-friendly web interface for PackingDB.

## Features

- 📱 **Mobile-First Design**: Responsive interface optimized for phones
- 🎒 **Trip Management**: Create and manage packing lists
- ⚙️ **Property Configuration**: Toggle trip properties to customize your packing list
- ✅ **Pack Items**: Check off items as you pack them
- 💾 **Auto-Save**: All changes are automatically saved

## Running the Web App

1. Build and run the server:
```bash
cd cmd/packingweb
go build
./packingweb
```

2. Open your browser to `http://localhost:8080`

The server will create a `public/trips/` directory to store your packing lists.

## API Endpoints

### Trips
- `GET /api/trips` - List all trips
- `POST /api/trips` - Create a new trip
- `GET /api/trips/{name}` - Get trip details
- `PUT /api/trips/{name}/update` - Update trip settings (nights, temperatures)

### Properties
- `GET /api/properties` - List all available properties
- `GET /api/trips/{name}/properties` - Get trip properties
- `POST /api/trips/{name}/properties/{property}/toggle` - Toggle a property

### Items
- `GET /api/trips/{name}/items` - Get packing items for a trip
- `POST /api/trips/{name}/items/{code}/toggle` - Toggle item packed status

## Development

The web frontend consists of:
- `cmd/packingweb/main.go` - Go server with REST API
- `static/index.html` - Main HTML page
- `static/styles.css` - Mobile-optimized CSS
- `static/app.js` - JavaScript client

## Usage

1. **Create a Trip**: Click "Create New Trip" and enter trip details
2. **Configure Properties**: Select trip properties (e.g., Camping, Swimming, Business)
3. **Pack Items**: Check off items as you pack them
4. **Edit Settings**: Modify nights and temperatures as needed

The app automatically saves all changes, so you can close and reopen your trip at any time.
