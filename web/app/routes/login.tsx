import { useState } from "react";
import { useNavigate } from "react-router";
import type { Route } from "./+types/login";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";
import { APIError } from "~/lib/api";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Login | Social" },
    { name: "description", content: "Login to your account" },
  ];
}

export default function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      await login(email, password);
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
      <div className="max-w-md mx-auto mt-12">
        <div className="border-2 border-black p-8">
          <h1 className="text-2xl font-bold mb-6">Login</h1>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block mb-2 text-sm">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field w-full"
                placeholder="user@example.com"
                required
                autoComplete="email"
              />
            </div>

            <div>
              <label className="block mb-2 text-sm">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input-field w-full"
                placeholder="••••••••"
                required
                autoComplete="current-password"
              />
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
                {isLoading ? "Loading..." : "Login"}
              </button>
              <button
                type="button"
                onClick={() => navigate("/register")}
                className="btn-secondary"
              >
                Register
              </button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
}
