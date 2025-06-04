import { useEffect, useRef, useState, useCallback } from 'react'
import { Row, Col, Button } from 'react-bootstrap'
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler,
    Plugin,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import zoomPlugin from 'chartjs-plugin-zoom'
import annotationPlugin from 'chartjs-plugin-annotation'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { convertUtcTimeToLocalTimeMonthly, logger } from 'helpers'
import { LabelSelect2 } from 'components'
import { useSocket } from 'hooks'

const FilterArea = styled(Row)`
margin-top: 5px;
margin-left: 5px;

& button {
    margin-right: 5px;
}
`

const ChartArea = styled(Row)`
height: 25rem;
padding-right: 0;

& canvas {
    padding-right: 0;
}
`

const ButtonArea = styled.div`
margin-top: 5px;
& button {
    margin-right: 5px;
}
`

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler, zoomPlugin, annotationPlugin, {
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

const accuracyDatasets = {
    labels: [[0, '0']],
    datasets: [
        {
            label: 'Train',
            data: [0],
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        },
        {
            label: 'Valid',
            data: [0],
            borderColor: 'rgb(53, 162, 235)',
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        },
    ],
}

const lossDatasets = {
    labels: [[0, '0']],
    datasets: [
        {
            label: 'Train',
            data: [0],
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        },
        {
            label: 'Valid',
            data: [0],
            borderColor: 'rgb(53, 162, 235)',
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        },
    ],
}

const chartTitle = (tooltipItems: any) => {
    if(tooltipItems[0].label.includes('train')) return tooltipItems[0].label
    return 'Epoch ' + tooltipItems[0].label.split(',')[0]
}

const chartFooter = (tooltipItems: any) => {
    return tooltipItems[0].label.split(',')[1]
}

const getLossLabelForValue = (val: any) => {
    if (!lossDatasets || !lossDatasets.labels || !lossDatasets.labels[val]) return ''
    return lossDatasets.labels[val][0]
}

const getAccuracyLabelForValue = (val: any) => {
    if (!accuracyDatasets || !accuracyDatasets.labels || !accuracyDatasets.labels[val]) return ''
    return accuracyDatasets.labels[val][0]
}

const accuracyOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
        intersect: false,
        mode: 'index' as const,
    },
    plugins: {
        legend: {
            position: 'top' as const,
            labels: {
                boxHeight: 1,
                color: '#888888',
                font: { size: 12, weight: '400', },
            },
        },
        title: {
            display: true,
            text: 'Accuracy',
            color: '#000000',
            font: { size: 14, weight: '600', },
        },
        zoom: {
            zoom: {
                wheel: { enabled: true, },
                mode: "x" as const,
            },
            pan: { enabled: true, mode: 'x' as const },
            limits: { y: { min: 0, max: 100, }, },
        },
        tooltip: {
            callbacks: {
                title: chartTitle,
                label: (context: any) => {
                    let label = context.parsed.y
                    return label + '%'
                },
                footer: chartFooter,
            },
        },
        annotation: {
            annotations: Array<any>()
        },
    },
    scales: {
        x: {
            axis: 'x' as const,
            position: 'bottom' as const,
            title: { display: true, text: 'Epoch', },
            ticks: {
                callback: (val: any, index: number) => {
                    return getAccuracyLabelForValue(val)
                },
            },
        },
        y: {
            axis: 'y' as const,
            position: 'left' as const,
            title: { display: true, text: 'Accuracy (%)', },
        },
    },
}

const lossOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: { intersect: false, mode: 'index' as const, },
    plugins: {
        legend: {
            position: 'top' as const,
            labels: {
                boxHeight: 1,
                color: '#888888',
                font: { size: 14, weight: '400', },
            },
        },
        title: {
            display: true,
            text: 'Loss',
            color: '#000000',
            font: { size: 18, weight: '600', },
        },
        zoom: {
            zoom: {
                wheel: { enabled: true, },
                mode: "x" as const
            },
            pan: { enabled: true, mode: 'x' as const },
            limits: { y: { min: 0, max: 100, }, },
        },
        tooltip: {
            callbacks: {
                title: chartTitle,
                label: (context: any) => {
                    let label = context.parsed.y
                    return label
                },
                footer: chartFooter,
            },
        },
        annotation: {
            annotations: Array<any>()
        }
    },
    scales: {
        x: {
            axis: 'x' as const,
            position: 'bottom' as const,
            title: { display: true, text: 'Epoch', },
            ticks: {
                callback: (val: any, index: number) => {
                    return getLossLabelForValue(val)
                },
            },
        },
        y: {
            axis: 'y' as const,
            position: 'left' as const,
            title: { display: true, text: 'Loss', },
        },
    },
}

