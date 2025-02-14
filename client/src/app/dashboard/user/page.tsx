// app/dashboard/user/page.tsx
"use client";
import { useAuth } from "@/app/_components/AuthContext";

import { Button } from "@/components/ui/button";
import withAuth from "@/app/_components/withAuth";
import {
  Settings,
  ChevronRight,
  MapPin,
  Star,
  Heart,
  Gift,
  HelpCircle,
  Bell,
  LogOut,
} from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Link from "next/link";

function UserDashboardPage() {
  const { user, logout } = useAuth();

  const menuItems = [
    { icon: Gift, label: "My Bookings", badge: "2 Active", href: "/bookings" },
    { icon: Heart, label: "Saved Places", badge: "12 Places", href: "/saved" },
    {
      icon: Bell,
      label: "Notifications",
      badge: "3 New",
      href: "/notifications",
    },
    { icon: HelpCircle, label: "Help Center", href: "/help" },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Main Content Area */}
      <div className="lg:pl-64">
        {" "}
        {/* Add padding-left on large screens for sidebar space */}
        {/* Profile Header */}
        <div className="bg-white px-4 pt-6 pb-4">
          <div className="flex justify-between items-start mb-6">
            <h1 className="text-2xl font-bold">Profile</h1>
            <Link href="/settings">
              <Button variant="ghost" size="icon">
                <Settings className="w-6 h-6 text-gray-600" />
              </Button>
            </Link>
          </div>

          <div className="flex items-center space-x-4">
            <Avatar>
              <AvatarImage src="" />
              <AvatarFallback>{user?.name?.charAt(0) || "U"}</AvatarFallback>
            </Avatar>
            <div>
              <h2 className="font-semibold text-lg">{user?.name || "User"}</h2>
              <div className="flex items-center text-gray-500 text-sm mt-1">
                <MapPin className="w-4 h-4 mr-1" />
                Tashkent, Uzbekistan
              </div>
            </div>
          </div>
        </div>
        {/* Stats */}
        <Card className="mt-2">
          <CardContent className="grid grid-cols-3 gap-4 p-4">
            <div className="text-center">
              <div
                className="text-2xl font-semibold"
                style={{ color: "hsl(var(--primary-app))" }}
              >
                12
              </div>
              <div className="text-sm text-gray-500">Places Visited</div>
            </div>
            <div className="text-center border-l border-r">
              <div
                className="text-2xl font-semibold"
                style={{ color: "hsl(var(--primary-app))" }}
              >
                43
              </div>
              <div className="text-sm text-gray-500">Reviews</div>
            </div>
            <div className="text-center">
              <div
                className="text-2xl font-semibold"
                style={{ color: "hsl(var(--primary-app))" }}
              >
                4.8
              </div>
              <div className="text-sm text-gray-500">Rating</div>
            </div>
          </CardContent>
        </Card>
        {/* Membership Card */}
        <Card className="mt-4 mx-4">
          <CardHeader className="bg-gradient-to-r from-primary-app to-primary-app-light rounded-t-lg text-white p-4">
            <div className="flex justify-between items-start mb-4">
              <div>
                <CardTitle className="text-sm opacity-90 mb-1">
                  Membership
                </CardTitle>
                <div className="font-semibold text-base">Premium Member</div>
              </div>
              <Star className="w-6 h-6 fill-current" />
            </div>
            <div className="flex justify-between items-end">
              <div>
                <div className="text-sm opacity-90 mb-1">Valid Until</div>
                <div className="font-semibold text-base">Dec 2024</div>
              </div>
              <div className="text-sm opacity-90">
                ID: {user?.id || "12345"}
              </div>
            </div>
          </CardHeader>
        </Card>
        {/* Menu Items */}
        <Card className="mt-4">
          <CardContent className="p-0">
            {menuItems.map((item, index) => (
              <Link href={item.href} key={index}>
                <Button
                  variant="ghost"
                  className="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50 rounded-none"
                >
                  <div className="flex items-center space-x-3">
                    <div className="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center">
                      <item.icon
                        className="w-5 h-5"
                        style={{ color: "hsl(var(--primary-app))" }}
                      />
                    </div>
                    <span className="font-medium">{item.label}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    {item.badge && (
                      <span className="text-sm text-gray-500">
                        {item.badge}
                      </span>
                    )}
                    <ChevronRight className="w-5 h-5 text-gray-400" />
                  </div>
                </Button>
              </Link>
            ))}
          </CardContent>
        </Card>
        {/* Logout Button */}
        <div className="px-4 mt-4">
          <Button
            variant="outline"
            className="w-full flex items-center justify-center space-x-2 px-4 py-3 text-red-500 bg-red-50 rounded-xl border-red-500"
            onClick={logout}
          >
            <LogOut className="w-5 h-5" />
            <span className="font-medium">Log Out</span>
          </Button>
        </div>
      </div>

      {/* Sidebar (Hidden on Small Screens) */}
      <aside className="fixed top-0 left-0 w-64 h-screen bg-white border-r transition-transform transform -translate-x-full lg:translate-x-0">
        <div className="mb-6 p-4">
          <h1 className="text-2xl font-bold">User Profile</h1>
        </div>

        <div className="flex items-center space-x-4 mb-6 px-4">
          <Avatar>
            <AvatarImage src="" />
            <AvatarFallback>{user?.name?.charAt(0) || "U"}</AvatarFallback>
          </Avatar>
          <div>
            <h2 className="font-semibold text-lg">{user?.name || "User"}</h2>
            <div className="flex items-center text-gray-500 text-sm mt-1">
              <MapPin className="w-4 h-4 mr-1" />
              Tashkent, Uzbekistan
            </div>
          </div>
        </div>
        <div className="mb-6 px-4">
          <div
            className="text-sm  mb-1"
            style={{ color: "hsl(var(--primary-app))" }}
          >
            Membership
          </div>
          <div className="font-semibold">Premium Member</div>
          <div className="text-sm mt-1">Valid Until: Dec 2024</div>
          <div className="text-sm ">ID: {user?.id || "12345"}</div>
        </div>

        <nav className="mb-6 px-4">
          <ul>
            {menuItems.map((item, index) => (
              <li key={index} className="mb-2">
                <Link href={item.href}>
                  <Button variant="ghost" className="w-full justify-start">
                    <item.icon className="mr-2 h-4 w-4" />
                    {item.label}
                  </Button>
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <div className="absolute bottom-4 left-4 right-4 px-4">
          <Button variant="outline" onClick={logout} className="w-full mb-2">
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
          <Link href="/settings">
            <Button variant="ghost" className="w-full justify-start">
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Button>
          </Link>
        </div>
      </aside>
    </div>
  );
}

export default withAuth(UserDashboardPage, { userType: "user" });
