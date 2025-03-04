"use client";
import React, { useState } from "react";
import { Sidebar, SidebarBody, SidebarLink } from "@/components/ui/sidebar";
import { LayoutDashboard, UserCog, Settings, LogOut } from "lucide-react";
import Image from "next/image";
import { Header } from "./components/header";
// import { Footer } from "./components/footer";
import { useAuth } from "@/app/_components/AuthContext";

type LayoutProps = {
  children: React.ReactNode;
};

export default function UserDashboardLayout({ children }: LayoutProps) {
  const [open, setOpen] = useState(false);
  const { logout } = useAuth();

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

  return (
    <div className="flex flex-col h-screen">
      {/* <Header /> */}
      <div className="flex flex-1">
        <Sidebar open={open} setOpen={setOpen}>
          <SidebarBody className="justify-between gap-10 h-full">
            <div className="flex flex-col flex-1 overflow-y-auto overflow-x-hidden h-full">
              <div className="mt-8 flex flex-col gap-2">
                {links.map((link, idx) => (
                  <SidebarLink key={idx} link={link} />
                ))}
                <button
                  onClick={logout}
                  className="flex items-center justify-start gap-2 group/sidebar py-2"
                >
                  <LogOut className="text-neutral-700 dark:text-neutral-200 h-5 w-5 flex-shrink-0" />
                  <span className="text-neutral-700 dark:text-neutral-200 text-sm group-hover/sidebar:translate-x-1 transition duration-150 whitespace-pre inline-block !p-0 !m-0">
                    Logout
                  </span>
                </button>
              </div>
            </div>
            <div>
              <SidebarLink
                link={{
                  label: "User Profile",
                  href: "/dashboard/user/profile",
                  icon: (
                    <Image
                      src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=3164&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                      className="h-7 w-7 flex-shrink-0 rounded-full"
                      width={50}
                      height={50}
                      alt="Avatar"
                    />
                  ),
                }}
              />
            </div>
          </SidebarBody>
        </Sidebar>
        <main className="flex-1 p-4">
          <Header />
          {children}
        </main>
      </div>
      {/* <Footer /> */}
    </div>
  );
}
