"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/app/_components/AuthContext";

import { FullPageLoader } from "@/components/ui/loading-spinner";
import withAuth from "@/app/_components/withAuth";
import { Button } from "@/components/ui/button";
import { Plus, ExternalLink, Edit, Trash2 } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Image from "next/image";
import { DashboardSkeleton } from "@/components/ui/skeleton";
import {
  getOwnerRestaurants,
  RestaurantResponse,
} from "@/app/_lib/restaurant-api";

function RestaurantsPage() {
  const [restaurants, setRestaurants] = useState<RestaurantResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { ownerAccessToken } = useAuth();
  const router = useRouter();

  useEffect(() => {
    const fetchRestaurants = async () => {
      if (!ownerAccessToken) return;

      try {
        setLoading(true);
        const data = await getOwnerRestaurants(ownerAccessToken);
        setRestaurants(data);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        console.error("Error fetching restaurants:", err);
        setError("Failed to load restaurants. Please try again.");
      } finally {
        setLoading(false);
      }
    };

    fetchRestaurants();
  }, [ownerAccessToken]);

  if (!ownerAccessToken) {
    return <FullPageLoader />;
  }

  if (loading) {
    return <DashboardSkeleton />;
  }

  const handleCreateRestaurant = () => {
    router.push("/dashboard/owner/restaurants/create");
  };

  const handleEditRestaurant = (id: number) => {
    router.push(`/dashboard/owner/restaurants/edit/${id}`);
  };

  const handleViewRestaurant = (id: number) => {
    router.push(`/dashboard/owner/restaurants/${id}`);
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">My Restaurants</h1>
        <Button onClick={handleCreateRestaurant}>
          <Plus className="h-4 w-4 mr-2" />
          Create New Restaurant
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 p-4 rounded-md mb-4">
          {error}
        </div>
      )}

      {!loading && restaurants.length === 0 ? (
        <Card className="bg-gray-50">
          <CardContent className="flex flex-col items-center justify-center p-6">
            <div className="rounded-full bg-gray-100 p-3">
              <Plus className="h-6 w-6 text-gray-400" />
            </div>
            <h3 className="mt-4 text-lg font-medium">No restaurants yet</h3>
            <p className="mt-1 text-center text-sm text-gray-500">
              Get started by creating your first restaurant
            </p>
            <Button onClick={handleCreateRestaurant} className="mt-4">
              Create Restaurant
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {restaurants.map((restaurant) => (
            <Card key={restaurant.id} className="overflow-hidden">
              <div className="relative h-48 w-full">
                <Image
                  src={
                    restaurant.main_image ||
                    "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?q=80&w=2070&auto=format&fit=crop"
                  }
                  alt={restaurant.name}
                  fill
                  className="object-cover"
                />
              </div>
              <CardHeader>
                <CardTitle>{restaurant.name}</CardTitle>
                <CardDescription>
                  {restaurant.location.address.city},{" "}
                  {restaurant.location.address.district}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="line-clamp-3 text-sm text-gray-500">
                  {restaurant.description}
                </p>
                <div className="mt-4 grid grid-cols-2 gap-2 text-xs">
                  <div className="bg-amber-50 text-amber-600 px-2 py-1 rounded">
                    Capacity: {restaurant.details.capacity} seats
                  </div>
                  <div className="bg-blue-50 text-blue-600 px-2 py-1 rounded">
                    Hours: {restaurant.details.opening_hours}
                  </div>
                </div>
              </CardContent>
              <CardFooter className="flex justify-between">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleViewRestaurant(restaurant.id)}
                >
                  <ExternalLink className="h-4 w-4 mr-2" />
                  View
                </Button>
                <div className="space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleEditRestaurant(restaurant.id)}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="text-red-500 border-red-200 hover:bg-red-50"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export default withAuth(RestaurantsPage, { userType: "owner" });
