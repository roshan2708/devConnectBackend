# 🚀 DevConnect Frontend Master Prompt

You are building a **professional frontend for DevConnect**, a social platform for developers. Use this comprehensive guide to generate a complete, production-ready React application.

## 📋 Project Overview

**DevConnect** is a developer networking platform where professionals can:
- Create and share posts about their work
- Follow other developers
- Like and comment on posts
- View personalized feeds based on followers
- Search for developers by skills
- View trending posts
- Update their profile with bio, skills, GitHub link, and location
- Receive real-time notifications

---

## 🏗️ Technical Stack

- **Framework**: React 18+
- **Build Tool**: Vite
- **State Management**: Redux Toolkit or Context API
- **HTTP Client**: Axios
- **Styling**: Tailwind CSS
- **Responsive Design**: Mobile-first approach
- **Authentication**: JWT (stored in localStorage)
- **Routing**: React Router v6

---

## 🔌 Backend Connection Details

**Base URL**: `http://localhost:3000` (Development)

**Authentication**:
- Token stored as `authToken` in localStorage
- Include token in all protected requests via `Authorization: Bearer {token}` header
- OAuth 2.0 Google Login integration
- Token expiry: 30 days

---

## 📡 API Endpoints Reference

### Authentication Endpoints
```
GET  /auth/google                    - Initiate Google OAuth
GET  /auth/google/callback          - Google OAuth callback (returns JWT token)
GET  /health                         - Health check
```

### User Endpoints
```
GET  /me                            - Get current logged-in user
GET  /users/{userID}                - Get user profile by ID
PUT  /profile                       - Update current user profile (bio, github, skills, location)
GET  /users/{userID}/posts          - Get all posts by a user
GET  /search?skill=javascript       - Search users by skill
```

### Post Endpoints
```
POST /posts                         - Create a new post (Content: string)
GET  /posts?page=1&limit=10        - Get all posts (paginated)
PUT  /posts/{postID}                - Edit a post
DELETE /posts/{postID}              - Delete a post
GET  /trending                      - Get trending posts (top 10, scored by likes*2 + comments*3)
GET  /feed?page=1&limit=10         - Get personalized feed (only posts from followed users)
```

### Like Endpoints
```
POST /posts/{postID}/like           - Like a post
DELETE /posts/{postID}/like         - Unlike a post
GET  /posts/{postID}/likes          - Get all users who liked a post
```

### Comment Endpoints
```
POST /posts/{postID}/comment        - Add comment to post (Content: string)
GET  /posts/{postID}/comments       - Get all comments for a post
DELETE /comments/{commentID}        - Delete a comment
```

### Follow Endpoints
```
POST /follow/{userID}               - Follow a user
DELETE /follow/{userID}             - Unfollow a user
GET  /followers/{userID}            - Get list of followers (array of user IDs)
GET  /following/{userID}            - Get list of following (array of user IDs)
```

### Notification Endpoints
```
GET  /notifications                 - Get all notifications (sorted by latest first)
PUT  /notifications/{id}/read       - Mark notification as read
```

---

