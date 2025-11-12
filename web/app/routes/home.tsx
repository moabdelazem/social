import { Link } from "react-router";
import type { Route } from "./+types/home";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Social | A Simple Social Network" },
    { name: "description", content: "A minimalist social network" },
  ];
}

export default function Home() {
  const { user } = useAuth();

  return (
    <Layout>
      <div className="space-y-8">
        {/* Hero Section */}
        <div className="border-2 border-black p-8">
          <h2 className="text-3xl font-bold mb-4">Welcome to Social</h2>
          <p className="text-lg mb-6">
            A minimalist social network for sharing thoughts and connecting with
            others.
          </p>

          {user ? (
            <div className="space-y-4">
              <p className="text-gray-600">
                Welcome back, <strong>{user.username}</strong>!
              </p>
              <div className="flex gap-4">
                <Link to="/feed" className="btn-primary no-underline">
                  View Feed
                </Link>
                <Link to="/create" className="btn-secondary no-underline">
                  Create Post
                </Link>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <p className="text-gray-600 mb-4">
                Please log in or register to continue.
              </p>
              <div className="flex gap-4">
                <Link to="/login" className="btn-primary no-underline">
                  Login
                </Link>
                <Link to="/register" className="btn-secondary no-underline">
                  Register
                </Link>
              </div>
            </div>
          )}
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-3 gap-6">
          <div className="border border-black p-6">
            <h3 className="text-xl font-bold mb-2">Posts</h3>
            <p className="text-sm text-gray-600">
              Share your thoughts with your followers
            </p>
          </div>
          <div className="border border-black p-6">
            <h3 className="text-xl font-bold mb-2">Follow</h3>
            <p className="text-sm text-gray-600">Connect with other users</p>
          </div>
          <div className="border border-black p-6">
            <h3 className="text-xl font-bold mb-2">Feed</h3>
            <p className="text-sm text-gray-600">Personalized content stream</p>
          </div>
        </div>

        {/* About */}
        <div className="border border-black p-6">
          <h3 className="text-xl font-bold mb-3">About</h3>
          <ul className="list-disc list-inside space-y-2 text-sm text-gray-600">
            <li>
              API Endpoint: <code>http://localhost:6767/v1</code>
            </li>
            <li>Authentication: JWT Bearer Tokens</li>
            <li>Database: PostgreSQL 16</li>
            <li>Features: Posts, Comments, Followers, Feed</li>
          </ul>
        </div>
      </div>
    </Layout>
  );
}
