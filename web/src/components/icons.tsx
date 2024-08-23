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

export function IconSparkles(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z" />
			<path d="M20 3v4" />
			<path d="M22 5h-4" />
			<path d="M4 17v2" />
			<path d="M5 18H3" />
		</Icon>
	)
}

export function IconTriangle(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3" />
			<path d="M12 9v4" />
			<path d="M12 17h.01" />
		</Icon>
	)
}

export function IconChevronDown(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m6 9 6 6 6-6" />
		</Icon>
	)
}

export function IconChevronUp(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m18 15-6-6-6 6" />
		</Icon>
	)
}

export function IconInfo(props: IconProps) {
	return (
		<Icon {...props}>
			<circle cx="12" cy="12" r="10" />
			<path d="M12 16v-4" />
			<path d="M12 8h.01" />
		</Icon>
	)
}

export function IconShare(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8" />
			<polyline points="16 6 12 2 8 6" />
			<line x1="12" x2="12" y1="2" y2="15" />
		</Icon>
	)
}
