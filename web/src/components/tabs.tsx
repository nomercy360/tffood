import { For } from 'solid-js'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { Link } from '~/components/link'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{
			href: '/log',
			icon: 'dinner_dining',
			name: '+Meal',
		},
		{
			href: '/',
			icon: 'home',
			name: 'Home',
		},
		{
			href: '/friends',
			icon: 'group',
			name: 'Friends',
		},
	]

	return (
		<>
			<div class="fixed bottom-0 z-50 grid h-[72px] w-full grid-cols-3 items-center border bg-background pb-2 shadow-sm">
				<For each={tabs}>
					{({ href, icon, name }) => (
						<Link
							href={href}
							state={{ from: location.pathname }}
							class={cn(
								'flex h-12 flex-col items-center justify-between text-sm text-secondary-foreground',
								{
									'text-foreground': location.pathname === href,
								},
							)}
						>
							<span class="material-symbols-rounded text-[32px]">{icon}</span>
							<span class="text-xs">{name}</span>
						</Link>
					)}
				</For>
			</div>
			{props.children}
		</>
	)
}
