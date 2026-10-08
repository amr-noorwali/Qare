import type { Config } from 'tailwindcss';

export default {
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#111111',
        paper: '#f7f7f5',
        mist: '#eeeeec',
        line: '#dededb',
        muted: '#696969',
        moss: '#111111',
        cream: '#f7f7f5',
        gold: '#9a9a9a',
      },
      fontFamily: { sans: ['Tajawal', 'Arial', 'sans-serif'] },
    },
  },
  plugins: [],
} satisfies Config;