## 📂 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── Auth/
│   │   │   ├── LoginPage.jsx
│   │   │   ├── GoogleLoginButton.jsx
│   │   │   └── ProtectedRoute.jsx
│   │   ├── Layout/
│   │   ���   ├── Header.jsx
│   │   │   ├── Navbar.jsx
│   │   │   ├── Sidebar.jsx
│   │   │   └── Footer.jsx
│   │   ├── Posts/
│   │   │   ├── PostCard.jsx
│   │   │   ├── PostForm.jsx
│   │   │   ├── PostList.jsx
│   │   │   ├── PostDetail.jsx
│   │   │   └── EditPostModal.jsx
│   │   ├── Comments/
│   │   │   ├── CommentSection.jsx
│   │   │   ├── CommentCard.jsx
│   │   │   └── CommentForm.jsx
│   │   ├── Likes/
│   │   │   ├── LikeButton.jsx
│   │   │   └── LikesList.jsx
│   │   ├── Users/
│   │   │   ├── UserCard.jsx
│   │   │   ├── UserProfile.jsx
│   │   │   ├── UserSearch.jsx
│   │   │   ├── FollowButton.jsx
│   │   │   └── FollowersList.jsx
│   │   ├── Notifications/
│   │   │   ├── NotificationBell.jsx
│   │   │   ├── NotificationPanel.jsx
│   │   │   └── NotificationCard.jsx
│   │   ├── Feed/
│   │   │   ├── FeedPage.jsx
│   │   │   ├── TrendingSection.jsx
│   │   │   └── PersonalFeed.jsx
│   │   └── Common/
│   │       ├── LoadingSpinner.jsx
│   │       ├── ErrorAlert.jsx
│   │       ├── Modal.jsx
│   │       └── Pagination.jsx
│   ├── pages/
│   │   ├── HomePage.jsx
│   │   ├── LoginPage.jsx
│   │   ├── FeedPage.jsx
│   │   ├── UserProfilePage.jsx
│   │   ├── SearchPage.jsx
│   │   ├── TrendingPage.jsx
│   │   ├── NotificationsPage.jsx
│   │   ├── SettingsPage.jsx
│   │   └── NotFoundPage.jsx
│   ├── services/
│   │   ├── api.js
│   │   ├── authService.js
│   │   ├── userService.js
│   │   ├── postService.js
│   │   ├── commentService.js
│   │   ├── likeService.js
│   │   ├── followService.js
│   │   └── notificationService.js
│   ├── hooks/
│   │   ├── useAuth.js
│   │   ├── usePost.js
│   │   ├── useUser.js
│   │   ├── useFetch.js
│   │   ├── useNotification.js
│   │   └── useFollow.js
│   ├── context/
│   │   ├── AuthContext.js
│   │   ├── UserContext.js
│   │   └── NotificationContext.js
│   ├── utils/
│   │   ├── formatDate.js
│   │   ├── validateInput.js
│   │   ├── tokenManager.js
│   │   └── errorHandler.js
│   ├── styles/
│   │   ├── tailwind.css
│   │   ├── globals.css
│   │   └── responsive.css
│   ├── App.jsx
│   ├── index.jsx
│   └── main.jsx
├── public/
├── .env.example
├── package.json
├── tailwind.config.js
├── vite.config.js
└── README.md
```

---

## 🎨 Key Features to Implement

### 1. **Authentication System**
- Google OAuth login integration
- JWT token management (30-day expiry)
- Protected routes that redirect to login
- Auto-logout on token expiry
- Remember login state on refresh

### 2. **Post Management**
- Create posts with text content
- Edit own posts
- Delete own posts
- Real-time content validation
- Character count indicator
- Timestamp display

### 3. **Social Interactions**
- Like/Unlike posts with visual feedback (heart icon state change)
- Add comments to posts with real-time display
- Delete own comments
- Follow/Unfollow users
- View follower/following lists with user cards

### 4. **Feed Features**
- **Personal Feed**: Only posts from followed users, sorted by engagement score
- **Trending Posts**: Top 10 posts sorted by (likes × 2 + comments × 3)
- Pagination support (page, limit parameters)
- Real-time updates on action

### 5. **User Profiles**
- View any user's profile with all details
- Edit own profile (bio, location, skills, GitHub link)
- View user's all posts
- Follow/Unfollow from profile
- User statistics (followers count, following count, posts count)
- Responsive layout

### 6. **Search & Discovery**
- Search users by skills (case-insensitive)
- Filter and display results
- Skill suggestions based on database

### 7. **Notifications**
- Real-time notification display
- Mark notifications as read
- Notification types: likes, comments, follows
- Notification bell with badge count
- View full notification history

---

## 📱 Responsive Design Requirements

### Breakpoints
- **Mobile**: 320px - 640px
- **Tablet**: 641px - 1024px
- **Laptop**: 1025px and above

### Mobile (320px - 640px)
- Single column layout
- Bottom navigation bar with: Feed, Search, Trending, Profile, Notifications
- Hamburger menu for additional options
- Simplified forms (no side panels)
- Touch-friendly buttons (minimum 44px)
- Stacked cards
- Hidden desktop-only features

### Tablet (641px - 1024px)
- 2-column layout: Sidebar (fixed) + Main content
- Sidebar navigation visible but narrower
- Medium-sized cards
- Slightly optimized spacing
- Some desktop features visible

### Laptop (1025px+)
- 3-column layout: Sidebar + Main feed + Right panel
- Full sidebar with all navigation options visible
- Right panel showing trending posts and suggested users
- Wider cards with more information
- Hover effects and tooltips
- Modal dialogs for editing

---

## 🔑 Data Models & Response Formats

### User Object (from /me or /users/{userID})
```javascript
{
  id: string,
  name: string,
  email: string,
  avatar_url: string,
  bio: string,
  github: string,
  skills: string,
  location: string
}
```

### Post Object (from /posts, /feed, /trending)
```javascript
{
  id: number,
  user_id: string,
  content: string,
  created_at: string,  // ISO timestamp
  score: number        // (likes*2 + comments*3) - optional in list endpoints
}
```

### Comment Object (from /posts/{postID}/comments)
```javascript
{
  user_id: string,
  content: string,
  created_at: string  // ISO timestamp
}
```

### Notification Object (from /notifications)
```javascript
{
  id: number,
  type: 'like' | 'comment' | 'follow',
  message: string,
  created_at: string,  // ISO timestamp
  read: boolean
}
```

### Like Object (from /posts/{postID}/likes)
```javascript
// Returns array of user IDs as strings
["user_id_1", "user_id_2", "user_id_3"]
```

### Followers/Following (from /followers/{userID}, /following/{userID})
```javascript
// Returns array of user IDs
["user_id_1", "user_id_2", "user_id_3"]
```

---

## 🛠️ Implementation Guidelines

### Error Handling
- Display toast notifications for user actions (success/error)
- Show user-friendly error messages, not raw API errors
- Implement retry mechanisms for failed requests
- Handle 401 (Unauthorized) by redirecting to login
- Handle 403 (Forbidden) with permission-denied message
- Handle 404 with appropriate not-found UI
- Handle 500 with generic error message and retry button

### Performance Optimization
- Lazy load images using `<img loading="lazy">`
- Implement pagination to limit data fetching
- Debounce search input (300ms delay)
- Memoize expensive computations using useMemo
- Code splitting by route with React.lazy()
- Virtualize long lists using react-window (optional but recommended)
- Cache API responses where appropriate

### Security
- Never store sensitive data in localStorage except JWT
- Validate all inputs on frontend before sending to backend
- Sanitize user-generated content to prevent XSS
- Use HTTPS in production
- Implement CSRF tokens if needed
- Use Content Security Policy headers

### State Management Architecture
- **Global State (Context/Redux)**:
  - Authentication state (logged in user, token)
  - Notifications (unread count, list)
  - Current user data
  - Global loading/error states
- **Local State (useState)**:
  - Form inputs
  - Modal open/close states
  - Component-specific UI states
  - Temporary filters

### API Request Pattern
```javascript
// services/api.js
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:3000'
});

