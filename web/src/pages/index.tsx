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

	let lastUserId: number | null = null

	return (
		<section class="px-1.5 py-4">
			<Show when={query.isSuccess} fallback={<Loader />}>
				<div class="grid grid-cols-2 gap-1.5">
					<For each={[0, 1]}>
						{(index) => (
							<div class="flex flex-col space-y-1.5">
								<For
									each={query.data.filter((_: any, i: number) => (i % 2 === index)) as Post[]}>
									{(item, _) => {
										const showAvatar = lastUserId !== item.user.id
										lastUserId = item.user.id
										return (
											<PostCard post={item} showAvatar={showAvatar} />
										)
									}}
								</For>
							</div>
						)}
					</For>
				</div>
			</Show>
		</section>
	)
}

function Loader() {
	return (
		<div class="grid grid-cols-2 gap-1.5">
			<div class="flex flex-col space-y-1.5">
				<div class="h-40 animate-pulse rounded-[20px] bg-border" />
				<div class="h-80 animate-pulse rounded-[20px] bg-border" />
				<div class="h-96 animate-pulse rounded-[20px] bg-border" />
				<div class="h-48 animate-pulse rounded-[20px] bg-border" />
			</div>
			<div class="flex flex-col space-y-1.5">
				<div class="h-96 animate-pulse rounded-[20px] bg-border" />
				<div class="h-72 animate-pulse rounded-[20px] bg-border" />
				<div class="h-96 animate-pulse rounded-[20px] bg-border" />
				<div class="h-96 animate-pulse rounded-[20px] bg-border" />
			</div>
		</div>
	)
}

export function PostCard(props: { post: Post; showAvatar: boolean }) {
	return (
		<Link href={`/posts/${props.post.id}`}>
			<div class="relative">
				<Show when={props.showAvatar}>
					<div class="absolute left-2 top-2 z-30">
						<Image src={props.post.user.avatar_url}
									 class="size-8 rounded-full" alt="Avatar"
						/>
					</div>
				</Show>
				<Image
					src={props.post.photo_url}
					class="h-auto w-full rounded-[20px] border object-cover"
					alt="Thumbnail"
				/>
			</div>
		</Link>
	)
}
