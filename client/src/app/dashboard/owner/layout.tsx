// app/dashboard/user/layout.tsx
import React, { ReactNode } from "react";

type LayoutProps = {
  children: ReactNode;
};

export default function OwnerDashboardLayout({ children }: LayoutProps) {
  return (
    <div className="flex flex-col min-h-screen">
      <main className="flex-1 p-4">{children}</main>
    </div>
  );
}
