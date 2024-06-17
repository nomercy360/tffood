import { For, onCleanup, onMount, Show } from 'solid-js'
import { IconUpset, IconSmile, IconNeutral, IconMap } from '~/components/icons'
import { cn, timeSince } from '~/lib/utils'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import {
	fetchAddPostReaction,
	fetchPosts,
	fetchRemovePostReaction,
} from '~/lib/api'
import { createMutation, createQuery } from '@tanstack/solid-query'
import { queryClient } from '~/App'
import { Link } from '~/components/link'

type Post = {
	id: number
	user_id: number
	text: string
	created_at: string
	updated_at: string
	hidden_at: string | null
	photo_url: string
	suggested_ingredients: string[]
	suggested_dish_name: string
	suggested_tags: string[]
	reactions: {
		frown: number
		smile: number
		meh: number
	}
	location_id: number | null
	location: {
		id: number
		latitude: number
		longitude: number
		address: string
	}
	user_reaction: {
		type: string
		has_reacted: boolean
	}
	user: {
		id: number
		username: string
		first_name: string
		last_name: string
		avatar_url: string
	}
}

export default function HomePage() {
	const query = createQuery(() => ({
		queryKey: ['posts'],
		queryFn: () => fetchPosts(),
	}))

	const handleMutate = async (id: number, type: string) => {
		await queryClient.cancelQueries({ queryKey: ['posts'] })

		queryClient.setQueryData(['posts'], (old: Post[]) =>
			old.map((post) => {
				if (post.id !== id) return post

				const reactionsCopy = { ...post.reactions }

				const isTogglingSameReaction =
					post.user_reaction.has_reacted && post.user_reaction.type === type

				if (isTogglingSameReaction) {
					reactionsCopy[type] = Math.max(0, reactionsCopy[type] - 1)
				} else {
					if (post.user_reaction.has_reacted) {
						reactionsCopy[post.user_reaction.type] = Math.max(
							0,
							reactionsCopy[post.user_reaction.type] - 1,
						)
					}
					reactionsCopy[type] = (reactionsCopy[type] || 0) + 1
				}

				const newUserReaction = {
					type: type,
					has_reacted: !isTogglingSameReaction,
				}

				return {
					...post,
					reactions: reactionsCopy,
					user_reaction: newUserReaction,
				}
			}),
		)

		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	const mutateReact = createMutation(() => ({
		mutationFn: async ({ item, type }: { item: Post; type: string }) => {
			if (item.user_reaction.has_reacted && item.user_reaction.type === type) {
				await fetchRemovePostReaction(item.id)
			} else {
				await fetchAddPostReaction(item.id, type as any)
			}
		},
		onMutate: ({ item, type }: { item: Post; type: string }) =>
			handleMutate(item.id, type),
	}))

	const mainButton = useMainButton()
	const navigate = useNavigate()

	function navigateToPost() {
		navigate('/post')
	}

	onMount(() => {
		mainButton.enable('Post Food').onClick(navigateToPost)
	})

	onCleanup(() => {
		mainButton.offClick(navigateToPost)
		mainButton.hide()
	})

	return (
		<section class="p-4">
			<div class="grid gap-2">
				<Show when={query.isSuccess} fallback={<Loader />}>
					<For each={query.data as Post[]}>
						{(item) => (
							<div class="rounded-lg border bg-section">
								<Link
									class="flex flex-row items-center justify-start gap-2 p-4"
									href={`/users/${item.user.username}`}
								>
									<img
										src={item.user.avatar_url}
										class="size-8 rounded-full"
										alt="User"
									/>
									<div class="flex flex-col items-start justify-start">
										<p class="text-sm font-semibold text-hint">
											{item.user.username}
										</p>
										<Show when={item.location}>
											<div class="flex flex-row items-center justify-start gap-1.5">
												<p class="line-clamp-1 text-xs text-hint">
													{item.location.address}
												</p>
											</div>
										</Show>
									</div>
								</Link>
								<img
									src={item.photo_url}
									class="aspect-[4/3] w-full object-cover"
									alt="Thumbnail"
								/>
								<div class="p-3.5">
									<p class="text-sm text-hint">
										{item.text || item.suggested_dish_name}
									</p>
									<div class="mt-4 flex flex-row flex-wrap items-center justify-start gap-1.5">
										<For each={item.suggested_tags}>
											{(ingredient) => (
												<span class="flex h-6 items-center justify-center rounded-lg bg-background px-2 py-0.5 text-xs text-hint">
													{ingredient}
												</span>
											)}
										</For>
									</div>
									<div class="mt-4 flex w-full flex-row items-center justify-between">
										<div class="flex flex-row items-center justify-start gap-2">
											<button
												class={cn(
													'flex h-8 w-12 items-center justify-start gap-1.5 rounded-lg px-1.5 text-sm',
													item.user_reaction.has_reacted &&
														item.user_reaction.type === 'smile' &&
														'bg-background',
												)}
												onClick={() =>
													mutateReact.mutate({ item, type: 'smile' })
												}
											>
												<IconSmile class="shrink-0" />
												{item.reactions.smile}
											</button>
											<button
												class={cn(
													'flex h-8 w-12 items-center justify-start gap-1.5 rounded-lg px-1.5 text-sm',
													item.user_reaction.has_reacted &&
														item.user_reaction.type === 'meh' &&
														'bg-background',
												)}
												onClick={() =>
													mutateReact.mutate({ item, type: 'meh' })
												}
											>
												<IconNeutral class="shrink-0" />
												{item.reactions.meh}
											</button>
											<button
												class={cn(
													'flex h-8 w-12 items-center justify-start gap-1.5 rounded-lg px-1.5 text-sm',
													item.user_reaction.has_reacted &&
														item.user_reaction.type === 'frown' &&
														'bg-background',
												)}
												onClick={() => {
													mutateReact.mutate({ item, type: 'frown' })
												}}
											>
												<IconUpset class="shrink-0" />
												{item.reactions.frown}
											</button>
										</div>
										<span class="text-xs text-hint">
											{timeSince(item.created_at)}
										</span>
									</div>
								</div>
							</div>
						)}
					</For>
				</Show>
			</div>
		</section>
	)
}

function Loader() {
	return (
		<div class="grid gap-2">
			<div class="h-80 animate-pulse rounded-lg border bg-section" />
			<div class="h-80 animate-pulse rounded-lg border bg-section" />
			<div class="h-80 animate-pulse rounded-lg border bg-section" />
			<div class="h-80 animate-pulse rounded-lg border bg-section" />
		</div>
	)
}
