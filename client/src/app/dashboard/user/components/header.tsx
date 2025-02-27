"use client";

import { useAuth } from "@/app/_components/AuthContext";
import { Button } from "@/components/ui/button";

export function Header() {
  const { user, logout } = useAuth();

  const waveColors = ["#38bdf8", "#818cf8", "#c084fc", "#e879f9", "#22d3ee"];

  const gradient = `linear-gradient(to right, ${waveColors.join(", ")})`;

  return (
    <header
      className="p-4 flex justify-between items-center text-white"
      style={{ background: gradient }}
    >
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
