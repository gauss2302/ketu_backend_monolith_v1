// app/dashboard/user/page.tsx
"use client";
import { useAuth } from "@/app/_components/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserDashboard() {
  const { user, loading, isAuthenticated, logout } = useAuth(); // Get logout function
  const router = useRouter();

  useEffect(() => {
    console.log("UserDashboard useEffect running");
    console.log("User:", user);
    console.log("Loading:", loading);
    console.log("isAuthenticated:", isAuthenticated());
    if (!loading && !isAuthenticated()) {
      router.push("/login");
    }
  }, [loading, isAuthenticated, router, user]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return null; // Or a loading indicator/placeholder
  }

  return (
    <div>
      <h1>User Dashboard</h1>
      <p>Welcome, {user.name}!</p>
      <p>This is a test string</p>
      <div>Test Dashboard Content</div>
      <button onClick={logout}>Logout</button> {/* Add logout button */}
    </div>
  );
}
