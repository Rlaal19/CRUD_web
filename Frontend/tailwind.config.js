import daisyui from "daisyui";

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  daisyui: {
    themes: [
      "corporate",
      {
        mytheme: {
          "star": "#ff7087",
          "primary": "#7c3aed",
          "secondary": "#0000FE",
          "accent": "#a3e635",
          "neutral": "#00ff00",
          "base-100": "#ffffff",
          "info": "#ff00ff",
          "success": "#00ff00",
          "warning": "#fbbf24",
          "error": "#ff0000",
        },
      },
    ],
  },
  plugins: [daisyui],
};
