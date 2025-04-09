"use client";
import React from "react";
import { cn } from "@/lib/utils";

interface LoadingSpinnerProps {
  size?: "sm" | "md" | "lg";
  className?: string;
  text?: string;
}

export function LoadingSpinner({
  size = "md",
  className,
  text,
}: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: "h-4 w-4 border-2",
    md: "h-8 w-8 border-3",
    lg: "h-12 w-12 border-4",
  };

  const waveColors = ["#38bdf8", "#818cf8", "#c084fc", "#e879f9", "#22d3ee"];
  const gradientBorder = `conic-gradient(from 0deg, ${waveColors.join(", ")})`;

  return (
    <div className={cn("flex flex-col items-center justify-center", className)}>
      <div
        className={cn("rounded-full animate-spin", sizeClasses[size])}
        style={{
          borderStyle: "solid",
          borderColor: "transparent",
          borderTopColor: "#c084fc", // You can change this to match your branding
          borderRightColor: "#e879f9",
          borderBottomColor: "#22d3ee",
          borderLeftColor: "#38bdf8",
        }}
      ></div>
      {text && (
        <p className="mt-3 text-center text-sm font-medium text-gray-600">
          {text}
        </p>
      )}
    </div>
  );
}

export function FullPageLoader() {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white/80 z-50">
      <div className="flex flex-col items-center">
        <LoadingSpinner size="lg" />
        <p className="mt-4 text-center text-lg font-medium bg-gradient-to-r from-purple-500 to-cyan-500 bg-clip-text text-transparent animate-pulse">
          Loading...
        </p>
      </div>
    </div>
  );
}

export function ContentLoader() {
  return (
    <div className="flex items-center justify-center p-8 rounded-lg">
      <LoadingSpinner text="Loading content..." />
    </div>
  );
}

export function ButtonLoader({ className }: { className?: string }) {
  return <LoadingSpinner size="sm" className={cn("ml-2", className)} />;
}
