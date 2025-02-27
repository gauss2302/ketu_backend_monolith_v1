"use client";

import { useAuth } from "@/app/_components/AuthContext";
import { Button } from "@/components/ui/button";

export function Header() {
  const { user, logout } = useAuth();

  return (
    <header className="bg-gray-200 p-4 flex justify-between items-center">
      <h1 className="text-xl font-semibold">My App</h1>
      {user && (
        <div className="flex items-center space-x-4">
          <span>Welcome, {user.name}</span>
          <Button variant="default" size="default" onClick={logout}>
            Logout
          </Button>
        </div>
      )}
    </header>
  );
}
