import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router";
import type { Route } from "./+types/posts";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";
import { postsAPI, usersAPI } from "~/lib/api";
import type { Post } from "~/lib/api";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "My Posts | Social" },
    { name: "description", content: "Your posts" },
  ];
}

export default function MyPosts() {
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

    const fetchPosts = async () => {
      try {
        const data = await usersAPI.getMyPosts();
        setPosts(data);
      } catch (err) {
        setError("Failed to load posts");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchPosts();
  }, [user, navigate]);
  const handleDelete = async (postId: number) => {
    if (!confirm("Are you sure you want to delete this post?")) {
      return;
    }

    try {
      await postsAPI.delete(postId);
      setPosts(posts.filter((p) => p.id !== postId));
    } catch (err) {
      alert("Failed to delete post");
      console.error(err);
    }
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="border-2 border-black p-6">
          <p>Loading posts...</p>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold mb-6">My Posts</h1>
          <Link to="/create" className="btn-secondary no-underline text-sm">
            Create New Post
          </Link>
        </div>

        {error && (
          <div className="text-error border-2 border-red-600 bg-red-50 p-4">
            {error}
          </div>
        )}

        {posts.length === 0 ? (
          <div className="border-2 border-black p-12 text-center">
            <p className="text-gray-600 mb-4">
              You haven't created any posts yet.
            </p>
            <Link
              to="/create"
              className="btn-primary no-underline inline-block"
            >
              Create Your First Post
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <article key={post.id} className="border-2 border-black p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <h2 className="text-xl font-bold mb-2">{post.title}</h2>
                    <p className="text-sm text-gray-600">
                      Posted on {new Date(post.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>

                <p className="mb-4 whitespace-pre-wrap">{post.content}</p>

                {post.tags && post.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2 mb-4">
                    {post.tags.map((tag, idx) => (
                      <span
                        key={idx}
                        className="text-xs border border-black px-2 py-1"
                      >
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <div className="flex items-center gap-4 text-sm border-t border-black pt-4">
                  <button
                    onClick={() => handleDelete(post.id)}
                    className="underline hover:no-underline"
                  >
                    Delete
                  </button>
                </div>
              </article>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}
