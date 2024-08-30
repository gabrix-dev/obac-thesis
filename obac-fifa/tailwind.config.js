/** @type {import('tailwindcss').Config} */
module.exports = {
  css: ["~/assets/css/main.css"],
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./nuxt.config.{js,ts}",
    "./app.vue",
  ],
  theme: {
    extend: {
      colors: {
        "neon-green": "#64feda",
        "neon-gray": "#10172b",
        "neon-pink": "#5E009A",
        "primary-bg": "#000000",
        "secondary-bg": "#FFFFFF",
        "accent-bg-dark": "#C882C8",
        "accent-bg-light": "#F2C2F2",
        "primary-text": "#FFFFFF",
      },
      fontFamily: {
        "press-start": ['"Press Start 2P"', "cursive"],
        "dot-gothic": ['"DotGothic16"', "sans-serif"],
        orbitron: ['"Orbitron"', "sans-serif"],
        audiowide: ['"Audiowide"', "cursive"],
        aldrich: ['"Aldrich"', "sans-serif"],
        exo: ['"Exo 2"', "sans-serif"],
      },
      backgroundImage: {
        desert: "url('./assets/img/desert.jpg')",
        desert2: "url('./assets/img/desert2.jpg')",
        stadium: "url('./assets/img/stadium_inside.jpeg')",
      },
    },
  },
  plugins: [require("tailwind-scrollbar")],
};
