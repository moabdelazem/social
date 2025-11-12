import { useState } from "react";
import { useNavigate } from "react-router";
import type { Route } from "./+types/register";
import { Layout } from "~/components/Layout";
import { useAuth } from "~/lib/auth-context";
import { APIError } from "~/lib/api";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Register | Social" },
    { name: "description", content: "Create a new account" },
  ];
}

export default function Register() {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setIsLoading(true);

    try {
      await register(username, email, password);
      setSuccess(
        "Registration successful! Check your email for activation link."
      );
      setTimeout(() => {
        navigate("/login");
      }, 2000);
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
          <h1 className="text-2xl font-bold mb-6">Register</h1>

          {success ? (
            <div className="border-2 border-green-600 bg-green-50 p-4">
              <p className="mb-2">Account created successfully!</p>
              <p className="text-sm">Check your email for activation link.</p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block mb-2 text-sm">Username</label>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="input-field w-full"
                  placeholder="username"
                  required
                  autoComplete="username"
                />
              </div>

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
                  autoComplete="new-password"
                  minLength={8}
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
                  {isLoading ? "Loading..." : "Register"}
                </button>
                <button
                  type="button"
                  onClick={() => navigate("/login")}
                  className="btn-secondary"
                >
                  Login
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </Layout>
  );
}
