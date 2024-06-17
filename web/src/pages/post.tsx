import { createStore } from 'solid-js/store'
import {
	createEffect,
	createSignal,
	For,
	Match,
	onCleanup,
	Show,
	Switch,
} from 'solid-js'
import { useMainButton } from '~/lib/useMainButton'
import { IconClose, IconMap } from '~/components/icons'
import { fetchCreatePost, fetchPresignedUrl, fetchTags } from '~/lib/api'
import { useNavigate } from '@solidjs/router'
import { createQuery } from '@tanstack/solid-query'
import { cn } from '~/lib/utils'
import { queryClient } from '~/App'

type CreatePost = {
	text: string
	photo: string
	location: string
	tags: number[]
}

async function uploadToS3(url: string, file: File) {
	const response = await fetch(url, {
		method: 'PUT',
		body: file,
		headers: {
			'Content-Type': file.type,
		},
	})
	if (!response.ok) {
		throw new Error('Failed to upload image to S3')
	}
}

export default function PostPage() {
	const [editPost, setEditPost] = createStore({
		text: '',
		location: '',
		photo: '',
		tags: [],
	} as CreatePost)

	const [loading, setLoading] = createSignal(false)
	const mainButton = useMainButton()

	const [imgFile, setImgFile] = createSignal<File | null>(null)
	const [previewUrl, setPreviewUrl] = createSignal('')

	const tagsQuery = createQuery(() => ({
		queryKey: ['tags'],
		queryFn: () => fetchTags(),
	}))

	const navigate = useNavigate()

	const savePost = async () => {
		if (imgFile() && imgFile() !== null) {
			mainButton.disable('Save').showProgress(true)
			try {
				const { file_name, url } = await fetchPresignedUrl(imgFile()!.name)
				await uploadToS3(url, imgFile()!)
				setEditPost('photo', file_name)
				await fetchCreatePost(editPost)
				await queryClient.invalidateQueries({ queryKey: ['posts'] })
				navigate('/')
			} catch (e) {
				console.error(e)
			} finally {
				mainButton.enable('Save').hideProgress()
				setImgFile(null)
				setPreviewUrl('')
			}
		}
	}

	createEffect(() => {
		if (imgFile()) {
			mainButton.enable('Save').onClick(savePost)
		} else {
			mainButton.disable('Save').onClick(savePost)
		}
	})

	onCleanup(() => {
		mainButton.hide().offClick(savePost)
	})

	const [currentLocation, setCurrentLocation] = createSignal<string>('')

	const getLocationName = async (latitude: number, longitude: number) => {
		setLoading(true)
		try {
			const response = await fetch(
				`https://nominatim.openstreetmap.org/reverse?format=json&lat=${latitude}&lon=${longitude}`,
			)
			const data = await response.json()
			return data.display_name || `${latitude}, ${longitude}`
		} catch (error) {
			console.error('Error fetching location name:', error)
			return `${latitude}, ${longitude}`
		} finally {
			setLoading(false)
		}
	}

	const requestLocation = () => {
		if (navigator.geolocation) {
			navigator.geolocation.getCurrentPosition(
				async (position) => {
					const { latitude, longitude } = position.coords
					const locationName = await getLocationName(latitude, longitude)
					setCurrentLocation(locationName)
					setEditPost('location', locationName)
				},
				(error) => {
					console.error('Error getting geolocation:', error)
				},
			)
		} else {
			console.error('Geolocation is not supported by this browser.')
		}
	}

	const handleFileChange = (event: any) => {
		const file = event.target.files[0]
		if (file) {
			const maxSize = 1024 * 1024 * 5 // 7MB

			if (file.size > maxSize) {
				window.Telegram.WebApp.showAlert('Try to select a smaller file')
				return
			}

			setImgFile(file)
			setPreviewUrl('')

			const reader = new FileReader()
			reader.onload = (e) => {
				setPreviewUrl(e.target?.result as string)
			}
			reader.readAsDataURL(file)
		}
	}

	const resolveImage = () => {
		return previewUrl() || null
	}

	return (
		<section class="min-h-screen bg-secondary px-4 pb-14 pt-5">
			<p class="text-2xl font-bold text-foreground">
				What are you cooking today?
			</p>
			<p class="text-hint">Share your delicious meal with the world</p>
			<label class="mt-8 block text-sm font-medium text-foreground">
				Description
				<textarea
					class="mt-2 h-32 w-full resize-none rounded-lg border bg-transparent p-2 text-foreground"
					placeholder="Describe what do you feel like sharing today"
					value={editPost.text}
					onInput={(e) => setEditPost('text', e.currentTarget.value)}
				/>
			</label>
			<Show
				when={!previewUrl()}
				fallback={
					<ImagePreview img={previewUrl()} onRemove={() => setImgFile(null)} />
				}
			>
				<label class="mt-4 flex h-10 items-center justify-start gap-4 rounded-lg border px-2 text-sm font-medium text-foreground">
					<span class="text-nowrap">Choose picture</span>
					<input
						type="file"
						class="sr-only mt-2 w-full rounded-lg bg-transparent p-2 text-foreground"
						placeholder="Enter image"
						accept="image/*"
						onChange={(e) => handleFileChange(e)}
					/>
				</label>
			</Show>
			<label
				class="mt-4 block w-full text-sm font-medium text-foreground"
				for="location"
			>
				Location
			</label>
			<div class="mt-2 flex flex-row items-center justify-between space-x-2">
				<input
					type="text"
					id="location"
					class="w-full rounded-lg border bg-transparent p-2 text-foreground"
					placeholder="Enter location"
					value={editPost.location}
					onInput={(e) => setEditPost('location', e.currentTarget.value)}
				/>
				<button
					class="flex h-full w-8 items-center justify-center rounded-full bg-transparent"
					onClick={() => requestLocation()}
					disabled={loading()}
				>
					<Switch>
						<Match when={!loading()}>
							<IconMap class="size-5 text-foreground" />
						</Match>
						<Match when={loading()}>
							<Spinner />
						</Match>
					</Switch>
				</button>
			</div>
			<div class="mt-4 flex flex-col items-start justify-between">
				<label class="text-sm">Optional tags</label>
				<div class="mt-2 flex flex-row flex-wrap items-center justify-start gap-2">
					<For each={tagsQuery.data}>
						{(tag) => (
							<button
								class={cn(
									'flex h-8 items-center justify-center rounded-lg bg-background px-4 text-sm font-medium text-foreground',
									editPost.tags.includes(tag.id) &&
										'bg-primary text-primary-foreground',
								)}
								onClick={() =>
									setEditPost(
										'tags',
										editPost.tags.includes(tag.id)
											? editPost.tags.filter((t) => t !== tag.id)
											: [...editPost.tags, tag.id],
									)
								}
							>
								{tag.name}
							</button>
						)}
					</For>
				</div>
			</div>
		</section>
	)
}

function Spinner() {
	return (
		<div role="status">
			<svg
				aria-hidden="true"
				class="size-5 animate-spin fill-accent-foreground text-hint"
				viewBox="0 0 100 101"
				fill="none"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
					fill="currentColor"
				/>
				<path
					d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
					fill="currentFill"
				/>
			</svg>
			<span class="sr-only">Loading...</span>
		</div>
	)
}

function ImagePreview(props: { img: string; onRemove: () => void }) {
	return (
		<div class="relative mt-4 flex flex-row items-center justify-between">
			<img
				src={props.img}
				class="aspect-[4/3] w-full rounded-xl object-cover"
				alt="Preview"
			/>
			<button
				onClick={() => props.onRemove()}
				class="absolute right-3 top-3 flex size-7 items-center justify-center rounded-full bg-secondary text-foreground"
			>
				<IconClose class="size-4" />
			</button>
		</div>
	)
}
