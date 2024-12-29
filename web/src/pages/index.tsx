import { createEffect, For, Show } from 'solid-js'
import { cn, timeSince } from '~/lib/utils'
import { useNavigate } from '@solidjs/router'
import { fetchMeals } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/link'

export type Meal = {
	id: number
	user_id: number
	text: string
	created_at: string
	updated_at: string
	hidden_at: string | null
	photo_url: string
	ingredients: {
		name: string
		amount: number
		weight: number
		calories: number
	}[]
	dish_name: string
	tags: string[]
	health_rating: number
	aesthetic_rating: number
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
	const query = createQuery(() => ({
		queryKey: ['meals'],
		queryFn: fetchMeals,
	}))

	createEffect(() => {
		if (query.data) {
			console.log(query.data)
		}
	})

	const navigate = useNavigate()

	function navigateToPost() {
		navigate('/meal')
	}

	return (
		<section class="p-4">
			<div class="grid gap-2 pb-20 text-3xl text-white">
				<Show when={query.isSuccess} fallback={<Loader />}>
					<For each={query.data}>
						{(meal: Meal[]) => (
							<PostCard meal={meal} class="rounded-lg border bg-section" />
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

export function UserProfileLink(props: { class: any; user: Meal['user'] }) {
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

export function PostCard(props: { meal: Meal; class: any }) {
	function sharePostURl(mealID: string) {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/eatzfood_bot/app?startapp=p' + mealID,
			}).toString() +
			'&text=Check out this meal'

		window.Telegram.WebApp.openTelegramLink(url)
	}

	const meal = props.meal

	return (
		<div class="w-full overflow-hidden rounded-lg bg-white">
			<img
				src={meal.photo_url}
				alt={meal.dish_name}
				class="h-64 w-full object-cover"
			/>
			<div class="p-6">
				{/* Header */}
				<div class="mb-4 flex items-center space-x-4">
					<img
						src={meal.user.avatar_url}
						alt={meal.user.username}
						class="size-10 rounded-full"
					/>
					<div>
						<h4 class="text-lg font-semibold text-neutral-800">
							{meal.dish_name}
						</h4>
						<p class="text-sm text-gray-500">@{meal.user.username}</p>
					</div>
				</div>
				<div class="mb-4 flex space-x-6">
					<div class="flex items-center space-x-2">
						<span class="material-symbols-rounded text-yellow-500">star</span>
						<span class="text-neutral-800">{meal.aesthetic_rating}%</span>
					</div>
					<div class="flex items-center space-x-2">
						<span class="material-symbols-rounded text-green-500">
							health_and_safety
						</span>
						<span class="text-neutral-800">{meal.health_rating}%</span>
					</div>
				</div>
				<div class="mb-4 text-sm text-gray-700">
					<p>
						<strong>Calories:</strong> {meal.food_insights?.calories} kcal
					</p>
					<p>
						<strong>Proteins:</strong> {meal.food_insights?.proteins} g
					</p>
					<p>
						<strong>Fats:</strong> {meal.food_insights?.fats} g
					</p>
					<p>
						<strong>Carbs:</strong> {meal.food_insights?.carbohydrates} g
					</p>
				</div>
				{meal.ingredients && (
					<div>
						<h5 class="mb-2 font-semibold">Ingredients:</h5>
						<ul class="list-disc pl-6 text-sm text-gray-700">
							<For each={meal.ingredients}>
								{(ingredient) => (
									<li>
										{ingredient.name} ({ingredient.calories} kcal,{' '}
										{ingredient.weight} g)
									</li>
								)}
							</For>
						</ul>
					</div>
				)}
			</div>
		</div>
	)
}
