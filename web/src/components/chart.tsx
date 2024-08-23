import { createEffect, createSignal, onCleanup } from 'solid-js'
import {
	Chart,
	BarElement,
	CategoryScale,
	LinearScale,
	Title,
	Tooltip,
	Legend,
	ChartItem,
	Filler,
	Colors,
	BarController,
} from 'chart.js'

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

const CaloricMacroBreakdownChart = () => {
	const [canvasRef, setCanvasRef] = createSignal<HTMLCanvasElement | null>()

	const exampleData = {
		labels: [
			'2024-08-20',
			'2024-08-21',
			'2024-08-22',
			'2024-08-23',
			'2024-08-24',
		],
		datasets: [
			{
				label: 'Protein (g)',
				backgroundColor: 'rgba(8,39,230,0.85)',
				data: [80, 90, 70, 85, 95],
				borderRadius: 10,
			},
			{
				label: 'Fat (g)',
				backgroundColor: 'rgba(244,123,4,0.98)',
				data: [70, 60, 65, 75, 80],
				borderRadius: 10,
			},
			{
				label: 'Carbohydrates (g)',
				backgroundColor: 'rgb(66,119,248)',
				data: [250, 230, 240, 260, 270],
				borderRadius: 10,
			},
			{
				label: 'Calories',
				backgroundColor: 'rgb(243,51,51)',
				data: [2000, 1900, 1800, 2100, 2200],
				borderRadius: 10,
			},
		],
	}

	createEffect(() => {
		const ctx = canvasRef()?.getContext('2d') as ChartItem
		const chart = new Chart(ctx, {
			type: 'bar',
			data: exampleData,
			options: {
				responsive: true,
				plugins: {
					legend: {
						position: 'top',
					},
					title: {
						display: true,
						text: 'Daily Caloric and Macronutrient Breakdown',
					},
				},
				scales: {
					x: {
						stacked: true,
					},
					y: {
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

export default CaloricMacroBreakdownChart
