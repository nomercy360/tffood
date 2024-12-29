import { createEffect, createSignal, onCleanup, onMount } from 'solid-js'
import {
	BarController,
	BarElement,
	CategoryScale,
	Chart,
	ChartItem,
	Colors,
	Filler,
	Legend,
	LinearScale,
	Title,
	Tooltip,
} from 'chart.js'
import { createQuery } from '@tanstack/solid-query'
import { fetchFoodInsights, fetchMeals } from '~/lib/api'
import { Meal } from '~/pages'
import { useTranslations } from '~/lib/locale-context'

Chart.register(
	Colors,
	Filler,
	Legend,
	Tooltip,
	CategoryScale,
	LinearScale,
	BarElement,
	BarController,
	Title,
)

const CaloricBreakdownChart = () => {
	const [canvasRef, setCanvasRef] = createSignal<HTMLCanvasElement | null>()
	const [chartData, setChartData] = createSignal()

	const query = createQuery(() => ({
		queryKey: ['food-insights'],
		queryFn: () => fetchFoodInsights(),
	}))

	const { t } = useTranslations()

	createEffect(() => {
		if (!query.data) return

		setChartData({
			labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
			datasets: [
				{
					label: 'Calories',
					data: query.data.caloric_breakdown,
					backgroundColor: window.Telegram.WebApp.themeParams.accent_text_color,
					borderRadius: 10,
				},
			],
		})

		const ctx = canvasRef()?.getContext('2d') as ChartItem
		const chart = new Chart(ctx, {
			type: 'bar',
			data: chartData() as any,
			options: {
				responsive: true,
				plugins: {
					legend: {
						display: false,
					},
					title: {
						display: true,
						text: t('common.caloric_breakdown_chart_title'),
					},
				},
				scales: {
					x: {
						stacked: true,
						grid: {
							display: false,
						},
					},
					y: {
						grid: {
							display: false,
						},
						stacked: true,
					},
				},
			},
		})

		onCleanup(() => {
			chart.destroy()
		})
	})

	return <canvas ref={(el) => setCanvasRef(el)} />
}

export { CaloricBreakdownChart }
