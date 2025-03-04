import React from "react";

export function Footer() {
  return (
    <footer className="p-4 text-center bg-neutral-100 dark:bg-neutral-900 border-t dark:border-neutral-800">
      <p>&copy; {new Date().getFullYear()} My App. All rights reserved.</p>
    </footer>
  );
}
