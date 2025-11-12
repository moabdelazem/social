import { useState } from "react";
import { useNavigate } from "react-router";
import type { Route } from "./+types/create";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";
import { postsAPI, APIError } from "~/lib/api";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Create Post | Social" },
    { name: "description", content: "Create a new post" },
  ];
}

export default function CreatePost() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  if (!user) {
    navigate("/login");
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const tagArray = tags
        .split(",")
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0);

      await postsAPI.create(
        title,
        content,
        tagArray.length > 0 ? tagArray : undefined
      );
      navigate("/feed");
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-2xl mx-auto mt-8">
        <div className="border-2 border-black p-8">
          <h1 className="text-2xl font-bold mb-6">Create Post</h1>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block mb-2 text-sm">Title</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="input-field w-full"
                placeholder="My awesome post"
                required
                maxLength={100}
              />
              <p className="text-xs text-gray-500 mt-1">
                {title.length}/100 characters
              </p>
            </div>

            <div>
              <label className="block mb-2 text-sm">Content</label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="input-field w-full min-h-[200px] resize-y"
                placeholder="Write your post content here..."
                required
                maxLength={1000}
              />
              <p className="text-xs text-gray-500 mt-1">
                {content.length}/1000 characters
              </p>
            </div>

            <div>
              <label className="block mb-2 text-sm">
                Tags (optional, comma-separated)
              </label>
              <input
                type="text"
                value={tags}
                onChange={(e) => setTags(e.target.value)}
                className="input-field w-full"
                placeholder="golang, react, tech"
              />
              <p className="text-xs text-gray-500 mt-1">
                Example: golang, react, programming
              </p>
            </div>

            {error && (
              <div className="text-error border-2 border-red-600 bg-red-50 p-4">
                {error}
              </div>
            )}

            <div className="flex gap-4 pt-4">
              <button
                type="submit"
                disabled={isLoading}
                className="btn-primary flex-1"
              >
                {isLoading ? "Publishing..." : "Publish Post"}
              </button>
              <button
                type="button"
                onClick={() => navigate("/feed")}
                className="btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
}
