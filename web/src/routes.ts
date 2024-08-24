import { lazy } from 'solid-js'
import type { RouteDefinition } from '@solidjs/router'

import HomePage from '~/pages/index'
import NavigationTabs from '~/components/tabs'
import ActivityPage from '~/pages/activity'
import JoinCommunity from '~/components/join-community'

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: NavigationTabs,
		children: [
			{
				component: HomePage,
				path: '',
			},
			{
				component: ActivityPage,
				path: '/activity',
			},
		],
	},
	{
		component: lazy(() => import('~/pages/post')),
		path: '/post',
	},
	{
		component: lazy(() => import('~/pages/posts/id')),
		path: '/posts/:id',
	},
	{
		component: lazy(() => import('~/pages/users/handle')),
		path: '/users/:username',
	},
	{
		component: lazy(() => import('~/pages/settings')),
		path: '/settings',
	},
	{
		path: '**',
		component: lazy(() => import('./pages/404')),
	},
]
