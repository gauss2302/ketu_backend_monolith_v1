// app/(auth)/layout.tsx
"use client";
import React from "react";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div className="overflow-hidden h-screen">{children}</div>;
}
