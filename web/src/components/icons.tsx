import { ComponentProps, splitProps } from 'solid-js'
import { cn } from '~/lib/utils'

type IconProps = ComponentProps<'svg'>

const Icon = (props: IconProps) => {
	const [, rest] = splitProps(props, ['class'])
	return (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			{...rest}
		/>
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
			<path
				d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z" />
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

export function IconProteins(props: IconProps) {
	return (
		<Icon {...props}>
			<rect width="24" height="24" rx="12" fill="#D6F961" />
			<mask id="mask0_23_5230" style={{ 'mask-type': 'alpha' }} maskUnits="userSpaceOnUse" x="5" y="5" width="14"
				height="14">
				<rect x="5" y="5" width="14" height="14" fill="#D9D9D9" />
			</mask>
			<g mask="url(#mask0_23_5230)">
				<path
					d="M12 13.75C12.4861 13.75 12.8993 13.5799 13.2396 13.2396C13.5799 12.8993 13.75 12.4861 13.75 12C13.75 11.5139 13.5799 11.1007 13.2396 10.7604C12.8993 10.4202 12.4861 10.25 12 10.25C11.5139 10.25 11.1007 10.4202 10.7604 10.7604C10.4201 11.1007 10.25 11.5139 10.25 12C10.25 12.4861 10.4201 12.8993 10.7604 13.2396C11.1007 13.5799 11.5139 13.75 12 13.75ZM12.0041 17.6C11.2333 17.6 10.5076 17.4542 9.82708 17.1625C9.14652 16.8709 8.55104 16.4698 8.04062 15.9594C7.5302 15.449 7.12916 14.8537 6.83749 14.1735C6.54583 13.4934 6.39999 12.7666 6.39999 11.9933C6.39999 11.22 6.54583 10.4955 6.83749 9.81982C7.12916 9.14412 7.5302 8.55107 8.04062 8.04065C8.55104 7.53023 9.14633 7.12919 9.82649 6.83752C10.5067 6.54586 11.2334 6.40002 12.0067 6.40002C12.78 6.40002 13.5045 6.54586 14.1802 6.83752C14.8559 7.12919 15.449 7.53023 15.9594 8.04065C16.4698 8.55107 16.8708 9.14519 17.1625 9.82302C17.4542 10.501 17.6 11.2253 17.6 11.9959C17.6 12.7667 17.4542 13.4924 17.1625 14.1729C16.8708 14.8535 16.4698 15.449 15.9594 15.9594C15.449 16.4698 14.8548 16.8709 14.177 17.1625C13.4991 17.4542 12.7748 17.6 12.0041 17.6Z"
					fill="#1C1B1F" />
			</g>
		</Icon>
	)
}

export function IconCarbs(props: IconProps) {
	return (
		<Icon {...props}>
			<rect width="24" height="24" rx="12" fill="#F9E5A8" />
			<mask id="mask0_23_5238" style={{ 'mask-type': 'alpha' }} maskUnits="userSpaceOnUse" x="5" y="5" width="14"
				height="14">
				<rect x="5" y="5" width="14" height="14" fill="#D9D9D9" />
			</mask>
			<g mask="url(#mask0_23_5238)">
				<path
					d="M12 17.6C11.8444 17.6 11.7083 17.5757 11.5917 17.5271C11.475 17.4785 11.3632 17.4007 11.2562 17.2938L6.70624 12.7438C6.5993 12.6368 6.52152 12.525 6.47291 12.4084C6.4243 12.2917 6.39999 12.1556 6.39999 12C6.39999 11.8445 6.4243 11.7084 6.47291 11.5917C6.52152 11.475 6.5993 11.3632 6.70624 11.2563L11.2562 6.70627C11.3632 6.59933 11.475 6.52155 11.5917 6.47294C11.7083 6.42433 11.8444 6.40002 12 6.40002C12.1555 6.40002 12.2917 6.42433 12.4083 6.47294C12.525 6.52155 12.6368 6.59933 12.7437 6.70627L17.2937 11.2563C17.4007 11.3632 17.4785 11.475 17.5271 11.5917C17.5757 11.7084 17.6 11.8445 17.6 12C17.6 12.1556 17.5757 12.2917 17.5271 12.4084C17.4785 12.525 17.4007 12.6368 17.2937 12.7438L12.7437 17.2938C12.6368 17.4007 12.525 17.4785 12.4083 17.5271C12.2917 17.5757 12.1555 17.6 12 17.6Z"
					fill="#1C1B1F" />
			</g>
		</Icon>
	)
}

export function IconFats(props: IconProps) {
	return (
		<Icon {...props}>
			<rect width="24" height="24" rx="12" fill="#EBBFCE" />
			<mask id="mask0_23_5246" style={{ 'mask-type': 'alpha' }} maskUnits="userSpaceOnUse" x="5" y="5" width="14"
				height="14">
				<rect x="5" y="5" width="14" height="14" fill="#D9D9D9" />
			</mask>
			<g mask="url(#mask0_23_5246)">
				<path
					d="M12.0041 17.6C11.2333 17.6 10.5076 17.4542 9.82708 17.1625C9.14652 16.8709 8.55104 16.4698 8.04062 15.9594C7.5302 15.449 7.12916 14.8537 6.83749 14.1735C6.54583 13.4934 6.39999 12.7666 6.39999 11.9933C6.39999 11.22 6.54583 10.4955 6.83749 9.81982C7.12916 9.14412 7.5302 8.55107 8.04062 8.04065C8.55104 7.53023 9.14633 7.12919 9.82649 6.83752C10.5067 6.54586 11.2334 6.40002 12.0067 6.40002C12.78 6.40002 13.5045 6.54586 14.1802 6.83752C14.8559 7.12919 15.449 7.53023 15.9594 8.04065C16.4698 8.55107 16.8708 9.14519 17.1625 9.82302C17.4542 10.501 17.6 11.2253 17.6 11.9959C17.6 12.7667 17.4542 13.4924 17.1625 14.1729C16.8708 14.8535 16.4698 15.449 15.9594 15.9594C15.449 16.4698 14.8548 16.8709 14.177 17.1625C13.4991 17.4542 12.7748 17.6 12.0041 17.6Z"
					fill="#1C1B1F" />
			</g>

		</Icon>
	)
}

export function IconChevronLeft(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m15 18-6-6 6-6" />
		</Icon>
	)
}

export function IconChevronRight(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m9 18 6-6-6-6" />
		</Icon>
	)
}