import { Link } from "react-router";
import type { ReactNode } from "react";
import { useAuth } from "~/lib/auth-context";

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-white text-black">
      {/* Header */}
      <header className="border-b-2 border-black py-6">
        <div className="max-w-4xl mx-auto px-4">
          <div className="flex items-center justify-between mb-6">
            <Link to="/" className="no-underline">
              <h1 className="text-2xl font-bold">SOCIAL</h1>
            </Link>

            <div className="flex items-center gap-6">
              {user ? (
                <>
                  <span className="text-sm">{user.username}</span>
                  <button
                    onClick={logout}
                    className="text-sm underline hover:no-underline"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link to="/login" className="text-sm">
                    Login
                  </Link>
                  <Link to="/register" className="text-sm">
                    Register
                  </Link>
                </>
              )}
            </div>
          </div>

          {user && (
            <nav className="flex gap-6 text-sm">
              <Link to="/feed">Feed</Link>
              <Link to="/posts">My Posts</Link>
              <Link to="/create">New Post</Link>
            </nav>
          )}
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 py-8">{children}</main>
    </div>
  );
}
