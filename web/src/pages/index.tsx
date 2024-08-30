import { For, Match, Show, Switch } from 'solid-js'
import { Navigate } from '@solidjs/router'
import { fetchPosts } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/link'
import { store } from '~/lib/store'
import Image from '~/components/image'

export type Post = {
	id: number
	user_id: number
	text: string
	created_at: string
	updated_at: string
	hidden_at: string | null
	photo_url: string
	ingredients: { name: string; weight: number }[]
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
	const isOnboarded = () => {
		return store.user.age && store.user.weight && store.user.height && store.user.fat_percentage && store.user.goal && store.user.gender
	}

	return (
		<Switch>
			<Match when={isOnboarded()}>
				<Feed />
			</Match>
			<Match when={!isOnboarded()}>
				<Navigate href="/onboard" />
			</Match>
		</Switch>
	)
}

function Feed() {
	const query = createQuery(() => ({
		queryKey: ['posts'],
		queryFn: () => fetchPosts(),
	}))

	return (
		<section class="px-1.5 py-4">
			<Show when={query.isSuccess} fallback={<Loader />}>
				<div class="grid grid-cols-2 gap-1.5">
					<div class="flex flex-col space-y-1.5">
						<For each={query.data.filter((_: any, index: number) => index % 2 === 0) as Post[]}>
							{(item) => (
								<PostCard post={item} class="" />
							)}
						</For>
					</div>
					<div class="flex flex-col space-y-1.5">
						<For each={query.data.filter((_: any, index: number) => index % 2 !== 0) as Post[]}>
							{(item) => (
								<PostCard post={item} class="" />
							)}
						</For>
					</div>
				</div>
			</Show>
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

export function PostCard(props: { post: Post; class: any }) {
	return (
		<Link href={`/posts/${props.post.id}`}>
			<div class="relative">
				<div class="absolute left-2 top-2 z-50">
					<Image src={props.post.user.avatar_url}
								 class="size-8 rounded-full" alt="Avatar"
					/>
				</div>
				<Image
					src={props.post.photo_url}
					class="h-auto w-full rounded-[20px] border object-cover"
					alt="Thumbnail"
				/>
			</div>
		</Link>
	)
}