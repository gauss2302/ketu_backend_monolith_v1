// app/_components/Sidebar.tsx
"use client";
import {
  Settings,
  MapPin,
  Heart,
  Gift,
  HelpCircle,
  Bell,
  LogOut,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuth } from "./AuthContext"; // Import useAuth
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export function Sidebar() {
  const { user, logout } = useAuth();
  const menuItems = [
    { icon: Gift, label: "My Bookings", href: "/bookings" },
    { icon: Heart, label: "Saved Places", href: "/saved" },
    { icon: Bell, label: "Notifications", href: "/notifications" },
    { icon: HelpCircle, label: "Help Center", href: "/help" },
  ];

  return (
    <aside className="w-64 bg-white border-r p-4">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">User Profile</h1>
      </div>

      <div className="flex items-center space-x-4 mb-6">
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
      <div className="mb-6">
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

      <nav className="mb-6">
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

      <Button variant="outline" onClick={logout} className="w-full">
        <LogOut className="mr-2 h-4 w-4" />
        Logout
      </Button>
      <div className="mt-4">
        <Link href="/settings">
          <Button variant="ghost" className="w-full justify-start">
            <Settings className="mr-2 h-4 w-4" />
            Settings
          </Button>
        </Link>
      </div>
    </aside>
  );
}
