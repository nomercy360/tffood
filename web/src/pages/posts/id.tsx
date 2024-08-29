import { createQuery } from '@tanstack/solid-query'
import { fetchPost } from '~/lib/api'
import { For, Show } from 'solid-js'
import { Post } from '~/pages'
import { useTranslations } from '~/lib/locale-context'
import { IconCarbs, IconFats, IconProteins } from '~/components/icons'

export default function UserProfilePage(props: any) {
	const query = createQuery(() => ({
		queryKey: ['posts', props.params.id],
		queryFn: () => fetchPost(props.params.id),
	}))

	const { t } = useTranslations()

	return (
		<div class="relative min-h-screen pb-10">
			<Show when={query.isSuccess} fallback={<Loading />}>
				<div
					class="absolute inset-x-0 top-0 flex w-full flex-row items-center justify-between bg-white/5 px-6 py-5 backdrop-blur">
					<div class="flex flex-row items-center space-x-2.5">
						<img
							src={query.data?.user.avatar_url}
							class="size-6 rounded-full"
							alt="Avatar"
						/>
						<span class="text-base font-medium text-white">
							{query.data?.user.first_name} {query.data?.user.last_name}
						</span>
					</div>
					<div class="flex flex-row items-center justify-center space-x-5">
						<button>
							<svg width="12" height="18" viewBox="0 0 12 18" fill="none" xmlns="http://www.w3.org/2000/svg"
									 class="shrink-0">
								<path
									d="M1.28516 17.897C1.00846 17.897 0.792643 17.8084 0.637695 17.6313C0.482747 17.4543 0.405273 17.2052 0.405273 16.8843V2.47412C0.405273 1.68278 0.601725 1.08789 0.994629 0.689453C1.38753 0.291016 1.97412 0.0917969 2.75439 0.0917969H9.24561C10.0259 0.0917969 10.6125 0.291016 11.0054 0.689453C11.3983 1.08789 11.5947 1.68278 11.5947 2.47412V16.8843C11.5947 17.2052 11.5173 17.4543 11.3623 17.6313C11.2074 17.8084 10.9915 17.897 10.7148 17.897C10.5101 17.897 10.3192 17.8278 10.1421 17.6895C9.97054 17.5511 9.69661 17.3021 9.32031 16.9424L6.07471 13.7466C6.0249 13.6912 5.9751 13.6912 5.92529 13.7466L2.67969 16.9424C2.30339 17.3076 2.02669 17.5566 1.84961 17.6895C1.67253 17.8278 1.48438 17.897 1.28516 17.897Z"
									fill="white" />
							</svg>
						</button>
						<button>
							<svg width="17" height="16" viewBox="0 0 17 16" fill="none" xmlns="http://www.w3.org/2000/svg"
									 class="shrink-0">
								<path
									d="M8.5 15.8843C8.41699 15.8843 8.31738 15.8594 8.20117 15.8096C8.09049 15.7653 7.99089 15.7155 7.90234 15.6602C6.3418 14.6641 4.98877 13.6209 3.84326 12.5308C2.69775 11.4351 1.81234 10.3117 1.18701 9.16064C0.56722 8.00407 0.257324 6.83643 0.257324 5.65771C0.257324 4.92171 0.373535 4.24935 0.605957 3.64062C0.843913 3.02637 1.17318 2.49512 1.59375 2.04688C2.01432 1.59863 2.50407 1.25277 3.06299 1.00928C3.62191 0.765788 4.22786 0.644043 4.88086 0.644043C5.69434 0.644043 6.4082 0.851562 7.02246 1.2666C7.63672 1.68164 8.12923 2.22949 8.5 2.91016C8.8763 2.22396 9.37158 1.67611 9.98584 1.2666C10.6001 0.851562 11.3112 0.644043 12.1191 0.644043C12.7721 0.644043 13.3781 0.765788 13.937 1.00928C14.5015 1.25277 14.9912 1.59863 15.4062 2.04688C15.8268 2.49512 16.1533 3.02637 16.3857 3.64062C16.6237 4.24935 16.7427 4.92171 16.7427 5.65771C16.7427 6.83643 16.43 8.00407 15.8047 9.16064C15.1849 10.3117 14.3022 11.4351 13.1567 12.5308C12.0168 13.6209 10.6665 14.6641 9.10596 15.6602C9.01742 15.7155 8.91504 15.7653 8.79883 15.8096C8.68815 15.8594 8.58854 15.8843 8.5 15.8843Z"
									fill="white" />
							</svg>
						</button>
					</div>
				</div>
				<img
					src={query.data?.photo_url}
					class="h-auto w-full object-cover"
					alt="Thumbnail"
				/>
				<div class="px-4">
					<div class="mt-4 flex flex-row flex-wrap items-center justify-start gap-1.5">
						<For each={query.data?.tags}>
							{(tag) => (
								<span
									class="flex h-6 items-center justify-center rounded-lg bg-background px-2 py-0.5 text-xs text-hint">
									{tag}
								</span>
							)}
						</For>
					</div>
					<p class="text-3xl font-medium">
						{query.data?.text || query.data?.dish_name}
					</p>
					<PostInsights insights={query.data!.food_insights} />
					<IngredientDetails ingredients={query.data!.ingredients} />
				</div>
			</Show>
		</div>
	)
}

function Loading() {
	return <p>Loading...</p>
}

function PostInsights(props: { insights: Post['food_insights'] }) {
	const { t } = useTranslations()
	return (
		<div class="mt-8 flex flex-col items-start justify-start rounded-xl bg-secondary px-5 pb-5 pt-4">
			<p class="text-xs font-medium uppercase">
				{t('common.nutrients')}
			</p>
			<p class="mt-1 text-2xl font-medium">
				{props.insights.calories} {t('common.kcal')}
			</p>
			<div class="mt-7 grid w-full gap-4">
				<div class="flex w-full flex-row items-center justify-between">
					<div class="flex flex-row items-center justify-start space-x-2.5">
						<IconProteins width="24" height="24" />
						<p>
							{t('common.proteins')}
						</p>
					</div>
					<p>
						{props.insights.proteins} {t('common.g')}
					</p>
				</div>
				<div class="flex w-full flex-row items-center justify-between">
					<div class="flex flex-row items-center justify-start space-x-2.5">
						<IconCarbs width="24" height="24" />
						<p>
							{t('common.carbohydrates')}
						</p>
					</div>
					<p>
						{props.insights.carbohydrates} {t('common.g')}
					</p>
				</div>
				<div class="flex w-full flex-row items-center justify-between">
					<div class="flex flex-row items-center justify-start space-x-2.5">
						<IconFats width="24" height="24" />
						<p>
							{t('common.fats')}
						</p>
					</div>
					<p>
						{props.insights.fats} {t('common.g')}
					</p>
				</div>
			</div>
		</div>
	)
}

function IngredientDetails(props: { ingredients: Post['ingredients'] }) {
	const { t } = useTranslations()

	return (
		<div class="mt-2 flex flex-col items-start justify-start rounded-xl bg-secondary px-5 pb-5 pt-4">
			<p class="text-xs font-medium uppercase">
				{t('common.ingredients')}
			</p>
			<div class="mt-7 grid w-full gap-4">
				<For each={props.ingredients}>
					{(ingredient) => (
						<div class="flex w-full flex-row items-center justify-between">
							<p>
								{ingredient.name}
							</p>
							<p>
								{ingredient.weight} {t('common.g')}
							</p>
						</div>
					)}
				</For>
			</div>
		</div>
	)
}