const plugins: Plugin[] = [
    {
        id: 'vertical-line',
        afterDraw: (chart: { tooltip?: any; scales?: any; ctx?: any }) => {
            // eslint-disable-next-line no-underscore-dangle
            if (chart.tooltip._active && chart.tooltip._active.length) {
                // find coordinates of tooltip
                const activePoint = chart.tooltip._active[0]
                const { ctx } = chart
                const { x } = activePoint.element
                const topY = chart.scales.y.top
                const bottomY = chart.scales.y.bottom

                // draw vertical line
                ctx.save()
                ctx.beginPath()
                ctx.moveTo(x, topY)
                ctx.lineTo(x, bottomY)
                ctx.lineWidth = 1
                ctx.strokeStyle = '#1C2128'
                ctx.stroke()
                ctx.restore()
            }
        },
    },
]

const initChartData = (charts: any, title: string[]) => {
    // chart 데이터 초기화
    charts.forEach((chart:any, index: number) => {
        if (chart && chart.current) {
            chart.current.options.plugins.title.text = title[index]
            chart.current.options.scales.y.title.text = title[index] + ' (%)'
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

            // annotation 초기화
            if (!chart.current || !chart.current.options || !chart.current.plugins) {

            } else {
                chart.current.options.plugins.annotation.annotations = []
            }

            chart.current.update()
        }
    })
}

const initAnnotation = (options: any) => {
    if (!options || options.length < 1) return
    options.forEach((option:any) => {
        option.plugins.annotation.annotations = []
    })
}

const addAnnotation = (option: any, label: string, value: any) => {
        option.plugins.annotation.annotations.push({
        type: 'line' as const,
        value: value,
        scaleID: 'x',
        borderColor: 'black',
        borderWidth: 1,
        borderDash: [3, 3],
        label: {
            display: true,
            backgroundColor: 'rgba(200,200,200,0.8)',
            drawTime: 'afterDatasetsDraw',
            content: label,
            position: 'end' as const
        },
    })
}

const refineChartData = (chartData: any, chartRefAccuracy: any, chartRefLoss: any) => {
    if (!chartData || chartData.length < 1) return

    for (let index = 0; index < chartData.length; index++) {
        const label = [
            `${chartData[index].cur_epoch}`,
            convertUtcTimeToLocalTimeMonthly(chartData[index].utc_ts * 1000),
        ]

        if (chartData[index].cur_epoch === 1) {
            addAnnotation(lossOptions, `t${chartData[index].local_id}`, label)
            addAnnotation(accuracyOptions, `t${chartData[index].local_id}`, label)
        }

        if (chartRefLoss && chartRefLoss.current) {
            chartRefLoss.current.data.labels.push(label)
            chartRefLoss.current.data.datasets[0].data.push(chartData[index].loss.toFixed(2))
            chartRefLoss.current.data.datasets[1].data.push(chartData[index].valid_loss.toFixed(2))
        }

        if (chartRefAccuracy && chartRefAccuracy.current) {
            chartRefAccuracy.current.data.labels.push(label)
            chartRefAccuracy.current.data.datasets[0].data.push(Number((chartData[index].acc * 100).toFixed(2)))
            chartRefAccuracy.current.data.datasets[1].data.push(Number((chartData[index].valid_acc * 100).toFixed(2)))
        }
    }

    chartRefLoss.current.update()
    chartRefAccuracy.current.update()
    
}

