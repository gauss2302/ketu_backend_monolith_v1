// app/dashboard/user/layout.tsx
import React, { ReactNode } from "react";

type LayoutProps = {
  children: ReactNode;
};

export default function UserDashboardLayout({ children }: LayoutProps) {
  return <>{children}</>;
}