// Add token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token expiry
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

---

## 📄 Pages to Build (Dynamic Routing)

### 1. **HomePage** (`/`)
- Landing page visible to unauthenticated users
- Call-to-action button for Google login
- Brief introduction to DevConnect
- Redirect to `/feed` if already logged in

### 2. **LoginPage** (`/login`)
- Google OAuth login button
- Redirect to `/feed` on successful login
- Redirect to `/` if already logged in

### 3. **FeedPage** (`/feed`)
- Protected route - redirect to login if not authenticated
- Display posts from followed users
- Pagination for posts
- Create post form at top
- Show loading state while fetching

### 4. **UserProfilePage** (`/users/:userID`)
- Display user information (bio, skills, location, GitHub)
- Show user's all posts
- Follow/Unfollow button
- User statistics (followers, following, posts count)
- Responsive cards for posts

### 5. **SearchPage** (`/search?skill=javascript`)
- Search input field
- Display filtered users by skill
- Show user cards with follow button
- Real-time filtering with debounce

### 6. **TrendingPage** (`/trending`)
- Display top 10 trending posts
- Sorted by engagement score
- Show post cards with engagement metrics
- No pagination (fixed 10 posts)

### 7. **NotificationsPage** (`/notifications`)
- List all notifications
- Mark notifications as read on click
- Group by type (likes, comments, follows)
- Show timestamp
- Delete/archive notifications (optional)

