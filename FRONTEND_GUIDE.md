# 🎮 Frontend Setup & Running

## Quick Start

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Start Frontend Development Server

```bash
npm start
```

Server will be available at: **http://localhost:3000**

### 3. Make Sure Backend is Running

The frontend expects the backend API on **http://localhost:8080**

```bash
# In another terminal
go run cmd/server/main.go
```

---

## Project Structure

```
frontend/
├── public/
│   └── index.html              # HTML entry point
├── src/
│   ├── pages/
│   │   ├── Login.jsx          # OAuth login page
│   │   ├── Dashboard.jsx      # Game list & creation
│   │   └── Game.jsx           # Game play screen
│   ├── components/
│   │   ├── Navigation.jsx     # Top navigation bar
│   │   ├── GameCard.jsx       # Game card component
│   │   └── GameBoard.jsx      # Game board display
│   ├── App.jsx                # Main app component
│   ├── index.jsx              # React entry point
│   └── index.css              # Global styles
├── package.json               # Dependencies
├── tailwind.config.js         # Tailwind configuration
├── postcss.config.js          # PostCSS configuration
└── .gitignore                 # Git ignore rules
```

---

## Technology Stack

- **React 18** - UI library
- **React Router v6** - Page routing
- **Axios** - HTTP client
- **Tailwind CSS** - Styling
- **React Icons** - Icon library

---

## Features

### ✅ Authentication
- OAuth2 login (Google, GitHub, Discord)
- JWT token management
- Protected routes
- User profile display

### ✅ Game Management
- Create new games
- Join existing games
- Copy room codes
- Game list with filters
- Real-time game state updates

### ✅ Game Play
- Display game status
- Show current player
- Display player cards
- Place cards on opponents
- Draw cards from deck
- View table cards

### ✅ Responsive Design
- Mobile-friendly UI
- Tailwind CSS styling
- Dark theme
- Intuitive navigation

---

## Pages

### /login
- OAuth provider selection
- Login instructions
- Getting started guide

### /dashboard
- Game list
- Create game form
- Join game cards
- Refresh games

### /game/:id
- Game status display
- Player list
- Your cards
- Game board (Phase 2)
- Game actions (place, draw)

---

## Components

### Navigation
- User profile display
- Logout button
- User role badge

### GameCard
- Game information
- Room code display/copy
- Join button
- Player count display

### GameBoard
- Game status
- Current turn display
- Your cards
- Other players
- Table cards (Phase 2)
- Action buttons

---

## Environment Setup

### Development (Default)
```bash
npm start
```
- Frontend: http://localhost:3000
- Backend: http://localhost:8080 (proxied)

### Production
```bash
npm run build
```
- Creates optimized build in `build/` folder
- Run with: `npm run serve`

---

## API Integration

The frontend communicates with the backend API:

```javascript
// Configured in package.json proxy
"proxy": "http://localhost:8080"

// API calls example:
axios.get('/api/v1/games')
axios.post('/api/v1/auth/google')
axios.post('/api/v1/games/{id}/place')
```

---

## Styling

Uses **Tailwind CSS** with custom configuration:

```tailwindcss
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### Color Scheme
- Background: Dark gray (#0f172a - #1e293b)
- Primary: Blue (#2563eb)
- Success: Green (#16a34a)
- Error: Red (#dc2626)
- Accent: Purple, Yellow

---

## Available Scripts

### `npm start`
Runs the app in development mode
- Open http://localhost:3000 to view
- Page reloads on code changes
- Errors shown in console

### `npm test`
Runs the test suite in interactive mode

### `npm run build`
Builds the app for production
- Optimized minified build
- Ready for deployment

### `npm run eject`
⚠️ One-way operation - exposes all configuration

---

## Troubleshooting

### Port 3000 already in use
```bash
PORT=3001 npm start
```

### Backend not responding
1. Check backend is running: `go run cmd/server/main.go`
2. Check backend is on port 8080
3. Check CORS is enabled in backend

### OAuth not working
1. Verify OAuth credentials in `.env`
2. Check redirect URI in OAuth provider settings
3. Check backend is running

### Styles not loading
```bash
# Clear node modules and reinstall
rm -rf node_modules package-lock.json
npm install
npm start
```

### CORS errors
- Ensure backend has CORS middleware configured
- Check CORS_ALLOWED_ORIGINS includes http://localhost:3000

---

## Deployment

### Vercel (Recommended)
```bash
npm install -g vercel
vercel
```

### Netlify
```bash
npm run build
# Deploy the `build/` folder
```

### Docker
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

---

## Next Steps

1. **Customize Styling** - Update colors in `tailwind.config.js`
2. **Add More Features** - Create new pages/components
3. **Improve Game UI** - Add card animations, sounds
4. **Add Tests** - Write Jest/React Testing Library tests
5. **Performance** - Add code splitting, lazy loading

---

## Useful Links

- [React Documentation](https://react.dev)
- [React Router](https://reactrouter.com)
- [Tailwind CSS](https://tailwindcss.com)
- [Axios](https://axios-http.com)
- [React Icons](https://react-icons.github.io/react-icons)

---

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review React console for errors
3. Check backend API is running
4. Verify OAuth credentials are configured

---

**Happy Frontend Coding! 🚀**

