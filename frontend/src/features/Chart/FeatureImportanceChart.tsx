import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js'
import { Bar } from 'react-chartjs-2'
import zoomPlugin from 'chartjs-plugin-zoom'
import styled from 'styled-components'
import { Row } from 'react-bootstrap'
import { useEffect, useRef } from 'react'


interface IChartArea {
    height?: number
}

const ChartArea = styled(Row) <IChartArea>`
height: ${(props) => (props.height ? props.height : 600)}px;
& canvas {
    padding-right: 0;
}
`

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, zoomPlugin, {
    id: 'train_line_charts',
    beforeDraw: (chart) => {
        const { ctx } = chart
        ctx.save()
        ctx.globalCompositeOperation = 'destination-over'
        ctx.fillStyle = '#ffffff'
        ctx.fillRect(0, 0, chart.width, chart.height)
        ctx.restore()
    },
})

const chartData = {
    labels: ['0'],
    datasets: [
        {
            data: [0],
            borderColor: 'rgb(53, 162, 235)',
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
        }
    ]
}

const options = {
    indexAxis: 'y' as const,
    maintainAspectRatio: false,
    elements: {
        bar: {
            borderWidth: 2,
        },
    },
    responsive: true,
    plugins: {
        legend: {
            display: false,
            position: 'right' as const,
        },
        title: {
            display: false,
        },
        zoom: {
            zoom: {
                wheel: {enabled: true},
                mode: "y" as const
            },
            pan: {enabled: true, mode: 'x' as const}
        }
    },
    scales: {
        x: {
            axis: 'x' as const,
            position: 'bottom' as const,
            title: {
                display: true,
                text: 'Feature Importance',
            }
        }
    }
}

const refineChartData = (data: any, chart: any) => {
    if (!data || data.length < 1) return

    const labels = Object.keys(data)
    labels.forEach((label: any) => {
        chart.current.data.labels.push(label)
        chart.current.data.datasets[0].data.push(data[label])
    })

    chart.current.update()
}

const initChartData = (chart: any) => {
    if (chart && chart.current) {
        // x축 라벨 초기화
        while (chart.current.data.labels.length > 0) {
            chart.current.data.labels.pop()
        }
        
        // y축 값 초기화
        chart.current.data.datasets.forEach((dataset: any) => {
            while (dataset.data.length > 0) {
                dataset.data.pop()
            }
        })

        chart.current.update()
    }
}

const FeatureImportanceChart = ({ feature_importance }: any) => {
    const chart: any = useRef<ChartJS | null>(null)

    useEffect(() => {
        if (!feature_importance) return

        initChartData(chart)
        const fe = JSON.parse(feature_importance)
        refineChartData(fe, chart)
    }, [feature_importance])

    return (
        <ChartArea height={500}>
            <Bar ref={chart} options={options} data={chartData} />
        </ChartArea>
    )
}

export default FeatureImportanceChart