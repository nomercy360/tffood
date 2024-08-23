import { createStore } from 'solid-js/store'
import { createEffect, createSignal, For, Match, onCleanup, onMount, Show, Switch } from 'solid-js'
import { useMainButton } from '~/lib/useMainButton'
import { IconClose, IconMap, IconSparkles } from '~/components/icons'
import { fetchCreatePost, fetchPostAISuggestions, fetchPresignedUrl, fetchUpdatePost } from '~/lib/api'
import { useNavigate } from '@solidjs/router'
import { queryClient } from '~/App'
import { cn } from '~/lib/utils'
import { addToast } from '~/components/toast'

type CreatePost = {
	text: string | null
	photo: string
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
	const [editPost, setEditPost] = createStore<CreatePost>({
		text: '',
		photo: '',
	})

	const [step, setStep] = createSignal(0)

	const [postLoading, setPostLoading] = createSignal(false)

	const mainButton = useMainButton()

	const [imgFile, setImgFile] = createSignal<File | null>(null)
	const [previewUrl, setPreviewUrl] = createSignal('')

	const navigate = useNavigate()

	const createPost = async () => {
		if (imgFile() && imgFile() !== null) {
			mainButton.showProgress(false)
			try {
				const { file_name, url } = await fetchPresignedUrl(imgFile()!.name)
				await uploadToS3(url, imgFile()!)
				setEditPost('photo', file_name)
				const resp = await fetchCreatePost(editPost)
				await queryClient.invalidateQueries({ queryKey: ['posts'] })
				navigate('/')
			} catch (e) {
				console.error(e)
			} finally {
				mainButton.hideProgress()
			}
		}
	}

	const errorToast = () => {
		addToast('Some error occurred')
	}

	onMount(async () => {
		mainButton.onClick(createPost)
	})

	createEffect(() => {
		if (step() === 0) {
			if (imgFile()) {
				mainButton.enable('Save')
			} else {
				mainButton.disable('Save')
			}
		}
	})

	onCleanup(() => {
		mainButton.hide()
			.offClick(createPost)
	})

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
		<>
			<Switch>
				<Match when={step() === 0}>
					<Layout
						title="Share your experience"
						subtitle="What are you eating today?"
					>
						<Show when={!imgFile()}
							fallback={<ImagePreview img={previewUrl()} onRemove={() => setImgFile(null)} />}>
							<label
								class="mt-4 flex h-10 items-center justify-start gap-4 rounded-lg border px-2 text-sm font-medium text-foreground">
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
						<Show when={imgFile()}>
							<label class="mt-6 block text-sm font-medium text-foreground">
								Description
								<div class="mt-2 flex flex-row items-center justify-between space-x-2">
									<input
										class="h-10 w-full resize-none rounded-lg border bg-transparent px-2 text-base font-normal text-foreground"
										placeholder="Describe what do you feel like sharing today"
										value={editPost.text || ''}
										onInput={(e) => setEditPost('text', e.currentTarget.value)}
										disabled={postLoading()}
									/>
								</div>
							</label>
						</Show>
					</Layout>
				</Match>
			</Switch>
		</>
	)
}


function Layout(props: {
	children: any, title: string, subtitle: string
}) {
	return (
		<section class="min-h-screen bg-secondary pb-14 pt-5">
			<div class="px-4">
				<p class="text-2xl font-bold text-foreground">
					{props.title}
				</p>
				<p class="text-hint">
					{props.subtitle}
				</p>
				{props.children}
			</div>
		</section>
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
