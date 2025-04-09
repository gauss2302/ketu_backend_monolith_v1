"use client";
import React, { useState } from "react";
import {
  LayoutDashboard,
  Store,
  UserCog,
  Settings,
  LogOut,
  Menu,
  X,
} from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/app/_components/AuthContext";

type LayoutProps = {
  children: React.ReactNode;
};

export default function OwnerDashboardLayout({ children }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { owner, ownerLogout } = useAuth();

  // Links for owner dashboard
  const links = [
    {
      label: "Dashboard",
      href: "/dashboard/owner",
      icon: <LayoutDashboard className="h-5 w-5" />,
    },
    {
      label: "My Restaurants",
      href: "/dashboard/owner/restaurants",
      icon: <Store className="h-5 w-5" />,
    },
    {
      label: "Profile",
      href: "/dashboard/owner/profile",
      icon: <UserCog className="h-5 w-5" />,
    },
    {
      label: "Settings",
      href: "/dashboard/owner/settings",
      icon: <Settings className="h-5 w-5" />,
    },
  ];

  // Gradient background color
  const bgColor = "bg-gradient-to-b from-orange-500 to-amber-500";

  return (
    <div className="flex flex-col h-screen">
      {/* Mobile Header */}
      <header className="md:hidden flex items-center justify-between p-4 bg-white shadow-sm">
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setSidebarOpen(true)}
            className="text-gray-500 focus:outline-none"
          >
            <Menu className="h-6 w-6" />
          </button>
          <h1 className="text-lg font-medium">Owner Dashboard</h1>
        </div>

        {owner && (
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium">{owner.name}</span>
            <div className="h-8 w-8 rounded-full bg-amber-100 flex items-center justify-center">
              <span className="text-sm font-medium text-amber-600">
                {owner.name?.charAt(0) || "O"}
              </span>
            </div>
          </div>
        )}
      </header>

      <div className="flex flex-1 h-full overflow-hidden">
        {/* Desktop Sidebar */}
        <aside
          className={`hidden md:flex flex-col fixed inset-y-0 left-0 w-64 ${bgColor} text-white p-4 z-20`}
        >
          <div className="flex flex-col h-full">
            <div className="mb-6">
              <h1 className="text-2xl font-bold">Owner Portal</h1>
              <p className="text-sm text-white/80 mt-1">Manage your business</p>
            </div>

            <nav className="space-y-1 flex-1">
              {links.map((link, index) => (
                <Link
                  key={index}
                  href={link.href}
                  className="flex items-center px-4 py-3 rounded-md text-white hover:bg-white/10 transition-colors"
                >
                  <span className="mr-3">{link.icon}</span>
                  <span>{link.label}</span>
                </Link>
              ))}

              <button
                onClick={ownerLogout}
                className="w-full flex items-center px-4 py-3 rounded-md text-white hover:bg-white/10 transition-colors"
              >
                <LogOut className="mr-3 h-5 w-5" />
                <span>Logout</span>
              </button>
            </nav>

            {owner && (
              <div className="mt-auto pt-4 border-t border-white/20">
                <div className="flex items-center space-x-3">
                  <div className="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center">
                    <span className="font-medium text-lg">
                      {owner.name?.charAt(0) || "O"}
                    </span>
                  </div>
                  <div>
                    <p className="font-medium">{owner.name}</p>
                    <p className="text-xs text-white/80">{owner.email}</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </aside>

        {/* Mobile Sidebar (Overlay) */}
        {sidebarOpen && (
          <div className="md:hidden fixed inset-0 z-30">
            <div
              className="absolute inset-0 bg-black/50"
              onClick={() => setSidebarOpen(false)}
            ></div>

            <div
              className={`absolute inset-y-0 left-0 w-64 max-w-[80%] ${bgColor} p-4 flex flex-col`}
            >
              <div className="flex justify-between items-center mb-6">
                <h1 className="text-xl font-bold text-white">Owner Portal</h1>
                <button
                  onClick={() => setSidebarOpen(false)}
                  className="text-white p-1 rounded-full hover:bg-white/10"
                >
                  <X className="h-6 w-6" />
                </button>
              </div>

              <nav className="space-y-1 flex-1">
                {links.map((link, index) => (
                  <Link
                    key={index}
                    href={link.href}
                    className="flex items-center px-4 py-3 rounded-md text-white hover:bg-white/10 transition-colors"
                    onClick={() => setSidebarOpen(false)}
                  >
                    <span className="mr-3">{link.icon}</span>
                    <span>{link.label}</span>
                  </Link>
                ))}

                <button
                  onClick={() => {
                    setSidebarOpen(false);
                    ownerLogout();
                  }}
                  className="w-full flex items-center px-4 py-3 rounded-md text-white hover:bg-white/10 transition-colors"
                >
                  <LogOut className="mr-3 h-5 w-5" />
                  <span>Logout</span>
                </button>
              </nav>

              {owner && (
                <div className="mt-auto pt-4 border-t border-white/20">
                  <div className="flex items-center space-x-3">
                    <div className="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center">
                      <span className="font-medium text-lg">
                        {owner.name?.charAt(0) || "O"}
                      </span>
                    </div>
                    <div>
                      <p className="font-medium text-white">{owner.name}</p>
                      <p className="text-xs text-white/80">{owner.email}</p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Main Content */}
        <main className="flex-1 p-4 overflow-auto md:ml-64">{children}</main>
      </div>
    </div>
  );
}
