// app/dashboard/user/layout.tsx
import React, { ReactNode } from "react";

import { Header } from "./components/header";
import { Footer } from "./components/footer";

type LayoutProps = {
  children: ReactNode;
};

export default function UserDashboardLayout({ children }: LayoutProps) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="flex-1 p-4">{children}</main>
      <Footer />
    </div>
  );
}
