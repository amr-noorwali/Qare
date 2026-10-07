import type { Config } from 'tailwindcss';
export default { content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}'], theme: { extend: { colors: { ink: '#172d2a', moss: '#346b5b', cream: '#f8f5ee', gold: '#e5af56' }, fontFamily: { sans: ['Arial', 'sans-serif'] } } }, plugins: [] } satisfies Config;
