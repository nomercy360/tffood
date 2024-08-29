import { createEffect, createSignal, For, Match, onCleanup, onMount, Show, Switch } from 'solid-js'
import { cn } from '~/lib/utils'
import { useBackButton } from '~/lib/useBackButton'
import { useTranslations } from '~/lib/locale-context'
import { createStore } from 'solid-js/store'
import { useNavigate } from '@solidjs/router'
import { fetchSaveUserOnboardData } from '~/lib/api'
import { setUser, store } from '~/lib/store'

type LayoutProps = {
	title: string;
	subtitle: string;
	step: number;
	children: any;
}

function ScreenLayout(props: LayoutProps) {
	return (
		<div class="flex h-screen flex-col items-center justify-between py-6">
			<div class="space-y-1.5 px-10 text-center">
				<p class="text-lg font-medium">{props.title}</p>
				<p class="text-sm leading-normal">{props.subtitle}</p>
			</div>
			{props.children}
		</div>
	)
}

type NumberInputProps = {
	value: string;
	setValue: (value: string) => void;
	units?: string;
}

function NumberInput(props: NumberInputProps) {
	const handleClick = (num: number) => {
		const val = props.value.length < 3 ? props.value + num : props.value
		props.setValue(val)
	}

	const handleDelete = () => {
		const val = props.value.slice(0, -1)
		props.setValue(val)
	}

	return (
		<div class="flex h-screen w-full flex-col items-center justify-center">
			<div class="mb-10 text-[80px]">
				{props.value || '0'}
				<Show when={props.units}>
					<span class="text-lg font-semibold">{props.units}</span>
				</Show>
			</div>
			<div class="grid w-full grid-cols-3 gap-4">
				<For each={[1, 2, 3, 4, 5, 6, 7, 8, 9, 0]}>{(num) => (
					<button
						type="button"
						onClick={() => handleClick(num)}
						class="flex h-12 w-full items-center justify-center text-lg font-medium focus:outline-none"
					>
						{num}
					</button>
				)}</For>
				<button
					type="button"
					onClick={handleDelete}
					class="flex h-12 w-full items-center justify-center text-xl focus:outline-none"
				>
					<svg width="21" height="18" viewBox="0 0 21 18" fill="none" xmlns="http://www.w3.org/2000/svg"
							 class={cn('text-black', props.value.length > 0 ? 'text-black' : 'text-secondary-foreground')}>
						<path
							d="M17.5044 17.4761L9.48877 17.4668C9.0682 17.4668 8.68473 17.4235 8.33838 17.3369C7.99821 17.2565 7.6766 17.1235 7.37354 16.938C7.07048 16.7586 6.77979 16.5174 6.50146 16.2144L1.24121 10.6572C0.907227 10.2923 0.666016 9.96452 0.517578 9.67383C0.375326 9.38314 0.304199 9.09245 0.304199 8.80176C0.304199 8.60384 0.335124 8.41211 0.396973 8.22656C0.458822 8.03483 0.55778 7.83382 0.693848 7.62354C0.829915 7.41325 1.01237 7.1875 1.24121 6.94629L6.49219 1.37061C6.77051 1.06755 7.0612 0.826335 7.36426 0.646973C7.66732 0.467611 7.98893 0.34082 8.3291 0.266602C8.67546 0.192383 9.05892 0.155273 9.47949 0.155273H17.5044C18.5063 0.155273 19.264 0.408854 19.7773 0.916016C20.2907 1.42318 20.5474 2.17464 20.5474 3.17041V14.5352C20.5474 15.5247 20.2907 16.2607 19.7773 16.7432C19.264 17.2318 18.5063 17.4761 17.5044 17.4761ZM9.08984 12.7354C9.34342 12.7354 9.5599 12.6519 9.73926 12.4849L12.1421 10.0635L14.5542 12.4849C14.7212 12.6519 14.9315 12.7354 15.1851 12.7354C15.4325 12.7354 15.6396 12.6519 15.8066 12.4849C15.9798 12.3117 16.0664 12.1045 16.0664 11.8633C16.0664 11.6097 15.9798 11.4025 15.8066 11.2417L13.3853 8.82031L15.8159 6.39893C15.9891 6.21956 16.0757 6.01237 16.0757 5.77734C16.0757 5.53613 15.9891 5.33203 15.8159 5.16504C15.6489 4.99186 15.4448 4.90527 15.2036 4.90527C14.9624 4.90527 14.7552 4.99186 14.582 5.16504L12.1421 7.58643L9.71143 5.16504C9.53825 5.00423 9.33105 4.92383 9.08984 4.92383C8.84863 4.92383 8.64144 5.00732 8.46826 5.17432C8.30127 5.33512 8.21777 5.54232 8.21777 5.7959C8.21777 6.02474 8.30436 6.22884 8.47754 6.4082L10.8989 8.82031L8.47754 11.251C8.30436 11.418 8.21777 11.6221 8.21777 11.8633C8.21777 12.1045 8.30127 12.3117 8.46826 12.4849C8.64144 12.6519 8.84863 12.7354 9.08984 12.7354Z"
							fill="currentColor" />
					</svg>
				</button>
			</div>
		</div>
	)
}

