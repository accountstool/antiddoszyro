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
        storm: "#334155",
        aqua: "#14b8a6",
        shell: "#08111d",
        cloud: "#ecf4f4"
      },
      boxShadow: {
        panel: "0 30px 80px rgba(15, 23, 42, 0.12)",
        lift: "0 18px 40px rgba(15, 23, 42, 0.10)",
        inset: "inset 0 1px 0 rgba(255,255,255,0.45)"
      },
      fontFamily: {
        display: ["Space Grotesk", "sans-serif"],
        sans: ["Manrope", "sans-serif"],
        mono: ["IBM Plex Mono", "monospace"]
      },
      backgroundImage: {
        "panel-grid":
          "linear-gradient(rgba(148,163,184,.08) 1px, transparent 1px), linear-gradient(90deg, rgba(148,163,184,.08) 1px, transparent 1px)",
        "mesh-teal":
          "radial-gradient(circle at top left, rgba(20,184,166,.24), transparent 28%), radial-gradient(circle at bottom right, rgba(249,115,22,.18), transparent 24%), linear-gradient(135deg, rgba(15,23,42,.96), rgba(8,17,29,.88))",
        "panel-wash":
          "linear-gradient(180deg, rgba(255,255,255,.92), rgba(244,248,250,.76))"
      }
    }
  },
  plugins: []
} satisfies Config;
