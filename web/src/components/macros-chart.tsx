import { createEffect, createSignal, onCleanup, onMount } from 'solid-js'
import {
	Chart,
	ChartItem,
	PieController,
	ArcElement,
	Tooltip,
	Legend,
} from 'chart.js'
import { createQuery } from '@tanstack/solid-query'
import { fetchPosts } from '~/lib/api'
import { Post } from '~/pages'

Chart.register(PieController, ArcElement, Tooltip, Legend)

const MacroPieChart = () => {
	const [canvasRef, setCanvasRef] = createSignal<HTMLCanvasElement | null>()
	const [chartData, setChartData] = createSignal()

	const query = createQuery(() => ({
		queryKey: ['posts', 'macro-pie-chart'],
		queryFn: () => fetchPosts(),
	}))

	const aggregateMacroData = (apiData: Post[]) => {
		let today = new Date().toISOString().split('T')[0]
		let macros = { proteins: 0, fats: 0, carbohydrates: 0 }

		apiData.forEach((entry) => {
			let entryDate = new Date(entry.created_at).toISOString().split('T')[0]
			if (entryDate === today) {
				macros.proteins += entry.food_insights.proteins || 0
				macros.fats += entry.food_insights.fats || 0
				macros.carbohydrates += entry.food_insights.carbohydrates || 0
			}
		})

		return macros
	}

	createEffect(() => {
		if (!query.data) return

		const macros = aggregateMacroData(query.data)

		setChartData({
			labels: ['Proteins', 'Fats', 'Carbohydrates'],
			datasets: [
				{
					data: [macros.proteins, macros.fats, macros.carbohydrates],
					backgroundColor: [
						'#FF6384', // Red for Proteins
						'#36A2EB', // Blue for Fats
						'#FFCE56', // Yellow for Carbohydrates
					],
					hoverBackgroundColor: [
						'#FF6384',
						'#36A2EB',
						'#FFCE56',
					],
				},
			],
		})

		const ctx = canvasRef()?.getContext('2d') as ChartItem
		const chart = new Chart(ctx, {
			type: 'pie',
			data: chartData() as any,
			options: {
				radius: '80%',
				responsive: true,
				plugins: {
					legend: {
						position: 'bottom',
					},
					title: {
						display: true,
						text: 'Macro Breakdown for Today',
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

export { MacroPieChart }
