// app/page.tsx
"use client";
import { useAuth } from "@/app/_components/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const { user, owner, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      if (user) {
        router.push("/dashboard/user");
      } else if (owner) {
        router.push("/dashboard/owner");
      } else {
        router.push("/login"); // Redirect to login if not authenticated
      }
    }
  }, [user, owner, loading, router]);

  if (loading) {
    return <div>Loading...</div>; // Show loading indicator
  }

  return null; // Or a brief loading message/spinner. We redirect immediately.
}
