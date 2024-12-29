import { createEffect, createSignal, Match, Switch } from 'solid-js'
import { setToken, setUser } from '~/lib/store'
import { API_BASE_URL } from '~/lib/api'
import { NavigationProvider, useNavigation } from '~/lib/useNavigation'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import { useLocation, useNavigate } from '@solidjs/router'
import Toast from '~/components/toast'
import { LocaleContextProvider } from '~/lib/locale-context'

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 1,
			staleTime: 1000 * 60 * 5, // 5 minutes
			gcTime: 1000 * 60 * 5, // 5 minutes
		},
		mutations: {
			retry: 1,
		},
	},
})

function transformStartParam(startParam?: string): string | null {
	if (!startParam) return null

	// Check if the parameter starts with "redirect-to-"
	if (startParam.startsWith('p')) {
		const path = startParam.slice('p'.length)

		return '/meals/' + path
	} else {
		return null
	}
}

export default function App(props: any) {
	const [isAuthenticated, setIsAuthenticated] = createSignal(false)
	const [isLoading, setIsLoading] = createSignal(true)

	const navigate = useNavigate()

	createEffect(async () => {
		const initData = window.Telegram.WebApp.initData

		console.log('WEBAPP:', window.Telegram)

		try {
			const resp = await fetch(`${API_BASE_URL}/auth/telegram?` + initData, {
				method: 'POST',
			})

			if (!resp.ok) {
				throw new Error('Failed to authenticate user')
			}

			const { user, token } = await resp.json()

			setUser(user)
			setToken(token)

			window.Telegram.WebApp.ready()
			window.Telegram.WebApp.expand()
			console.log('WEBAPP:', window.Telegram)
			window.Telegram.WebApp.SettingsButton.show()
			window.Telegram.WebApp.SettingsButton.onClick(() => {
				navigate('/settings')
			})

			setIsAuthenticated(true)
			setIsLoading(false)

			// if there is a redirect url, redirect to it
			// ?startapp=redirect-to=/users/

			const startapp = window.Telegram.WebApp.initDataUnsafe.start_param

			const redirectUrl = transformStartParam(startapp)

			if (redirectUrl) {
				navigate(redirectUrl)
				return
			}
		} catch (e) {
			console.error('Failed to authenticate user:', e)
			setIsAuthenticated(false)
			setIsLoading(false)
		}
	})

	return (
		<LocaleContextProvider>
			<NavigationProvider>
				<QueryClientProvider client={queryClient}>
					<Switch>
						<Match when={isAuthenticated()}>
							<div>{props.children}</div>
						</Match>
						<Match when={!isAuthenticated() && isLoading()}>
							<div class="h-screen w-full flex-col items-start justify-center bg-secondary" />
						</Match>
						<Match when={!isAuthenticated() && !isLoading()}>
							<div class="h-screen min-h-screen w-full flex-col items-start justify-center bg-secondary text-foreground">
								Something went wrong. Please try again later.
							</div>
						</Match>
					</Switch>
					<Toast />
				</QueryClientProvider>
			</NavigationProvider>
		</LocaleContextProvider>
	)
}
