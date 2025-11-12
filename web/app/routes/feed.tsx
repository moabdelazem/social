import { useState, useEffect } from "react";
import type { Route } from "./+types/feed";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";
import { usersAPI } from "~/lib/api";
import type { Post } from "~/lib/api";
import { useNavigate } from "react-router";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Feed | Social" },
    { name: "description", content: "Your social feed" },
  ];
}

export default function Feed() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!user) {
      navigate("/login");
      return;
    }

    const fetchFeed = async () => {
      try {
        const data = await usersAPI.getFeed(20, 0, "desc");
        setPosts(data);
      } catch (err) {
        setError("Failed to load feed");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchFeed();
  }, [user, navigate]);

  if (isLoading) {
    return (
      <Layout>
        <div className="border-2 border-black p-6">
          <p>Loading feed...</p>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold mb-6">Feed</h1>
          <button
            onClick={() => navigate("/create")}
            className="btn-primary text-sm"
          >
            New Post
          </button>
        </div>

        {error && (
          <div className="text-error border-2 border-red-600 bg-red-50 p-4">
            {error}
          </div>
        )}

        {posts.length === 0 ? (
          <div className="border-2 border-black p-6 text-center py-12">
            <p className="mb-4">No posts in your feed</p>
            <p className="text-sm text-gray-600">
              Start following users or create your first post
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <article
                key={post.id}
                className="border-2 border-black p-6 hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h2 className="text-lg font-bold mb-1">{post.title}</h2>
                    <p className="text-xs text-gray-500">
                      ID: {post.id} | User: {post.user_id}
                    </p>
                  </div>
                  <span className="text-xs text-gray-500">
                    {new Date(post.created_at).toLocaleDateString()}
                  </span>
                </div>

                <p className="mb-3 whitespace-pre-wrap">{post.content}</p>

                {post.tags && post.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2 mb-3">
                    {post.tags.map((tag, idx) => (
                      <span
                        key={idx}
                        className="text-xs border-2 border-black px-2 py-1"
                      >
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <div className="flex items-center gap-4 text-sm border-t-2 border-black pt-3">
                  <span className="text-gray-600">
                    Comments: {post.comments_count || 0}
                  </span>
                </div>
              </article>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}
