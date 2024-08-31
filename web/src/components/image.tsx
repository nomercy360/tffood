import { createSignal, onCleanup } from 'solid-js'
import { cn } from '~/lib/utils'

const Image = (props: { src: string; alt: string, class: string }) => {
	const [loading, setLoading] = createSignal(true)
	const [displayed, setDisplayed] = createSignal(false)

	const minDisplayTime = 10000

	const handleImageLoad = () => {
		setTimeout(() => {
			setDisplayed(true)
			if (!loading()) {
				setLoading(false)
			}
		}, minDisplayTime)
	}

	onCleanup(() => {
		setLoading(false)
		setDisplayed(false)
	})

	return (
		<div class="relative">
			{(loading() || !displayed()) && (
				<div class="absolute inset-0 flex items-center justify-center">
					<div class="size-full rounded-[20px]" />
				</div>
			)}
			<img
				src={props.src}
				alt={props.alt}
				class={cn('h-auto w-full object-cover transition-opacity duration-500',
					props.class,
					loading() && !displayed() ? 'opacity-0 min-h-40' : 'opacity-100')
				}
				loading="lazy"
				onLoad={() => {
					handleImageLoad()
					setLoading(false)
				}}
			/>
		</div>
	)
}

export default Image
