// app/dashboard/user/_layout/desktop.tsx
import React, { ReactNode } from "react";
import { Sidebar } from "@/app/_components/Sidebar";

type DesktopLayoutProps = {
  children: ReactNode;
};

export default function DesktopLayout({ children }: DesktopLayoutProps) {
  return (
    <div className="flex h-screen md:hidden">
      {" "}
      {/* Hidden on mobile */}
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-6">{children}</main>
    </div>
  );
}
