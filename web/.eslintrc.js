module.exports = {
	'env': {
		'browser': true,
		'es2021': true,
	},
	'extends': [
		'plugin:tailwindcss/recommended',
		'plugin:solid/typescript',
	],
	'overrides': [
		{
			'env': {
				'node': true,
			},
			'files': [
				'.eslintrc.{js,cjs}',
			],
			'parserOptions': {
				'sourceType': 'script',
			},
		},
	],
	'parser': '@typescript-eslint/parser',
	'parserOptions': {
		'ecmaVersion': 'latest',
		'sourceType': 'module',
	},
	'plugins': [
		'@typescript-eslint',
		'solid',
	],
	'rules': {
		'indent': [
			'error',
			'tab',
		],
		'linebreak-style': [
			'error',
			'unix',
		],
		'quotes': [
			'error',
			'single',
		],
		'semi': [
			'error',
			'never',
		],
	},
}
