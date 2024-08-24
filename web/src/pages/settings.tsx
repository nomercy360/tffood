import { createSignal, onCleanup, onMount } from 'solid-js'
import { useMainButton } from '~/lib/useMainButton'
import { store } from '~/lib/store'
import { fetchUpdateUserSettings } from '~/lib/api'
import { useTranslations } from '~/lib/locale-context'

export default function SettingsPage() {
	const [notificationsEnabled, setNotificationsEnabled] = createSignal(store.user.notifications_enabled)
	const [language, setLanguage] = createSignal(store.user.language)

	const { t } = useTranslations()

	const mutate = async () => {
		await fetchUpdateUserSettings({
			notifications_enabled: notificationsEnabled(),
			language: language(),
		})
	}

	const mainButton = useMainButton()

	onMount(async () => {
		mainButton.setVisible(t('common.save_changes'))
		mainButton.onClick(mutate)
	})

	onCleanup(() => {
		mainButton.hide()
			.offClick(mutate)
	})


	return (
		<form class="mx-auto min-h-screen space-y-4 bg-white p-6">
			<label class="inline-flex w-full cursor-pointer items-center justify-between">
				<span class="text-sm font-medium text-foreground dark:text-secondary">
					{t('common.notifications_enabled')}
				</span>
				<input
					type="checkbox"
					checked={notificationsEnabled()}
					onChange={() => setNotificationsEnabled(!notificationsEnabled())}
					class="peer sr-only" />
				<div
					class="peer relative h-6 w-11 rounded-full bg-background after:absolute after:start-[2px] after:top-[2px] after:size-5 after:rounded-full after:border after:bg-white after:transition-all after:content-[''] peer-checked:bg-accent-foreground peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none rtl:peer-checked:after:-translate-x-full" />
			</label>
			<div class="space-y-1">
				<label for="language" class="block text-sm font-medium text-foreground">
					{t('common.language')}
				</label>
				<select
					id="language"
					value={store.user.language}
					onChange={(e) => setLanguage(e.target.value)}
					class="block h-10 w-full rounded-md border bg-white px-4"
				>
					<option value="en">
						{t('common.english')}
					</option>
					<option value="ru">
						{t('common.russian')}
					</option>
				</select>
			</div>
		</form>
	)
}
