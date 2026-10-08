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
      fontSize: {
        xs: ['0.875rem', { lineHeight: '1.5' }],
        sm: ['1rem', { lineHeight: '1.55' }],
        base: ['1.0625rem', { lineHeight: '1.7' }],
        lg: ['1.1875rem', { lineHeight: '1.6' }],
        xl: ['1.375rem', { lineHeight: '1.45' }],
        '2xl': ['1.625rem', { lineHeight: '1.4' }],
        '3xl': ['1.875rem', { lineHeight: '1.35' }],
        '4xl': ['2.25rem', { lineHeight: '1.28' }],
        '5xl': ['2.75rem', { lineHeight: '1.2' }],
        '6xl': ['3.375rem', { lineHeight: '1.14' }],
        '7xl': ['4rem', { lineHeight: '1.12' }],
      },
    },
  },
  plugins: [],
} satisfies Config;
