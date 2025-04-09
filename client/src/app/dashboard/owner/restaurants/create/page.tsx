"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/app/_components/AuthContext";

import { FullPageLoader } from "@/components/ui/loading-spinner";
import withAuth from "@/app/_components/withAuth";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { RestaurantForm } from "@/app/_components/RestaurantForm";
import {
  CreateRestaurantDTO,
  createRestaurant,
} from "@/app/_lib/restaurant-api";

function CreateRestaurantPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { owner, ownerAccessToken, refreshOwnerToken } = useAuth();
  const router = useRouter();

  const handleSubmit = async (formData: CreateRestaurantDTO) => {
    if (!ownerAccessToken) {
      setError("You must be logged in to create a restaurant");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Try to refresh the token before making the request
      const freshToken = await refreshOwnerToken();

      if (!freshToken) {
        throw new Error("Your session has expired. Please log in again.");
      }

      await createRestaurant(formData, freshToken);
      router.push("/dashboard/owner/restaurants");
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      console.error("Error creating restaurant:", err);

      if (err.response?.status === 401) {
        setError("Your session has expired. Please log in again.");
        setTimeout(() => {
          router.push("/owner-login");
        }, 2000);
      } else {
        setError(
          err.response?.data?.error ||
            "Failed to create restaurant. Please try again."
        );
      }
    } finally {
      setLoading(false);
    }
  };

  if (!ownerAccessToken) {
    return <FullPageLoader />;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center">
        <Button variant="ghost" onClick={() => router.back()} className="mr-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <h1 className="text-2xl font-bold">Create New Restaurant</h1>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="max-w-4xl mx-auto">
        <RestaurantForm
          onSubmit={handleSubmit}
          loading={loading}
          error={error}
        />
      </div>
    </div>
  );
}

export default withAuth(CreateRestaurantPage, { userType: "owner" });