export const [onboardStore, setOnboardStore] = createStore<{
	age: string;
	weight: string;
	height: string;
	bodyFat: string;
	goal: string;
	gender: string;
}>({
	age: '',
	weight: '',
	height: '',
	bodyFat: '',
	goal: '',
	gender: 'male',
})

export default function OnboardPage() {
	const [step, setStep] = createSignal(1)

	const backButton = useBackButton()

	const decrementStep = () => {
		setStep(step() - 1)
	}

	onMount(() => {
		backButton.onClick(decrementStep)
	})

	onCleanup(() => {
		backButton.offClick(decrementStep)
	})

	createEffect(() => {
		if (step() > 1) {
			backButton.setVisible()
		} else {
			backButton.hide()
		}
	})

	const { t } = useTranslations()

	const navigate = useNavigate()

	const texts = [
		{
			title: t('onboard.gender.question'),
			subtitle: t('onboard.gender.description'),
		},
		{
			title: t('onboard.age.question'),
			subtitle: t('onboard.age.description'),
		},
		{
			title: t('onboard.weight.question'),
			subtitle: t('onboard.weight.description'),
		},
		{
			title: t('onboard.height.question'),
			subtitle: t('onboard.height.description'),
		},
		{
			title: t('onboard.body_fat.question'),
			subtitle: t('onboard.body_fat.description'),
		},
		{
			title: t('onboard.goals.question'),
			subtitle: t('onboard.goals.description'),
		},
	]

	function getButtonState() {
		if (step() === 1) {
			return onboardStore.gender.length > 0
		} else if (step() === 2) {
			return onboardStore.age.length > 0
		} else if (step() === 3) {
			return onboardStore.weight.length > 0
		} else if (step() === 4) {
			return onboardStore.height.length > 0
		} else if (step() === 5) {
			return onboardStore.bodyFat.length > 0
		} else if (step() === 6) {
			return onboardStore.goal.length > 0
		}
	}

	const goals = [
		{
			title: t('onboard.goals.options.gain_muscles.title'),
			description: t('onboard.goals.options.gain_muscles.description'),
			value: 'gain_muscles',
		},
		{
			title: t('onboard.goals.options.lose_weight.title'),
			description: t('onboard.goals.options.lose_weight.description'),
			value: 'lose_weight',
		},
		{
			title: t('onboard.goals.options.track_nutrition.title'),
			description: t('onboard.goals.options.track_nutrition.description'),
			value: 'track_nutrition',
		},
		{
			title: t('onboard.goals.options.improve_health.title'),
			description: t('onboard.goals.options.improve_health.description'),
			value: 'improve_health',
		},
		{
			title: t('onboard.goals.options.count_calories.title'),
			description: t('onboard.goals.options.count_calories.description'),
			value: 'count_calories',
		},
	]

	async function saveUserOnboardData() {
		try {
			const user = await fetchSaveUserOnboardData({
				age: Number(onboardStore.age),
				weight: Number(onboardStore.weight),
				height: Number(onboardStore.height),
				fat_percentage: Number(onboardStore.bodyFat),
				goal: onboardStore.goal,
				gender: onboardStore.gender
			})

			setUser(user)
			navigate('/')

		} catch (e) {
			console.error(e)
		}
	}

	return (
		<ScreenLayout
			title={texts[step() - 1].title}
			subtitle={texts[step() - 1].subtitle}
			step={step()}>
			<Switch>
				<Match when={step() === 1}>
					<div class="flex w-full flex-col space-y-1.5 px-1.5">
						<label
							class={cn('flex flex-col space-y-1 rounded-2xl border border-gray-200 px-5 py-4 text-start', 'male' === onboardStore.gender && 'bg-primary-foreground')}>
							<input
								type="radio"
								name="gender"
								value="male"
								class="hidden"
								checked={'male' === onboardStore.gender}
								onChange={(e) => setOnboardStore('gender', e.currentTarget.value)}
							/>
							<p class="text-lg font-medium">Male</p>
							<p class="text-sm text-secondary-foreground">Wishd is made for friends. Search for people you know or
								share your link, so they can find you.</p>
						</label>
						<label
							class={cn('flex flex-col space-y-1 rounded-2xl border border-gray-200 px-5 py-4 text-start', 'female' === onboardStore.gender && 'bg-primary-foreground')}>
							<input
								type="radio"
								name="gender"
								value="female"
								class="hidden"
								checked={'female' === onboardStore.gender}
								onChange={(e) => setOnboardStore('gender', e.currentTarget.value)}
							/>
							<p class="text-lg font-medium">Female</p>
							<p class="text-sm text-secondary-foreground">Wishd is made for friends. Search for people you know or
								share your link, so they can find you.</p>
						</label>
					</div>
				</Match>
				<Match when={step() === 2}>
					<NumberInput
						value={onboardStore.age}
						setValue={(value) => setOnboardStore('age', value)} />
				</Match>
				<Match when={step() === 3}>
					<NumberInput
						units="kg"
						value={onboardStore.weight}
						setValue={(value) => setOnboardStore('weight', value)} />
				</Match>
				<Match when={step() === 4}>
					<NumberInput
						units="cm"
						value={onboardStore.height}
						setValue={(value) => setOnboardStore('height', value)} />
				</Match>
				<Match when={step() === 5}>
					<NumberInput
						units="%"
						value={onboardStore.bodyFat}
						setValue={(value) => setOnboardStore('bodyFat', value)} />
				</Match>
				<Match when={step() === 6}>
					<div class="flex w-full flex-col space-y-1.5 px-1.5">
						<For each={goals}>{(item) => (
							<label
								class={cn('flex flex-col space-y-1 rounded-2xl border border-gray-200 px-5 py-4 text-start', item.value === onboardStore.goal && 'bg-primary-foreground')}>
								<input
									type="radio"
									name="goal"
									value={item.value}
									class="hidden"
									checked={item.value === onboardStore.goal}
									onChange={(e) => setOnboardStore('goal', e.currentTarget.value)}
								/>
								<p class="text-lg font-medium">{item.title}</p>
								<p class="text-sm text-secondary-foreground">{item.description}</p>
							</label>
						)}</For>
					</div>
				</Match>
			</Switch>
			<button
				onClick={() => step() === 6 ? saveUserOnboardData() : setStep(step() + 1)}
				class={cn('h-11 w-[150px] rounded-full bg-black text-base font-semibold text-white focus:outline-none', !getButtonState() && 'opacity-50')}
				disabled={!getButtonState()}
			>
				{step() === 6 ? t('onboard.save_button') : step() === 1 ? 'Continue as ' + onboardStore.gender : t('onboard.continue_button')}
			</button>
		</ScreenLayout>
	)
}