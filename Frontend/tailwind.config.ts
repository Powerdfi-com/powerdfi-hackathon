import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "selector",
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {
        primary: "#C02CF4",
        secondary: "#1BDEFD",
        background: "#111111",
        "white-text": "#FFFFFF",
        "text-grey": "#303030",
        "black-shade": "#191919",
        "white-shade": "#191919",
        border: "#303030",
        shade: "#191919",
        "blue-shade": "#272D37",
        "text-grey2": "#828F9C",
        "text-blue": "#1BDEFD26",
      },
    },
  },
  plugins: [],
};
export default config;
