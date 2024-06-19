import { createSignal, createEffect, onCleanup, For } from 'solid-js'
import { IconTriangle } from '~/components/icons'

export const [toasts, setToasts] = createSignal<
	{ id: number; message: string }[]
>([])

export const addToast = (message: string) => {
	const id = Date.now()
	setToasts([...toasts(), { id, message }])

	// Remove the toast after 3 seconds
	setTimeout(() => {
		setToasts(toasts().filter(toast => toast.id !== id))
	}, 3000)
}

const Toast = () => {
	createEffect(() => {
		const currentToasts = toasts()
		if (currentToasts.length > 5) {
			const newToasts = currentToasts.slice(1)
			setToasts(newToasts)
		}
	})

	return (
		<div class="fixed bottom-4 left-1/2 -translate-x-1/2 space-y-2">
			<For each={toasts()}>{toast => (
				<div
					class="flex h-9 w-[calc(100vw-2rem)] items-center justify-start rounded-xl bg-red-500 px-4 py-2 text-sm font-medium text-white">
					<div class="mr-2 flex size-6 items-center justify-center rounded-full bg-red-700">
						<IconTriangle class="size-4 shrink-0 text-white" />
					</div>
					{toast.message}
				</div>
			)}</For>
		</div>
	)
}

export default Toast
