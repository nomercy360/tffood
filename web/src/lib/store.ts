import { createStore } from 'solid-js/store'

type User = {
	id: number
	username: string
	language: string
	notifications_enabled: boolean
}

type AuthStore = {
	user: User
	token: string
	showSubmitAppPopup: boolean | undefined
}

export const [store, setStore] = createStore<{
	user: User
	token: string
	showSubmitAppPopup: boolean | undefined
}>({} as AuthStore)

export const setUser = (user: User) => setStore('user', user)

export const setToken = (token: string) => setStore('token', token)

export const setShowSubmitAppPopup = (showSubmitAppPopup: boolean) =>
	setStore('showSubmitAppPopup', showSubmitAppPopup)
