/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'dark-bg': '#1E1E2E',
        'dark-surface': '#2B2B3D',
        'dark-primary': '#3E7BFA',
        'dark-secondary': '#8E8EA0',
        'dark-error': '#E53E3E',
        'dark-on-bg': '#E1E1E1',
        'dark-on-surface': '#FFFFFF',
      },
    },
  },
  plugins: [],
}
