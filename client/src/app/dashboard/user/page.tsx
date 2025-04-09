// app/dashboard/user/page.tsx
"use client";
import { useAuth } from "@/app/_components/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { FullPageLoader } from "@/components/ui/loading-spinner";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { LogOut } from "lucide-react";

export default function UserDashboard() {
  const { user, loading, isAuthenticated, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated()) {
      router.push("/login");
    }
  }, [loading, isAuthenticated, router, user]);

  if (loading) {
    return <FullPageLoader />;
  }

  if (!user) {
    return null;
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Welcome, {user.name}!</CardTitle>
          <CardDescription>This is your personal dashboard.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h3 className="text-lg font-medium mb-2">Account Information</h3>
              <div className="space-y-1 text-sm">
                <p>
                  <span className="font-medium">User ID:</span> {user.id}
                </p>
                <p>
                  <span className="font-medium">Username:</span>{" "}
                  {user.username || "Not set"}
                </p>
                <p>
                  <span className="font-medium">Email:</span> {user.email}
                </p>
              </div>
            </div>
            <div>
              <h3 className="text-lg font-medium mb-2">Quick Actions</h3>
              <div className="space-y-2">
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => router.push("/dashboard/user/profile")}
                >
                  Edit Profile
                </Button>
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => router.push("/dashboard/user/settings")}
                >
                  Account Settings
                </Button>
                <Button
                  variant="outline"
                  className="w-full justify-start text-red-500"
                  onClick={logout}
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
        <Card>
          <CardHeader>
            <CardTitle>Favorites</CardTitle>
            <CardDescription>Your saved places</CardDescription>
          </CardHeader>
          <CardContent className="h-40 flex items-center justify-center">
            <p className="text-center text-sm text-muted-foreground">
              You haven't saved any places yet.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Your recent browsing history</CardDescription>
          </CardHeader>
          <CardContent className="h-40 flex items-center justify-center">
            <p className="text-center text-sm text-muted-foreground">
              No recent activity to display.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Notifications</CardTitle>
            <CardDescription>Your latest alerts</CardDescription>
          </CardHeader>
          <CardContent className="h-40 flex items-center justify-center">
            <p className="text-center text-sm text-muted-foreground">
              You have no new notifications.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
