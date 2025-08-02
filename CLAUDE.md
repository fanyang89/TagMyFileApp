# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Wails desktop application** built with Go backend and Vue.js + TypeScript frontend. The app is called "TagMyFileApp" and appears to be a file tagging application based on the project name.

## Architecture

- **Backend**: Go (1.23) using Wails v2 framework
  - Entry point: `main.go` - Sets up the Wails application window and binds the Go backend
  - Core logic: `app.go` - Contains the main App struct and business logic methods
  - Currently has a simple `Greet` method as an example of frontend-backend communication

- **Frontend**: Vue 3 + TypeScript + Vite
  - Located in `frontend/` directory
  - Uses Vue 3 Composition API with `<script setup>` syntax
  - Entry point: `frontend/src/main.ts`
  - Main component: `frontend/src/App.vue`
  - Example component: `frontend/src/components/HelloWorld.vue` demonstrates calling Go backend methods

- **Build System**: Wails handles the entire build process
  - Configuration: `wails.json`
  - Frontend dependencies managed via npm in `frontend/package.json`

## Development Commands

### Development Mode
```bash
wails dev
```
- Runs the application in development mode with hot reload
- Frontend dev server runs on http://localhost:34115 for browser development
- Automatically rebuilds frontend when files change

### Building
```bash
wails build
```
- Creates a production build for the current platform
- Outputs executable to `build/bin/` directory

### Frontend-only Development
```bash
cd frontend
npm install     # Install dependencies
npm run dev     # Start Vite dev server
npm run build   # Build frontend for production
```

## Key Files and Patterns

### Backend-Backend Communication
- Go methods in `app.go` are automatically exposed to frontend via Wails
- TypeScript bindings are generated in `frontend/wailsjs/go/main/App.d.ts`
- Frontend calls Go methods using promises (e.g., `Greet(name).then(result => ...)`)

### Project Structure
```
├── main.go              # Application entry point
├── app.go               # Main application logic
├── wails.json           # Wails configuration
├── go.mod               # Go module dependencies
├── frontend/            # Vue.js frontend
│   ├── src/
│   │   ├── App.vue      # Main Vue component
│   │   ├── main.ts      # Frontend entry point
│   │   └── components/  # Vue components
│   ├── dist/            # Built frontend assets
│   └── package.json     # Frontend dependencies
└── build/               # Build artifacts and platform-specific files
```

## Development Notes

- The application uses Wails' built-in asset server to serve the frontend
- Frontend changes are automatically reloaded during development
- TypeScript types for Go backend methods are auto-generated
- The default window size is 1024x768 with dark background (#1B2636)
- No existing linting or testing commands are configured in package.json