import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function clamp(val: number, min: number, max: number) {
	return val > max ? max : val < min ? min : val
}

export function timeSince(dateString: string) {
	const date = new Date(dateString)
	const now = new Date()
	const seconds = Math.floor((now - date) / 1000)
	const minutes = Math.floor(seconds / 60)
	const hours = Math.floor(minutes / 60)
	const days = Math.floor(hours / 24)
	const months = Math.floor(days / 30)
	const years = Math.floor(days / 365)

	if (years > 0) {
		return years + 'y ago'
	} else if (months > 0) {
		return months + 'mo ago'
	} else if (days > 0) {
		return days + 'd ago'
	} else if (hours > 0) {
		return hours + 'h ago'
	} else if (minutes > 0) {
		return minutes + 'm ago'
	} else if (seconds <= 10) {
		return 'just now'
	} else {
		return seconds + 's ago'
	}
}
