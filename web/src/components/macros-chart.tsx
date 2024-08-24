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
import { fetchFoodInsights } from '~/lib/api'

Chart.register(PieController, ArcElement, Tooltip, Legend)

const MacroPieChart = () => {
	const [canvasRef, setCanvasRef] = createSignal<HTMLCanvasElement | null>()
	const [chartData, setChartData] = createSignal()

	const query = createQuery(() => ({
		queryKey: ['posts', 'macro-pie-chart'],
		queryFn: () => fetchFoodInsights(),
	}))

	createEffect(() => {
		if (!query.data) return

		const macros = query.data.macros

		setChartData({
			labels: ['Proteins', 'Fats', 'Carbohydrates'],
			datasets: [
				{
					data: [macros.proteins, macros.fats, macros.carbohydrates],
					backgroundColor: [
						'#fb2954', // Red for Proteins
						'#4277f8', // Blue for Fats
						'#f8c542', // Yellow for Carbohydrates
					],
					hoverBackgroundColor: ['#fb2954', '#4277f8', '#f8c542'],
					borderWidth: 0,
				},
			],
		})

		const ctx = canvasRef()?.getContext('2d') as ChartItem
		const chart = new Chart(ctx, {
			type: 'pie',
			data: chartData() as any,
			options: {
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
