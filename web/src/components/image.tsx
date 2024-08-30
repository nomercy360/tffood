import { createSignal } from 'solid-js'
import { cn } from '~/lib/utils'

const Image = (props: { src: string; alt: string, class: string }) => {
	const [loading, setLoading] = createSignal(true)

	const handleImageLoad = () => {
		setLoading(false)
	}

	return (
		<div class="relative">
			{loading() && (
				<div class="absolute inset-0 flex items-center justify-center">
					<div class="size-full rounded-[20px] bg-border" />
				</div>
			)}
			<img
				src={props.src}
				alt={props.alt}
				class={cn('h-auto w-full object-cover transition-opacity duration-500',
					props.class,
					loading() ? 'opacity-0 min-h-40' : 'opacity-100')
				}
				loading="lazy"
				onLoad={handleImageLoad}
			/>
		</div>
	)
}

export default Image
