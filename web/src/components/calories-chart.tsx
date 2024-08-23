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
import { fetchPosts } from '~/lib/api'
import { Post } from '~/pages'

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
		queryKey: ['posts', 'caloric-macro-breakdown'],
		queryFn: () => fetchPosts(),
	}))

	const getDayName = (dateStr: string) => {
		const date = new Date(dateStr)
		return date.toLocaleDateString('en-US', { weekday: 'short' })
	}

	const transformData = (apiData: Post[]) => {
		const dataByDay = {
			Mon: 0, Tue: 0, Wed: 0, Thu: 0, Fri: 0, Sat: 0, Sun: 0,
		}

		apiData.forEach((entry) => {
			const day = getDayName(entry.created_at)
			if (dataByDay.hasOwnProperty(day)) {
				// @ts-ignore
				dataByDay[day] += entry.food_insights.calories
			}
		})

		return Object.values(dataByDay)
	}

	createEffect(() => {
		if (!query.data) return

		setChartData({
			labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
			datasets: [
				{
					label: 'Calories',
					data: transformData(query.data),
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
						text: 'Daily Caloric Breakdown',
					},
				},
				scales: {
					x: {
						stacked: true,
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
