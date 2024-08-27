import { For, Match, onCleanup, onMount, Show, Switch } from 'solid-js'
import { cn, timeSince } from '~/lib/utils'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { fetchPosts } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/link'
import { IconInfo, IconShare } from '~/components/icons'
import { store } from '~/lib/store'
import JoinCommunity from '~/components/join-community'

export type Post = {
	id: number
	user_id: number
	text: string
	created_at: string
	updated_at: string
	hidden_at: string | null
	photo_url: string
	ingredients: { name: string; amount: number }[]
	dish_name: string
	tags: string[]
	user: {
		id: number
		username: string
		first_name: string
		last_name: string
		avatar_url: string
	}
	food_insights: {
		calories: number
		proteins: number
		fats: number
		carbohydrates: number
		dietary_information: string[]
	}
}

export default function HomePage() {
	return (
		<Switch>
			<Match when={store.user.community_status == 'member'}>
				<Feed />
			</Match>
			<Match when={store.user.community_status == 'none'}>
				<JoinCommunity />
			</Match>
		</Switch>
	)
}

function Feed() {
	const query = createQuery(() => ({
		queryKey: ['posts'],
		queryFn: () => fetchPosts(),
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
						{(item) => (<PostCard post={item} class="" />)}
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

export function UserProfileLink(props: { class: any; user: Post['user'] }) {
	return (
		<Link
			class={cn('flex flex-row items-center justify-start gap-2', props.class)}
			href={`/users/${props.user.username}`}
		>
			<img src={props.user.avatar_url} class="size-8 rounded-full" alt="User" />
			<div class="flex flex-col items-start justify-start">
				<p class="text-sm font-semibold text-hint">@{props.user.username}</p>
				<p class="text-sm font-semibold text-foreground">
					{props.user?.first_name} {props.user?.last_name}
				</p>
			</div>
		</Link>
	)
}

export function PostCard(props: { post: Post, class: any }) {
	function sharePostURl(postID: string) {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/eatzfood_bot/app?startapp=p' + postID,
			}).toString() +
			'&text=Check out this post'

		window.Telegram.WebApp.openTelegramLink(url)
	}

	return (
		<div class={cn('rounded-lg border bg-section', props.class)}>
			<UserProfileLink user={props.post.user} class="p-4" />
			<img
				src={props.post.photo_url}
				class="aspect-[4/3] w-full object-cover"
				alt="Thumbnail"
			/>
			<div class="p-3.5">
				<p class="text-sm text-hint">
					{props.post.text || props.post.dish_name}
				</p>
				<div class="mt-4 flex w-full flex-row items-center justify-between">
					<div class="flex flex-row items-center justify-start gap-2">
						<button
							class="flex size-8 flex-row items-center justify-center gap-1.5 rounded-lg"
							onClick={() => sharePostURl(props.post.id.toString())}
						>
							<IconShare class="size-5" />
						</button>
						<Link
							class="flex size-8 flex-row items-center justify-center gap-1.5 rounded-lg"
							href={`/posts/${props.post.id}`}
						>
							<IconInfo class="size-5" />
						</Link>
					</div>
					<span class="text-xs text-hint">
						{timeSince(props.post.created_at)}
					</span>
				</div>
			</div>
		</div>
	)
}