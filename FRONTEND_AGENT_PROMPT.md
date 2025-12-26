# Frontend Agent Prompt (Separate Repository)

You are an expert Vue 3 and TailwindCSS developer building the frontend for the Work Track application.

## Project Context: Separate Repositories
- **Backend**: Go API (Running at `http://localhost:8080` or `https://work-track-api.onrender.com`)
- **Frontend**: Vue 3 application (This is your NEW repository)

## Your Task
Initialize and build the Vue 3 frontend in this new repository.

### Requirements
1. **Tech Stack**:
   - Vue 3 (Composition API, `<script setup>`)
   - Vite (Build tool)
   - TailwindCSS (Styling)
   - Vue Router (Routing)
   - Pinia (State Management)
   - Axios (HTTP Client)

2. **Design System**:
   - Use a modern, premium design (dark mode, glassmorphism, smooth animations).
   - **Do NOT** use component libraries (like Vuetify/ElementUI) unless requested. Use headless UI + Tailwind.

3. **Features to Implement**:
   - **Authentication**: Login/Register pages.
   - **Dashboard**: View daily track items.
   - **Tracking**: Add/Edit/Delete track items (Type, Emergency Call, Holiday Call, Hours, Shifts).
   - **Date Navigation**: Switch between days/weeks.

4. **Development Workflow**:
   - Run `npm create vite@latest .` (or `work-track-frontend`).
   - **CORS Handling**:
     - In `vite.config.js`, configure the proxy:
     ```js
     server: {
       proxy: {
         '/api': {
           target: 'http://localhost:8080',
           changeOrigin: true
         }
       }
     }
     ```

5. **Deliverables**:
   - Fully functional Vue 3 app.
   - `README.md` with setup instructions.

## Important
- **API Documentation**: Refer to the backend's `API_DOCUMENTATION.md` (provided by the user or available in the backend repo).
- **Authentication**: The backend uses JWT. Store the token in `localStorage` and send it in the `Authorization: Bearer <token>` header.
- **Error Handling**: Handle 401 (Unauthorized) by redirecting to login.
- Make it look **amazing**!
