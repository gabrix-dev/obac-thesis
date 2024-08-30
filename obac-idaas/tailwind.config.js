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
        "entrust-purple": "#690070",
        "entrust-dark-purple": "#54005a",
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
    },
  },
  plugins: [require("tailwind-scrollbar")],
};
