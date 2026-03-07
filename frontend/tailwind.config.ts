import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#0f172a",
        mist: "#f8fafc",
        tide: "#0f766e",
        ember: "#c2410c",
        storm: "#334155"
      },
      boxShadow: {
        panel: "0 22px 60px rgba(15, 23, 42, 0.14)"
      },
      fontFamily: {
        display: ["Manrope", "sans-serif"],
        mono: ["IBM Plex Mono", "monospace"]
      },
      backgroundImage: {
        "panel-grid":
          "linear-gradient(rgba(148,163,184,.08) 1px, transparent 1px), linear-gradient(90deg, rgba(148,163,184,.08) 1px, transparent 1px)"
      }
    }
  },
  plugins: []
} satisfies Config;
