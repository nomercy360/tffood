import { createQuery } from '@tanstack/solid-query'
import { fetchPost } from '~/lib/api'
import { For, Show } from 'solid-js'
import { Post, UserProfileLink } from '~/pages'

export default function UserProfilePage(props: any) {
	const query = createQuery(() => ({
		queryKey: ['posts', props.params.id],
		queryFn: () => fetchPost(props.params.id),
	}))

	return (
		<div class="min-h-screen bg-secondary p-2">
			<Show when={query.isSuccess} fallback={<Loading />}>
				<img
					src={query.data?.photo_url}
					class="aspect-[4/3] w-full rounded-lg object-cover"
					alt="Thumbnail"
				/>
				<div class="p-2">
					<p class="text-sm font-medium">
						{query.data?.text || query.data?.suggested_dish_name}
					</p>
					<Show when={query.data?.suggested_ingredients}>
						<p class="mt-2 text-xs text-hint">
							Ingredients:{' '}
							{query.data?.suggested_ingredients
								.map((i) => `${i.name} (${i.amount}g)`)
								.join(', ')}
						</p>
					</Show>
					<div class="mt-4 flex flex-row flex-wrap items-center justify-start gap-1.5">
						<For each={query.data?.suggested_tags}>
							{(tag) => (
								<span class="flex h-6 items-center justify-center rounded-lg bg-background px-2 py-0.5 text-xs text-hint">
									{tag}
								</span>
							)}
						</For>
					</div>
					<PostInsights insights={query.data!.food_insights} />
					<UserProfileLink
						user={query.data!.user}
						class={'mt-3 rounded-lg bg-background p-2'}
					/>
				</div>
			</Show>
		</div>
	)
}

function Loading() {
	return <p>Loading...</p>
}

function PostInsights(props: { insights: Post['food_insights'] }) {
	return (
		<div class="mt-4 flex flex-col items-start justify-start gap-2 rounded-lg bg-background p-2">
			<div class="text-sm">
				<strong>Calories:</strong> {props.insights.calories} kcal
			</div>
			<div class="text-sm">
				<strong>Proteins:</strong> {props.insights.proteins} g
			</div>
			<div class="text-sm">
				<strong>Fats:</strong> {props.insights.fats} g
			</div>
			<div class="text-sm">
				<strong>Carbohydrates:</strong> {props.insights.carbohydrates} g
			</div>
			<p class="text-sm">
				<strong>Dietary Information:</strong>{' '}
				{props.insights.dietary_information.join(', ')}
			</p>
		</div>
	)
}
