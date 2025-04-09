// app/(auth)/layout.tsx
"use client";
import React from "react";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="overflow-hidden h-screen w-full flex items-center justify-center px-4 sm:px-6 lg:px-8">
      {children}
    </div>
  );
}
