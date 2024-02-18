/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{templ,go}"],
  theme: {
    extend: {
      fontFamily: {
        workbench: ["Workbench", "sans-serif"],
        montserrat: ["Montserrat", "sans-serif"],
      },
    },
  },
  plugins: [],
};
