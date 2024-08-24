import { CaloricBreakdownChart } from '~/components/calories-chart'
import { MacroPieChart } from '~/components/macros-chart'

export default function ActivityPage() {
	return (
		<section class="mt-4 p-2 text-foreground">
			<div class="w-full rounded-lg bg-section p-4">
				<CaloricBreakdownChart />
			</div>
			<div class="mt-2 w-full rounded-lg bg-section p-4">
				<MacroPieChart />
			</div>
		</section>
	)
}
