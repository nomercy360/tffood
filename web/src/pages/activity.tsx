import { useTranslations } from '~/lib/locale-context'
import { fetchFoodInsights } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { createEffect, createSignal, Show } from 'solid-js'
import { IconCarbs, IconChevronLeft, IconChevronRight, IconFats, IconProteins } from '~/components/icons'

type DaySelectorProps = {
	date: Date;
	setDate: (date: Date) => void;
}

function DaySelector(props: DaySelectorProps) {
	const [nextDayActive, setNextDayActive] = createSignal(true)

	function handlePrevDay() {
		props.setDate(new Date(props.date.setDate(props.date.getDate() - 1)))
	}

	function handleNextDay() {
		props.setDate(new Date(props.date.setDate(props.date.getDate() + 1)))
	}

	createEffect(() => {
		const today = new Date()
		const tomorrow = new Date(today)
		tomorrow.setDate(tomorrow.getDate())
		setNextDayActive(props.date.toDateString() !== tomorrow.toDateString())
	})

	return (
		<div class="flex flex-row items-center justify-center space-x-5">
			<button class="flex size-6 items-center justify-center"
				onClick={handlePrevDay}>
				<IconChevronLeft class="size-5 shrink-0 text-black" stroke="currentColor" stroke-width="2"
												 stroke-linecap="round" stroke-linejoin="round" />
			</button>
			<p class="font-semibold text-foreground">{props.date.toDateString()}</p>
			<button class="flex size-6 items-center justify-center"
				onClick={handleNextDay}
				disabled={!nextDayActive()}>
				<IconChevronRight class="size-5 shrink-0 text-black"
					stroke={nextDayActive() ? 'currentColor' : 'transparent'}
					stroke-width="2"
					stroke-linecap="round" stroke-linejoin="round" />
			</button>
		</div>
	)
}

export default function ActivityPage() {
	const { t } = useTranslations()

	const [date, setDate] = createSignal(new Date())

	const query = createQuery(() => ({
		queryKey: ['activity', date()],
		queryFn: () => fetchFoodInsights(date()),
	}))


	return (
		<section class="mt-4 px-4 py-8 text-foreground">
			<Show when={query.isSuccess} fallback={<p>Loading...</p>}>
				<div class="flex w-full flex-col items-center justify-between text-center">
					<DaySelector
						date={date()}
						setDate={setDate}
					/>
					<p class="mt-12 text-3xl font-medium">{query.data.calories_left} kcal left</p>
					<p class="mt-2.5 max-w-[305px]">
						{query.data.calories_consumed > 0 ? 'Start adding your meals and wellie will count everything. Here is yor plan for today.'
							: 'You are keeping your calories in deficit. Keep going and you’ll achieve your goal.'}
					</p>
					<div class="mt-9 flex flex-row items-center justify-center">
						<div class="flex w-20 flex-col items-center justify-center gap-3.5">
							<IconProteins width="24" height="24" />
							<div class="space-y-1">
								<p class="text-xs font-medium uppercase">Prot</p>
								<p class="text-base uppercase">{query.data.macros.proteins} g</p>
							</div>
						</div>
						<div class="flex w-20 flex-col items-center justify-center gap-3.5">
							<IconCarbs width="24" height="24" />
							<div class="space-y-1">
								<p class="text-xs font-medium uppercase">Carbs</p>
								<p class="text-base uppercase">{query.data.macros.carbohydrates} g</p>
							</div>
						</div>
						<div class="flex w-20 flex-col items-center justify-center gap-3.5">
							<IconFats width="24" height="24" />
							<div class="space-y-1">
								<p class="text-xs font-medium uppercase">Fats</p>
								<p class="text-base uppercase">{query.data.macros.fats} g</p>
							</div>
						</div>
					</div>
				</div>
			</Show>
		</section>
	)
}
