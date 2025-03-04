"use client";

import * as React from "react";
import { MoonIcon, SunIcon } from "@radix-ui/react-icons";
import { useTheme } from "next-themes";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

export function ThemeToggle() {
  const { setTheme, theme } = useTheme();

  return (
    <ToggleGroup
      value={theme}
      onValueChange={setTheme}
      type="single"
      aria-label="Theme"
    >
      <ToggleGroupItem value="light" aria-label="Light">
        <SunIcon className="h-[1.2rem] w-[1.2rem]" />
      </ToggleGroupItem>
      <ToggleGroupItem value="dark" aria-label="Dark">
        <MoonIcon className="h-[1.2rem] w-[1.2rem]" />
      </ToggleGroupItem>
      <ToggleGroupItem value="system" aria-label="System">
        <span className="sr-only">System</span>
        <span className="h-[1.2rem] w-[1.2rem] text-sm">Sys</span>
      </ToggleGroupItem>
    </ToggleGroup>
  );
}
