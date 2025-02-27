// app/dashboard/owner/page.tsx
"use client";

import { Button } from "@/components/ui/button";
import withAuth from "@/app/_components/withAuth";
import { useAuth } from "@/app/_components/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

function OwnerDashboardPage() {
  const { owner, ownerLogout, loading, ownerAccessToken } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !ownerAccessToken) {
      router.push("/owner-login");
    }
    if (!loading && !owner) {
      router.push("/owner-login");
    }
  }, [loading, owner, ownerAccessToken, router]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!owner) {
    return <div>Owner not found</div>;
  }

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Owner Dashboard</h1>
      {owner && (
        <>
          <p>Welcome, Owner {owner.name}!</p>
          <p>Owner ID: {owner.owner_id}</p>
          <p>Email: {owner.email}</p>
          <Button onClick={ownerLogout}>Logout</Button>
        </>
      )}
    </div>
  );
}
export default withAuth(OwnerDashboardPage);
