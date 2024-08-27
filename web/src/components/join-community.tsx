import { For, onCleanup, onMount } from 'solid-js'


import { PostCard } from '~/pages'
import { useTranslations } from '~/lib/locale-context'
import { setUser, store } from '~/lib/store'
import { fetchSubmitJoinRequest, fetchUpdateUserSettings } from '~/lib/api'

const dummyData = [
	{
		'id': 14,
		'user_id': 1,
		'text': null,
		'photo_url': 'https://assets.peatch.io/media/1/5kPfCvgC.jpg',
		'user': {
			'id': 1,
			'username': 'mkkksim',
			'avatar_url': 'https://assets.peatch.io/1/927635965.jpg',
			'first_name': 'Maksim',
			'last_name': null,
			'title': null,
		},
		'dish_name': 'Датский сливочный рулет с клубникой',
	},
	{
		'id': 13,
		'user_id': 1,
		'text': null,
		'photo_url': 'https://assets.peatch.io/media/1/zJT1WRPj.jpg',
		'user': {
			'id': 1,
			'username': 'mkkksim',
			'avatar_url': 'https://assets.peatch.io/1/927635965.jpg',
			'first_name': 'Maksim',
			'last_name': null,
			'title': null,
		},
		'dish_name': 'Томатный суп',
	},
	{
		'id': 12,
		'user_id': 1,
		'text': null,
		'created_at': '2024-08-24T03:09:44Z',
		'updated_at': '2024-08-24T03:09:44Z',
		'hidden_at': '2024-08-24T03:09:44Z',
		'photo_url': 'https://assets.peatch.io/media/1/i2KjWd5w.jpg',
		'user': {
			'id': 1,
			'username': 'mkkksim',
			'avatar_url': 'https://assets.peatch.io/1/927635965.jpg',
			'first_name': 'Maksim',
			'last_name': null,
			'title': null,
		},
		'dish_name': 'Хлеб с творогом и вареньем',
	},
]


export default function JoinCommunity() {

	onMount(() => {
		const root = document.getElementById('root')
		if (root) {
			root.style.overflow = 'hidden'
		}
	})

	onCleanup(() => {
		const root = document.getElementById('root')
		if (root) {
			root.style.overflow = 'auto'
		}
	})

	const { t } = useTranslations()

	const requestToJoin = async () => {
		if (store.user.request_to_join_at !== null) {
			return
		}

		try {
			await fetchSubmitJoinRequest()
			setUser({ ...store.user, request_to_join_at: new Date().toISOString() })
		} catch (e) {
			console.error(e)
		}
	}

	return (
		<section class="relative max-h-screen overflow-y-hidden p-4 text-center">
			<h1 class="text-2xl font-bold">
				{t('common.join_community')}
			</h1>
			<p class="mt-2 text-sm text-hint">
				{t('common.join_community_description')}
			</p>
			<button
				onClick={requestToJoin}
				disabled={store.user.request_to_join_at !== null}
				class="mb-8 mt-4 h-10 w-full rounded-lg bg-primary px-2 text-sm font-medium text-white disabled:opacity-50">
				{store.user.request_to_join_at !== null ? t('common.join_community_requested') : t('common.join_community_button')}
			</button>

			<div class="grid gap-2">
				<For each={dummyData}>
					{(item) => (<PostCard post={item as any} class="pointer-events-none blur-md" />)}
				</For>
			</div>
		</section>
	)
}
