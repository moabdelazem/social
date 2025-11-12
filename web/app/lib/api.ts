// API Configuration
const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:6767/v1";

// API Error class
export class APIError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: any
  ) {
    super(message);
    this.name = "APIError";
  }
}

// Generic API request function
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem("token");

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  // Merge additional headers
  if (options.headers) {
    Object.entries(options.headers).forEach(([key, value]) => {
      headers[key] = String(value);
    });
  }

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new APIError(
      data.error || "An error occurred",
      response.status,
      data
    );
  }

  return data.data || data;
}

// Auth API
export const authAPI = {
  register: async (username: string, email: string, password: string) => {
    return apiRequest<{ user: User; token: string }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    });
  },

  login: async (email: string, password: string) => {
    return apiRequest<{ user: User; token: string }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  activate: async (token: string) => {
    return apiRequest<{ message: string }>("/auth/activate", {
      method: "PUT",
      body: JSON.stringify({ token }),
    });
  },
};

// Posts API
export const postsAPI = {
  create: async (title: string, content: string, tags?: string[]) => {
    return apiRequest<Post>("/posts", {
      method: "POST",
      body: JSON.stringify({ title, content, tags }),
    });
  },

  getById: async (id: number) => {
    return apiRequest<Post>(`/posts/${id}`);
  },

  update: async (id: number, data: { title?: string; content?: string }) => {
    return apiRequest<Post>(`/posts/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete: async (id: number) => {
    return apiRequest<void>(`/posts/${id}`, {
      method: "DELETE",
    });
  },
};

// Users API
export const usersAPI = {
  getById: async (id: number) => {
    return apiRequest<User>(`/users/${id}`);
  },

  follow: async (id: number) => {
    return apiRequest<void>(`/users/${id}/follow`, {
      method: "PUT",
      body: JSON.stringify({ user_id: id }),
    });
  },

  unfollow: async (id: number) => {
    return apiRequest<void>(`/users/${id}/unfollow`, {
      method: "PUT",
    });
  },

  getFeed: async (limit = 20, offset = 0, sort = "desc") => {
    return apiRequest<Post[]>(
      `/users/feed?limit=${limit}&offset=${offset}&sort=${sort}`
    );
  },

  getMyPosts: async () => {
    return apiRequest<Post[]>("/users/me/posts");
  },
};

// Types
export interface User {
  id: number;
  username: string;
  email: string;
  is_active: boolean;
  created_at: string;
}

export interface Post {
  id: number;
  user_id: number;
  title: string;
  content: string;
  tags: string[];
  version: number;
  created_at: string;
  updated_at: string;
  comments?: Comment[];
  comments_count?: number;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  created_at: string;
  user?: User;
}
