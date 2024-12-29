import { lazy } from 'solid-js'
import type { RouteDefinition } from '@solidjs/router'

import HomePage from '~/pages/index'
import NavigationTabs from '~/components/tabs'

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: NavigationTabs,
		children: [
			{
				component: HomePage,
				path: '',
			},
		],
	},
	{
		component: lazy(() => import('~/pages/post')),
		path: '/log',
	},
	{
		component: lazy(() => import('~/pages/meals/id')),
		path: '/meals/:id',
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
