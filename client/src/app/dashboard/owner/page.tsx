// app/dashboard/owner/page.tsx
"use client";

import { Button } from "@/components/ui/button";
import withAuth from "@/app/_components/withAuth";
import { useAuth } from "@/app/_components/AuthContext";

function OwnerDashboardPage() {
  const { owner, ownerLogout } = useAuth();

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
