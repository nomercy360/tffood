import { store } from '~/lib/store'

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

export async function fetchPosts() {
	const response = await apiFetch({ endpoint: '/posts' })
	return response as any
}

export async function fetchCreatePost(post: any) {
	const response = await apiFetch({
		endpoint: '/posts',
		method: 'POST',
		body: post,
	})

	return response as any
}

export async function fetchUpdatePost(id: number, post: any) {
	const response = await apiFetch({
		endpoint: `/posts/${id}`,
		method: 'PUT',
		body: post,
	})

	return response as any
}

export async function fetchPostAISuggestions(id: number) {
	const response = await apiFetch({
		endpoint: `/posts/${id}/ai`,
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

export async function fetchAddPostReaction(
	id: number,
	type: 'frown' | 'meh' | 'smile',
) {
	return await apiFetch({
		endpoint: `/posts/${id}/react/${type}`,
		method: 'POST',
	})
}

export async function fetchRemovePostReaction(id: number) {
	return await apiFetch({
		endpoint: `/posts/${id}/react`,
		method: 'DELETE',
	})
}

export async function fetchTags() {
	const response = await apiFetch({ endpoint: '/tags' })
	return response as any
}
