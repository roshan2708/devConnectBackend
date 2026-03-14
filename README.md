# devConnectBackend

Base URL: `https://devconnectbackend-wuej.onrender.com`

## API Documentation

### Authentication & Health
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/` | API root check | No |
| GET | `/health` | API health check | No |
| GET | `/auth/{provider}` | Begin OAuth (e.g., `google`) | No |
| GET | `/auth/{provider}/callback` | OAuth callback endpoint | No |

### Users & Profile
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/me` | Get current authenticated user details | Yes |
| GET | `/users/{userID}` | Get a public user profile | No |
| PUT | `/profile` | Update current user's profile | Yes |

**PUT `/profile` Payload:**
```json
{
  "bio": "string",
  "github": "string",
  "skills": "string (comma separated)",
  "location": "string"
}
```

### Posts
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/posts` | Get all posts (Pagination: `?page=1&limit=10`) | Yes |
| POST | `/posts` | Create a new post | Yes |
| GET | `/users/{userID}/posts` | Get all posts by a specific user | No |
| PUT | `/posts/{postID}` | Edit a post | Yes |
| DELETE | `/posts/{postID}` | Delete a post | Yes |

**POST & PUT `/posts...` Payload:**
```json
{
  "content": "string"
}
```

### Comments & Likes
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/posts/{postID}/comments`| Get comments on a post | No |
| POST | `/posts/{postID}/comment` | Add a comment to a post | Yes |
| DELETE | `/comments/{commentID}` | Delete a comment | Yes |
| GET | `/posts/{postID}/likes` | Get users who liked a post | No |
| POST | `/posts/{postID}/like` | Like a post | Yes |
| DELETE | `/posts/{postID}/like` | Unlike a post | Yes |

**POST `/posts/{postID}/comment` Payload:**
```json
{
  "content": "string"
}
```

### Follow System
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/followers/{userID}` | Get followers of a user | No |
| GET | `/following/{userID}` | Get users followed by a user | No |
| POST | `/follow/{userID}` | Follow a user | Yes |
| DELETE | `/follow/{userID}` | Unfollow a user | Yes |

### Feed & Discover
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/feed` | Get personalized feed of posts from followed users | Yes |
| GET | `/trending` | Get trending posts based on engagement score | No |
| GET | `/search` | Search for users (`?query=string`) | No |

### Notifications
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/notifications` | Get user's notifications | Yes |
| PUT | `/notifications/{id}/read`| Mark a notification as read | Yes |
