/**@type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{ts,tsx}'],
	theme: {
		container: {
			center: true,
			padding: '2rem',
			screens: {
				'2xl': '1400px',
			},
		},
		extend: {
			colors: {
				input: 'var(--input)',
				background: 'var(--background)',
				foreground: 'var(--foreground)',
				hint: 'var(--hint)',
				link: 'var(--link)',
				primary: {
					DEFAULT: 'var(--primary-bg)',
					foreground: 'var(--primary-foreground)',
				},
				secondary: {
					DEFAULT: 'var(--secondary-bg)',
					foreground: 'var(--secondary-foreground)',
				},
				header: 'var(--header-bg)',
				accent: {
					foreground: 'var(--accent-foreground)',
				},
				section: {
					DEFAULT: 'var(--section-bg)',
					foreground: 'var(--section-text)',
				},
				'subtitle-text': 'var(--subtitle-text)',
				'destructive-text': 'var(--destructive-text)',
			},
		},
	},
	plugins: [require('tailwindcss-animate')],
}
