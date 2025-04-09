// app/dashboard/owner/page.tsx
"use client";

import { Button } from "@/components/ui/button";
import withAuth from "@/app/_components/withAuth";
import { useAuth } from "@/app/_components/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { FullPageLoader } from "@/components/ui/loading-spinner";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Store,
  LogOut,
  Plus,
  BarChart3,
  Calendar,
  Settings,
} from "lucide-react";

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
    return <FullPageLoader />;
  }

  if (!owner) {
    return null;
  }

  const stats = [
    {
      title: "Restaurants",
      value: "0",
      description: "Total restaurants managed",
      icon: <Store className="h-5 w-5 text-amber-500" />,
      action: "View All",
      onClick: () => router.push("/dashboard/owner/restaurants"),
    },
    {
      title: "Bookings",
      value: "0",
      description: "Total bookings received",
      icon: <Calendar className="h-5 w-5 text-emerald-500" />,
      action: "Manage",
      onClick: () => router.push("/dashboard/owner/bookings"),
    },
    {
      title: "Revenue",
      value: "$0",
      description: "Total revenue this month",
      icon: <BarChart3 className="h-5 w-5 text-blue-500" />,
      action: "Analytics",
      onClick: () => router.push("/dashboard/owner/analytics"),
    },
  ];

  return (
    <div className="space-y-6">
      <Card className="border-amber-200 bg-gradient-to-r from-amber-50 to-amber-100">
        <CardHeader>
          <CardTitle className="text-2xl">Welcome, {owner.name}!</CardTitle>
          <CardDescription className="text-amber-700">
            Manage your restaurants and monitor your business
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h3 className="text-lg font-medium mb-2">Account Information</h3>
              <div className="space-y-1 text-sm">
                <p>
                  <span className="font-medium">Owner ID:</span>{" "}
                  {owner.owner_id}
                </p>
                <p>
                  <span className="font-medium">Name:</span> {owner.name}
                </p>
                <p>
                  <span className="font-medium">Email:</span> {owner.email}
                </p>
                <p>
                  <span className="font-medium">Phone:</span> {owner.phone}
                </p>
                <p>
                  <span className="font-medium">Created:</span>{" "}
                  {new Date(owner.created_at).toLocaleDateString()}
                </p>
              </div>
            </div>
            <div>
              <h3 className="text-lg font-medium mb-2">Quick Actions</h3>
              <div className="space-y-2">
                <Button
                  className="w-full justify-start bg-amber-500 hover:bg-amber-600"
                  onClick={() =>
                    router.push("/dashboard/owner/restaurants/create")
                  }
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Add New Restaurant
                </Button>
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => router.push("/dashboard/owner/profile")}
                >
                  <Settings className="w-4 h-4 mr-2" />
                  Account Settings
                </Button>
                <Button
                  variant="outline"
                  className="w-full justify-start text-red-500"
                  onClick={ownerLogout}
                >
                  <LogOut className="w-4 h-4 mr-2" />
                  Logout
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {stats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              {stat.icon}
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground mt-1">
                {stat.description}
              </p>
            </CardContent>
            <CardFooter>
              <Button
                variant="ghost"
                className="w-full justify-center text-sm h-8"
                onClick={stat.onClick}
              >
                {stat.action}
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Getting Started</CardTitle>
          <CardDescription>
            Follow these steps to set up your restaurant business
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ol className="space-y-4 list-decimal list-inside text-sm">
            <li className="p-3 bg-gray-50 rounded-md">
              <span className="font-medium">Create your first restaurant</span>
              <p className="pl-5 mt-1 text-muted-foreground">
                Add details about your restaurant including description,
                location, and menu.
              </p>
            </li>
            <li className="p-3 bg-gray-50 rounded-md">
              <span className="font-medium">Set up your menu</span>
              <p className="pl-5 mt-1 text-muted-foreground">
                Add categories and items to your restaurant menu.
              </p>
            </li>
            <li className="p-3 bg-gray-50 rounded-md">
              <span className="font-medium">Configure opening hours</span>
              <p className="pl-5 mt-1 text-muted-foreground">
                Set your restaurant's opening and closing times.
              </p>
            </li>
            <li className="p-3 bg-gray-50 rounded-md">
              <span className="font-medium">Add photos of your restaurant</span>
              <p className="pl-5 mt-1 text-muted-foreground">
                Upload high-quality images to showcase your restaurant.
              </p>
            </li>
          </ol>
        </CardContent>
      </Card>
    </div>
  );
}
export default withAuth(OwnerDashboardPage);
