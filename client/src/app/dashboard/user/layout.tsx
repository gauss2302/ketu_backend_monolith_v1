"use client";
import React, { useState, useEffect } from "react";

import {
  LayoutDashboard,
  UserCog,
  Settings,
  LogOut,
  Menu,
  X,
} from "lucide-react";
import Image from "next/image";
import { useAuth } from "@/app/_components/AuthContext";
import Link from "next/link";

type LayoutProps = {
  children: React.ReactNode;
};

export default function UserDashboardLayout({ children }: LayoutProps) {
  const [open, setOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const { logout, user } = useAuth();

  // Detect mobile view
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);

    return () => {
      window.removeEventListener("resize", checkMobile);
    };
  }, []);

  const links = [
    {
      label: "Dashboard",
      href: "/dashboard/user",
      icon: (
        <LayoutDashboard className="text-neutral-700 dark:text-neutral-200 h-5 w-5 flex-shrink-0" />
      ),
    },
    {
      label: "Profile",
      href: "/dashboard/user/profile",
      icon: (
        <UserCog className="text-neutral-700 dark:text-neutral-200 h-5 w-5 flex-shrink-0" />
      ),
    },
    {
      label: "Settings",
      href: "/dashboard/user/settings",
      icon: (
        <Settings className="text-neutral-700 dark:text-neutral-200 h-5 w-5 flex-shrink-0" />
      ),
    },
  ];

  // Wave colors for gradient
  const waveColors = ["#38bdf8", "#818cf8", "#c084fc", "#e879f9", "#22d3ee"];
  const gradient = `linear-gradient(to bottom, ${waveColors.join(", ")})`;

  return (
    <div className="flex flex-col h-screen">
      {/* Mobile Header */}
      <header className="md:hidden flex items-center justify-between p-4 bg-white shadow-sm">
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setOpen(!open)}
            className="text-gray-500 focus:outline-none"
          >
            <Menu className="h-6 w-6" />
          </button>
          <h1 className="text-lg font-medium">My App</h1>
        </div>
        {user && (
          <div className="flex items-center space-x-2">
            <span className="text-sm">{user.name}</span>
            <div className="h-8 w-8 rounded-full bg-purple-100 flex items-center justify-center overflow-hidden">
              {user.name ? (
                <span className="text-sm font-medium">
                  {user.name.charAt(0)}
                </span>
              ) : (
                <UserCog className="h-4 w-4" />
              )}
            </div>
          </div>
        )}
      </header>

      <div className="flex flex-1 h-full overflow-hidden">
        {/* Sidebar for larger screens */}
        <div
          className={`hidden md:block h-full w-64 fixed inset-y-0 left-0 z-20`}
          style={{ background: gradient }}
        >
          <div className="flex flex-col h-full p-4">
            <div className="mb-6 flex justify-between items-center">
              <h1 className="text-white text-xl font-bold">Dashboard</h1>
            </div>

            <div className="mt-6 space-y-4 flex-1">
              {links.map((link, idx) => (
                <Link
                  key={idx}
                  href={link.href}
                  className="flex items-center text-white p-2 rounded-md hover:bg-white/10"
                >
                  {link.icon}
                  <span className="ml-3">{link.label}</span>
                </Link>
              ))}

              <button
                onClick={logout}
                className="w-full flex items-center text-white p-2 rounded-md hover:bg-white/10"
              >
                <LogOut className="h-5 w-5" />
                <span className="ml-3">Logout</span>
              </button>
            </div>

            <div className="mt-auto pt-4 border-t border-white/20">
              <div className="flex items-center space-x-3">
                <div className="h-10 w-10 rounded-full bg-white/10 flex items-center justify-center overflow-hidden">
                  {user?.name ? (
                    <span className="text-white font-medium">
                      {user.name.charAt(0)}
                    </span>
                  ) : (
                    <Image
                      src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=3164&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                      width={40}
                      height={40}
                      alt="User"
                      className="h-full w-full object-cover"
                    />
                  )}
                </div>
                <div>
                  <p className="text-sm font-medium text-white">
                    {user?.name || "User"}
                  </p>
                  <p className="text-xs text-white/70">
                    {user?.email || "user@example.com"}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Mobile Sidebar (Overlay) */}
        {open && (
          <div className="md:hidden fixed inset-0 z-30">
            <div
              className="absolute inset-0 bg-black/50"
              onClick={() => setOpen(false)}
            ></div>
            <div
              className="absolute inset-y-0 left-0 w-64 max-w-[80%] bg-white shadow-lg p-4 flex flex-col h-full"
              style={{ background: gradient }}
            >
              <div className="flex justify-between items-center mb-6">
                <h1 className="text-white text-xl font-bold">Dashboard</h1>
                <button onClick={() => setOpen(false)} className="text-white">
                  <X className="h-6 w-6" />
                </button>
              </div>

              <div className="mt-6 space-y-4 flex-1">
                {links.map((link, idx) => (
                  <Link
                    key={idx}
                    href={link.href}
                    className="flex items-center text-white p-2 rounded-md hover:bg-white/10"
                    onClick={() => setOpen(false)}
                  >
                    {link.icon}
                    <span className="ml-3">{link.label}</span>
                  </Link>
                ))}

                <button
                  onClick={() => {
                    setOpen(false);
                    logout();
                  }}
                  className="w-full flex items-center text-white p-2 rounded-md hover:bg-white/10"
                >
                  <LogOut className="h-5 w-5" />
                  <span className="ml-3">Logout</span>
                </button>
              </div>

              <div className="mt-auto pt-4 border-t border-white/20">
                <div className="flex items-center space-x-3">
                  <div className="h-10 w-10 rounded-full bg-white/10 flex items-center justify-center overflow-hidden">
                    {user?.name ? (
                      <span className="text-white font-medium">
                        {user.name.charAt(0)}
                      </span>
                    ) : (
                      <Image
                        src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=3164&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                        width={40}
                        height={40}
                        alt="User"
                        className="h-full w-full object-cover"
                      />
                    )}
                  </div>
                  <div>
                    <p className="text-sm font-medium text-white">
                      {user?.name || "User"}
                    </p>
                    <p className="text-xs text-white/70">
                      {user?.email || "user@example.com"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Main Content */}
        <main className="flex-1 p-4 overflow-auto md:ml-64">{children}</main>
      </div>
    </div>
  );
}
