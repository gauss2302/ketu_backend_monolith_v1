import React from "react";

export function Footer() {
  const waveColors = ["#38bdf8", "#818cf8", "#c084fc", "#e879f9", "#22d3ee"];

  const gradient = `linear-gradient(to right, ${waveColors.join(", ")})`;

  return (
    <footer
      className="p-4 text-center text-white"
      style={{ background: gradient }}
    >
      <p>&copy; {new Date().getFullYear()} My App. All rights reserved.</p>
    </footer>
  );
}