### 8. **SettingsPage** (`/settings`)
- Edit profile form (name, bio, location, skills, GitHub)
- Save changes functionality
- Logout button
- Account settings section

### 9. **NotFoundPage** (`*`)
- 404 error message
- Link back to home

---

## 🎯 Component Specifications

### Layout Components

**Header/Navbar**
- Logo/branding on left
- Search bar in center (redirect to search page)
- Right side: notification bell badge, user dropdown menu
- On mobile: hamburger menu replaces search and dropdown
- Sticky positioning

**Sidebar (Desktop/Tablet)**
- Fixed width (250px on desktop, narrower on tablet)
- Navigation links: Feed, Trending, Search, Profile, Settings
- Current logged-in user card at top
- Logout button at bottom
- Hidden on mobile, shown as bottom nav

**Bottom Navigation (Mobile Only)**
- Fixed at bottom
- 5 icons: Feed, Search, Trending, Profile, Notifications
- Active state styling
- Hidden on tablet/desktop

**Main Content Area**
- Central column for primary content
- Responsive padding and margins
- Max-width constraint on desktop

### Feature Components

**PostCard**
- User avatar (clickable to profile)
- User name (clickable to profile)
- Post timestamp (relative, e.g., "2 hours ago")
- Post content (truncated with "read more" if long)
- Engagement metrics: likes count, comments count
- Like button (toggles state, shows heart icon)
- Comment button (expandable section)
- Edit button (if owner)
- Delete button (if owner) with confirmation
- Hover effects showing additional actions

**PostForm**
- Text input field for post content
- Character count (0/280)
- Submit button (disabled if empty)
- Cancel button
- Loading state on submit
- Error message on failure

**CommentSection**
- List of comments
- Each comment shows: user avatar, name, content, timestamp
- Delete button for own comments
- CommentForm at bottom

**CommentCard**
- Compact display of single comment
- User avatar and name (clickable to profile)
- Comment content
- Timestamp
- Delete button (if owner)

**CommentForm**
- Text input field
- Submit button
- Loading state

**LikeButton**
- Toggleable heart icon
- Shows like count next to it
- Visual feedback on click (animation)
- Disabled while loading

**FollowButton**
- Toggle button: Follow / Following / Unfollow
- Loading state
- Different styling for followed state
- Confirmation on unfollow (optional)

**UserCard**
- User avatar (large)
- User name (heading)
- Bio/skills preview
- Location
- Follow/Unfollow button
- Link to full profile
- Responsive sizing

**NotificationBell**
- Bell icon in header
- Badge showing unread count (red circle)
- Dropdown on click

**NotificationPanel**
- Dropdown list of last 5 notifications
- Each notification: type icon, message, timestamp
- Mark as read on click
- Link to full notifications page at bottom

**UserProfile**
- Header section: avatar, name, stats (followers, following, posts)
- Bio section with editable fields (if owner)
- Skills display
- Location
- GitHub link (clickable)
- Posts section below

**SearchForm**
- Text input with placeholder "Search by skills..."
- Debounced search
- Real-time filtering
- Clear button
- Display results as user cards

**LoadingSpinner**
- Centered spinner/skeleton
- Full-screen overlay option
- Transparent background
- Animated rotation

**ErrorAlert**
- Dismissible alert box
- Error icon
- Error message
- Retry button (if applicable)

**Modal**
- Overlay background
- Centered content box
- Close button (X)
- Responsive sizing
- Animation on open/close

