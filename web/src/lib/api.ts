import { store } from '~/lib/store'
import { Meal } from '~/pages'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string

export const apiFetch = async ({
	endpoint,
	method = 'GET',
	body = null,
	showProgress = true,
	contentType = 'application/json',
	responseContentType = 'json' as 'json' | 'blob',
}: {
	endpoint: string
	method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
	body?: unknown
	showProgress?: boolean
	responseContentType?: string
	contentType?: string
	token?: string
}) => {
	const headers: { [key: string]: string } = {}

	headers.Authorization = `Bearer ${store.token}`

	let bodyToSend = body
	if (contentType === 'application/json' && body) {
		bodyToSend = JSON.stringify(body)
		headers['Content-Type'] = 'application/json'
	} else if (contentType === 'multipart/form-data' && body) {
		bodyToSend = new FormData()
		Object.entries(body).forEach(([key, value]) => {
			;(bodyToSend as FormData).append(key, value)
		})
	} else {
		bodyToSend = undefined
	}

	const response = await fetch(`${API_BASE_URL}/api${endpoint}`, {
		method,
		headers,
		body: bodyToSend as BodyInit,
	})

	if (!response.ok) {
		const errorResponse = await response.json()
		throw { code: response.status, message: errorResponse.message }
	}

	switch (response.status) {
		case 204:
			return true
		default:
			return response[responseContentType as 'json' | 'blob']()
	}
}

export async function fetchMeals() {
	const response = await apiFetch({ endpoint: '/meals' })
	return response as any
}

export async function fetchFoodInsights() {
	const response = await apiFetch({ endpoint: '/food-insights' })
	return response as any
}

export async function fetchPost(id: number) {
	const response = await apiFetch({ endpoint: `/meals/${id}` })
	return response as Meal
}

export async function fetchUpdateUserSettings(settings: any) {
	const response = await apiFetch({
		endpoint: '/user/settings',
		method: 'PUT',
		body: settings,
	})

	return response as any
}

export async function fetchCreatePost(post: any) {
	const response = await apiFetch({
		endpoint: '/meals',
		method: 'POST',
		body: post,
	})

	return response as any
}

export async function fetchUpdatePost(id: number, post: any) {
	const response = await apiFetch({
		endpoint: `/meals/${id}`,
		method: 'PUT',
		body: post,
	})

	return response as any
}

export async function fetchPostAISuggestions(id: number) {
	const response = await apiFetch({
		endpoint: `/meals/${id}/ai`,
	})

	return response as any
}

export async function fetchPresignedUrl(filename: string) {
	const response = await apiFetch({
		endpoint: '/presigned-url',
		method: 'POST',
		body: { file_name: filename },
	})

	return response as any
}

export async function fetchTags() {
	const response = await apiFetch({ endpoint: '/tags' })
	return response as any
}

export async function fetchSubmitJoinRequest() {
	await apiFetch({ endpoint: '/community/join', method: 'POST' })
}
