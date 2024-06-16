import { ComponentProps, splitProps } from 'solid-js'
import { cn } from '~/lib/utils'

type IconProps = ComponentProps<'svg'>

const Icon = (props: IconProps) => {
	const [, rest] = splitProps(props, ['class'])
	return (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class={cn('size-5', props.class)}
			{...rest}
		/>
	)
}

export function IconUpset(props: IconProps) {
	return (
		<Icon {...props}>
			<circle cx="12" cy="12" r="10" />
			<path d="M16 16s-1.5-2-4-2-4 2-4 2" />
			<line x1="9" x2="9.01" y1="9" y2="9" />
			<line x1="15" x2="15.01" y1="9" y2="9" />
		</Icon>
	)
}

export function IconNeutral(props: IconProps) {
	return (
		<Icon {...props}>
			<circle cx="12" cy="12" r="10" />
			<line x1="8" x2="16" y1="15" y2="15" />
			<line x1="9" x2="9.01" y1="9" y2="9" />
			<line x1="15" x2="15.01" y1="9" y2="9" />
		</Icon>
	)
}

export function IconSmile(props: IconProps) {
	return (
		<Icon {...props}>
			<circle cx="12" cy="12" r="10" />
			<path d="M8 14s1.5 2 4 2 4-2 4-2" />
			<line x1="9" x2="9.01" y1="9" y2="9" />
			<line x1="15" x2="15.01" y1="9" y2="9" />
		</Icon>
	)
}

export function IconClose(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M18 6L6 18" />
			<path d="M6 6l12 12" />
		</Icon>
	)
}

export function IconMap(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z" />
			<circle cx="12" cy="10" r="3" />
		</Icon>
	)
}