**Pagination**
- Previous/Next buttons
- Page numbers (show 5 around current)
- Current page indicator
- Disabled states for first/last page
- Jump to page input (optional)

---

## 🔄 State Management Pattern Example

```javascript
// context/AuthContext.js
import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../services/api';

const AuthContext = createContext();

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check if user is already logged in on mount
  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (token) {
      fetchCurrentUser();
    } else {
      setLoading(false);
    }
  }, []);

  const fetchCurrentUser = async () => {
    try {
      const { data } = await api.get('/me');
      setUser(data);
    } catch (err) {
      localStorage.removeItem('authToken');
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, error, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// Custom hooks
export function usePost() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPosts = async (page = 1, limit = 10) => {
    setLoading(true);
    try {
      const { data } = await api.get(`/posts?page=${page}&limit=${limit}`);
      setPosts(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createPost = async (content) => {
    try {
      await api.post('/posts', { content });
      await fetchPosts(); // Refresh list
    } catch (err) {
      throw new Error(err.response?.data || 'Failed to create post');
    }
  };

  return { posts, loading, error, fetchPosts, createPost };
}
```

---

## ✅ Implementation Checklist

- [ ] Setup React project with Vite
- [ ] Install all dependencies (axios, react-router, tailwind, etc.)
- [ ] Configure Tailwind CSS and responsive utilities
- [ ] Setup API service layer with Axios interceptors
- [ ] Implement authentication context and protected routes
- [ ] Create layout components (Header, Navbar, Sidebar, Footer)
- [ ] Implement Google OAuth login flow
- [ ] Build HomePage and LoginPage
- [ ] Create post CRUD operations and PostCard component
- [ ] Build FeedPage with posts list and pagination
- [ ] Implement Like/Unlike functionality
- [ ] Create comment system (add, delete, display)
- [ ] Build Follow/Unfollow functionality
- [ ] Create UserProfilePage with editable profile
- [ ] Implement user search by skills
- [ ] Build SearchPage with results
- [ ] Create TrendingPage with top posts
- [ ] Implement notification system with bell icon
- [ ] Build NotificationsPage
- [ ] Create SettingsPage for profile edit
- [ ] Add error handling and loading states throughout
- [ ] Test all API integrations
- [ ] Make all components fully responsive (mobile, tablet, desktop)
- [ ] Add form validation
- [ ] Optimize images and performance
- [ ] Test on multiple devices/screen sizes
- [ ] Add keyboard navigation
- [ ] Deploy to production (Vercel/Netlify)

---

## 🚀 Quick Start Commands

```bash
# Create React project with Vite
npm create vite@latest devconnect-frontend -- --template react
cd devconnect-frontend

# Install dependencies
npm install axios react-router-dom tailwindcss postcss autoprefixer

# Initialize Tailwind CSS
npx tailwindcss init -p

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

---

## 📌 Critical Implementation Notes

1. **CORS**: Backend has CORS middleware enabled at port 3000
2. **Token Management**: Store JWT in localStorage, include in Authorization header
3. **Token Expiry**: 30 days - implement re-login flow
4. **Pagination**: All list endpoints support `page` and `limit` (default: page=1, limit=10)
5. **Timestamps**: Format all dates consistently (use date utility function)
6. **Rate Limiting**: Backend has rate limiter - handle 429 errors gracefully
7. **Error Messages**: Display user-friendly messages, not raw error codes
8. **Loading States**: Show spinners/skeletons while fetching data
9. **Real-time Updates**: Consider WebSocket for live notifications (future enhancement)
10. **Mobile-First**: Design for mobile first, then enhance for larger screens

---

## 🎨 Responsive Breakpoint Configuration (Tailwind)

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    screens: {
      'sm': '320px',   // Mobile
      'md': '641px',   // Tablet
      'lg': '1025px',  // Laptop/Desktop
    },
  },
}
```

---

**This master prompt gives your LLM everything needed to build a complete, production-ready DevConnect frontend! 🎉**