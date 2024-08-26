const flowbite = require("flowbite-react/tailwind");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./src/**/**/*.{js,ts,jsx,tsx}",
    "./src/**/**/**/*.{js,ts,jsx,tsx}",
    flowbite.content(),
  ],
  theme: {
    extend: {
      gridTemplateRows: {
        "1fr": "auto",
        "2fr": "auto 1fr",
        "3fr": "auto 1fr auto",
        "4fr": "auto 1fr auto auto",
      },
      gridTemplateColumns: {
        "1fr": "auto",
        "2fr": "auto 1fr",
        "3fr": "auto 1fr auto",
        "4fr": "auto 1fr auto auto",
      },
    },
  },
  plugins: [
    flowbite.plugin(),
    require("flowbite/plugin")({
      charts: true,
    }),
  ],
};
