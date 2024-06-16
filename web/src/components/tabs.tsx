import { A, useLocation } from '@solidjs/router'
import { cn } from '~/lib/utils'
import { store } from '~/lib/store'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	return (
		<div class="pt-16">
			<div class="fixed inset-0 flex h-16 w-full items-center bg-background px-4">
				<div class="flex w-full flex-row items-center justify-start space-x-4">
					<ul
						class="grid w-full grid-cols-2 gap-2 rounded-xl bg-secondary p-1 text-center text-sm font-medium text-foreground"
						id="default-tab"
						role="tablist"
					>
						<li role="presentation">
							<A
								href={'/'}
								class={cn(
									'flex h-8 items-center justify-center rounded-lg bg-transparent px-4 text-sm font-medium transition-all duration-500 ease-in-out',
									location.pathname === '/' && 'bg-background',
								)}
								id="feed"
								role="tab"
								aria-controls="feed"
								aria-selected="false"
							>
								Latest
							</A>
						</li>
						<li role="presentation">
							<A
								href={'/activity'}
								class={cn(
									'flex h-8 items-center justify-center rounded-lg bg-transparent px-4 text-sm font-medium transition-all duration-500 ease-in-out',
									location.pathname === '/activity' && 'bg-background',
								)}
								id="posts-tab"
								role="tab"
								aria-controls="posts"
								aria-selected="false"
							>
								My Activity
							</A>
						</li>
					</ul>
				</div>
			</div>
			{props.children}
		</div>
	)
}
