/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{templ,go}"],
  theme: {
    extend: {
      fontFamily: {
        Workbench: ["Workbench", "sans-serif"],
        RobotoMono: ["RobotoMono", "monospace"],
      },
    },
  },
  plugins: [],
};
