import { CaloricBreakdownChart } from '~/components/calories-chart'
import { MacroPieChart } from '~/components/macros-chart'

export default function ActivityPage() {
	return (
		<section class="p-8 text-foreground">
			<h1>Nutrition Dashboard</h1>
			<CaloricBreakdownChart />
			<MacroPieChart />
		</section>
	)
}
