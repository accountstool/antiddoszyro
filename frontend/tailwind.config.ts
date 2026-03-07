import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#111111",
        mist: "#f5f5f5",
        tide: "#111111",
        ember: "#404040",
        storm: "#525252",
        aqua: "#171717",
        shell: "#050505",
        cloud: "#f3f3f3"
      },
      boxShadow: {
        panel: "0 24px 64px rgba(0, 0, 0, 0.10)",
        lift: "0 16px 36px rgba(0, 0, 0, 0.18)",
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
          "radial-gradient(circle at top left, rgba(255,255,255,.08), transparent 28%), radial-gradient(circle at bottom right, rgba(255,255,255,.04), transparent 24%), linear-gradient(135deg, rgba(10,10,10,.98), rgba(28,28,28,.94))",
        "panel-wash":
          "linear-gradient(180deg, rgba(255,255,255,.92), rgba(244,244,245,.82))"
      }
    }
  },
  plugins: []
} satisfies Config;