const VisionCLSTrainChart = ({ trial, trains }: any) => {
    const [t] = useTranslation('translation')
    const chartRefAccuracy: any = useRef<ChartJS | null>(null)
    const chartRefLoss: any = useRef<ChartJS | null>(null)

    const [socketConnected, setSocketConnected] = useState(false)
    const [trainOptions, setTrainOptions] = useState<any>([])
    const [selectedTrains, setSelectedTrains] = useState<any>([])

    const onResetZoom = () => {
        if (!chartRefAccuracy || !chartRefLoss) return
        chartRefAccuracy.current.resetZoom()
        chartRefLoss.current.resetZoom()
    }

    const onZoomIn = () => {
        chartRefAccuracy.current.zoom(1.01)
        chartRefLoss.current.zoom(1.01)
    }

    const onZoomOut = () => {
        chartRefAccuracy.current.zoom(0.99)
        chartRefLoss.current.zoom(0.99)
    }

    const handleSocketMessage = useCallback((e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            if (trains && trains.length > 0 && msg.length > 0) {
                const train_ids = msg.reduce((acc: number[], m:any) => {
                    if (acc.includes(m.local_id)) {
                        return acc
                    } 
                    
                    acc.push(m.local_id)
                    return acc 
                }, [])
                
                setSelectedTrains(
                    trains
                    .filter((t:any) => {
                        if (train_ids.includes(t.local_id)) return true
                        return false
                    })
                    .map((t:any) => {
                        return {
                            label: 'trial' + t.local_id,
                            value: t.uuid
                        }
                    })
                )
            }
            refineChartData(msg, chartRefAccuracy, chartRefLoss)
        } catch (e) {
            logger.error(e)
        }
    }, [trains, chartRefAccuracy, chartRefLoss])
    const ws = useSocket(`/trials/chart/visioncls/${trial && trial.trial_id}`, 'Chart', handleSocketMessage, {setSocketConnected, shouldCleanup: true, shouldConnect: !!trial && trial.trial_id !== 0})

    useEffect(() => {
        if (!trial || trial.trial_id === 0) return

        initChartData([chartRefAccuracy, chartRefLoss], [t(`metric.${trial.target_metric}`), 'Loss'])
        initAnnotation([accuracyOptions, lossOptions])
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial])


    useEffect(() => {
        if (!trains) return
        const tos = trains.map((train: any) => {
            return {
                value: train.uuid,
                label: `Trial${train.local_id}`
            }
        })
        setTrainOptions(tos)
    }, [trains])

    const onClickAllChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ filter: 'ALL' }))

        initChartData([chartRefAccuracy, chartRefLoss], [t(`metric.${trial.target_metric}`), 'Loss'])
        initAnnotation([accuracyOptions, lossOptions])
    }

    const onClickSelectedChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        
        const uuids = selectedTrains.map((train: any) => {
            return "'" + train.value + "'"
        })
        ws.current.send(JSON.stringify({ filter: 'SELECTED', train_uuids: uuids }))

        initChartData([chartRefAccuracy, chartRefLoss], [t(`metric.${trial.target_metric}`), 'Loss'])
        initAnnotation([accuracyOptions, lossOptions])
    }

    const onClickCurrentChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ filter: 'CURRENT' }))

        initChartData([chartRefAccuracy, chartRefLoss], [t(`metric.${trial.target_metric}`), 'Loss'])
        initAnnotation([accuracyOptions, lossOptions])
    }

    const onClickTopNChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ filter: 'TOPN', topn: 10 }))

        initChartData([chartRefAccuracy, chartRefLoss], [t(t(`metric.${trial.target_metric}`)), 'Loss'])
        initAnnotation([accuracyOptions, lossOptions])
    }

    return (
        <Col>
            <Row>
                <FilterArea>
                    <Col sm='4'>
                        <LabelSelect2
                            title={t('ui.train.select')}
                            name={'multiselect-train'}
                            options={trainOptions}
                            value={selectedTrains}
                            onChange={(options: any) => setSelectedTrains(options)}
                        />
                    </Col>
                    <Col>
                        <Button onClick={onClickSelectedChart}>
                        {t('button.viewTrain.selected')}
                        </Button>
                        <Button onClick={onClickAllChart}>
                        {t('button.viewTrain.all')}
                        </Button>
                        <Button onClick={onClickCurrentChart}>
                        {t('button.viewTrain.recent')}
                        </Button>
                        <Button onClick={onClickTopNChart}>
                        {t('button.viewTrain.toprated')}
                        </Button>
                    </Col>
                </FilterArea>
            </Row>
            <Row>
                <ButtonArea>
                    <Button variant='info' onClick={() => onResetZoom()}>
                        <i className='fas fa-sync' /> {t(`button.zoom.init`)}
                    </Button>
                    <Button variant='info' onClick={() => onZoomIn()}>
                        <i className='fas fa-plus' /> {t(`button.zoom.in`)}
                    </Button>
                    <Button variant='info' onClick={() => onZoomOut()}>
                        <i className='fas fa-minus' /> {t(`button.zoom.out`)}
                    </Button>
                </ButtonArea>
                <ChartArea>
                    <Line ref={chartRefAccuracy} data={accuracyDatasets} options={accuracyOptions} redraw={true} plugins={plugins} />
                </ChartArea>
            </Row>
            <Row>
                <ChartArea>
                    <Line ref={chartRefLoss} data={lossDatasets} options={lossOptions} redraw={true} plugins={plugins} />
                </ChartArea>
            </Row>
        </Col>
    )
}

export default VisionCLSTrainChart