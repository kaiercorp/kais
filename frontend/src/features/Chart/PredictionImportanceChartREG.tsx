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
import { Card, Col, Row } from 'react-bootstrap'

const ChartArea = styled(Row)`
height: 600px;
`

const LeftCol = styled(Col)`
text-align: end;
color: rgb(255, 99, 132);
`

const RightCol = styled(Col)`
color: rgb(53, 162, 235);
`

// const ButtonArea = styled.div`
//   margin-bottom: 5px;
// `

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


const PredictionImportanceChartREG = ({ data }: any) => {
    if (!data) return <></>

    var label: any = []
    var _data: any = []
    data.forEach((d: any) => {
        label.push(d[0])
        _data.push(d[1])
    })

    const borderColors = _data.map((v: any) => v >= 0 ? 'rgb(53, 162, 235)' : 'rgb(255, 99, 132)')
    const backgroundColors = _data.map((v: any) => v >= 0 ? 'rgba(53, 162, 235, 0.5)' : 'rgba(255, 99, 132, 0.2)')

    const chartData = {
        labels: label,
        datasets: [
            {
                data: _data,
                borderColor: borderColors,
                backgroundColor: backgroundColors,
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
                    wheel: { enabled: true, },
                    mode: "x" as const
                },
                pan: { enabled: true, mode: 'x' as const },
            }
        },
        scales: {
            x: {
                axis: 'x' as const,
                position: 'bottom' as const,
                title: {
                    display: true,
                }
            }
        }
    }

    return (
        <Card>
            <Card.Header>
            Local explanation
            </Card.Header>
            <Card.Body>
            <ChartArea>
                <Row>
                    <LeftCol>Decreasing</LeftCol>
                    <RightCol>Increasing</RightCol>
                </Row>
                <Row>
                    <Bar options={options} data={chartData} />
                </Row>
            </ChartArea>
            </Card.Body>
        </Card>
    )
}

export default PredictionImportanceChartREG