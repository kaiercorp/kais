import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler,
    Plugin
} from 'chart.js'
import { Chart } from 'react-chartjs-2'
import zoomPlugin from 'chartjs-plugin-zoom'
import annotationPlugin from 'chartjs-plugin-annotation'
import { Button, Col, Row } from 'react-bootstrap'
import { useEffect, useRef, useState, useCallback } from 'react'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'

import { LabelSelectTyped } from 'components'
import { logger, ApiDownloadFile } from 'helpers'
import { useSocket } from 'hooks'

const FilterArea = styled(Row)`
margin-top: 5px;
margin-left: 5px;

& button {
    margin-right: 5px;
}
`

interface IChartArea {
    height?: number
}

const ChartArea = styled(Row) <IChartArea>`
height: ${(props) => (props.height ? props.height : 600)}px;
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

const DownloadBtnArea = styled.div`
display: inline-block;
float: right;
`

ChartJS.register(
    CategoryScale, 
    LinearScale, 
    BarElement, 
    LineElement,
    PointElement,
    Title, 
    Tooltip, 
    Legend, 
    Filler, 
    zoomPlugin, 
    annotationPlugin, 
    {
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
            xAxisID: 'x'
            },
        {
            data: [0],
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            pointRadius: 5,
            pointBorderWidth: 2,
            lineTension: 0,
            type: 'line' as const,
            fill: false,
            xAxisID: 'x1',
            segment: {
                borderColor: 'rgb(0,0,0,0.2)',
                borderDash: [1,1],
            },
            spanGaps: true
        }
    ]
}

const chartTitle = (tooltipItems: any) => {
    if (tooltipItems[0].label.includes('train')) return tooltipItems[0].label
    return tooltipItems[0].label.split(',')[0]
}

const chartFooter = (tooltipItems: any) => {
    return tooltipItems[0].label.split(',')[1]
}

const getChartLabelForValue = (val: any) => {
    if (!chartData || !chartData.labels || !chartData.labels[val]) return ''
    return chartData.labels[val][0]
}

const chartOption = {
    indexAxis: 'y' as const,
    responsive: true,
    maintainAspectRatio: false,
    elements: {
        bar: {
            borderWidth: 2,
        },
    },
    interaction: {
        intersect: false,
        mode: 'point' as const,
    },
    plugins: {
        legend: {
            display: false,
        },
        title: {
            display: false,
        },
        zoom: {
            zoom: {
                wheel: { enabled: true, },
                mode: "y" as const
            },
            pan: { 
                enabled: true, 
                mode: 'y' as const
            },
        },
        tooltip: {
            callbacks: {
                title: chartTitle,
                label: (context: any) => {
                    let label = Number(context.parsed.x).toFixed(3)
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
            position: 'top' as const,
            title: {
                display: true,
                text: 'Feature Importance',
                color: 'blue'
            },
        },
        x1: {
            axis: 'x' as const,
            type: 'linear' as const,
            display: true,
            position: 'bottom' as const,
            title: {
                display: true,
                text: 'Target Metric',
                color: 'red'
            },
            grid: {
                drawOnChartArea: false,
            },
        },
        y: {
            axis: 'y' as const,
            position: 'top' as const,
            title: {
                display: true,
            },
            ticks: {
                callback: (val: any, index: number) => {
                    return getChartLabelForValue(val)
                },
            },
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

        // annotation 초기화
        if (!chart.current || !chart.current.options || !chart.current.plugins) {

        } else {
            chart.current.options.plugins.annotation.annotations = []
        }

        chart.current.update()
    }
}

const initAnnotation = (option: any) => {
    if (!option) return
    option.plugins.annotation.annotations = []
}

const addAnnotation = (option: any, label: string, value: any) => {
    option.plugins.annotation.annotations.push({
        type: 'line' as const,
        value: value,
        scaleID: 'y',
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

const refineChartData = (data: any, chart: any) => {
    if (!data || data.length < 1) return

    const fi_labels = Object.keys(JSON.parse(data[0].feature_importance.String))
    const mean_fi_labels = Math.floor(fi_labels.length / 2)

    for (let index = 0; index < data.length; index++) {

        const target_metric = data[index].target_metric.String
        const all_result = JSON.parse(data[index].all_result.String)
        const feature_importance = JSON.parse(data[index].feature_importance.String)

        fi_labels.forEach((label: any, idx: number) => {
            const xlabel = [
                `${label}`,
                data[index].updated_at,
            ]
            
            if (idx === 0) addAnnotation(chartOption, `t${data[index].train_id}`, xlabel)

            chart.current.data.labels.push(xlabel)
            chart.current.data.datasets[0].data.push(feature_importance[label])
            if (mean_fi_labels === idx) {
                chart.current.data.datasets[1].data.push(all_result[target_metric])
            } else {
                chart.current.data.datasets[1].data.push(NaN)
            }
        })
    }

    chart.current.update()
}

const TableClsFIContainer = ({ trial, trains }: any) => {
    const [t] = useTranslation('translation')
    const chart: any = useRef<ChartJS | null>(null)

    const [socketConnected, setSocketConnected] = useState(false)
    const [trainOptions, setTrainOptions] = useState<any>([])
    const [selectedTrains, setSelectedTrains] = useState<any>([])

    const queryClient = useQueryClient()

    const onResetZoom = (chart: any) => {
        if (!chart) return
        chart.resetZoom()
    }

    const onZoomIn = (chart: any) => {
        chart.zoom(1.01)
    }

    const onZoomOut = (chart: any) => {
        chart.zoom(0.99)
    }

    const handleSocketMessage = useCallback((e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            if (trains && trains.length > 0 && msg.length > 0) {
                const train_ids = msg.map((m:any) => {
                    return m.train_id
                })
                
                setSelectedTrains(
                    trains
                    .filter((t:any) => {
                        if (train_ids.includes(t.local_id)) return true
                        return false
                    })
                    .map((t:any) => {
                        return {
                            label: 'trial' + t.local_id,
                            value: t.id
                        }
                    })
                )
            }
            refineChartData(msg, chart)
        } catch (e) {
            logger.error(e)
        }
    }, [trains, chart])
    const ws = useSocket(`/trials/chart/tablecls/${trial.trial_id}`, 'Chart', handleSocketMessage, {setSocketConnected, shouldCleanup: true, shouldConnect: !!trial && trial.trial_id && trains && trains.length})

    useEffect(() => {
        if (!trial || trial.trial_id === 0 || !trains || trains.length < 1) return

        initChartData(chart)
        initAnnotation(chartOption)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial, trains])

    useEffect(() => {
        if (!trains) return
        const tos = trains.map((train: any) => {
            return {
                value: train.id,
                label: `Trial${train.local_id}`
            }
        })
        
        setTrainOptions(tos)
    }, [trains])

    // const onClickAllChart = () => {
    //     if (!socketConnected || !ws || !ws.current) return
    //     ws.current.send(JSON.stringify({ filter: 'ALL' }))

    //     initChartData(chart)
    //     initAnnotation(chartOption)
    // }

    const onClickSelectedChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        if (selectedTrains.length < 1) return
        
        const train_ids = selectedTrains.map((train: any) => {
            return train.value
        })
        ws.current.send(JSON.stringify({ filter: 'SELECTED', train_ids: train_ids }))

        initChartData(chart)
        initAnnotation(chartOption)
    }

    const onClickCurrentChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ filter: 'CURRENT', topn: 20 }))

        initChartData(chart)
        initAnnotation(chartOption)
    }

    const onClickTopNChart = () => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ filter: 'TOPN', topn: 10 }))

        initChartData(chart)
        initAnnotation(chartOption)
    }
        
    const onClickDownloadAllFeatureImportances = () => {
        ApiDownloadFile(queryClient, `/train/feature-importance/download`, `${trial.trial_name}_${trial.trial_id}.zip`, {trial_id: trial.trial_id}) 
    }

    return (
        <Col>
            <Row>
                <FilterArea>
                    <Col sm='4'>
                        {/* <LabelSelect
                            title={t('ui.train.select')}
                            name={'multiselect-train'}
                            value={selectedTrains[0]}
                            onChange={(e: any) => handleSelectTrain(e.target.value)}
                        >
                            <option key='' value=''>{t('ui.label.let.select')}</option>
                            {trains && trains.map((train: any) => {
                                return <option key={`Trial${train.local_id}`} value={train.uuid}>train{train.local_id}</option>
                            })}
                        </LabelSelect> */}
                        <LabelSelectTyped 
                            title={t('ui.train.select')}
                            name={'multiselect-train'}
                            options={trainOptions}
                            value={selectedTrains}
                            onChange={(options: any) => setSelectedTrains(options)}
                        />
                    </Col>
                    <Col>
                        <Button onClick={onClickSelectedChart} disabled={selectedTrains.length < 1}>
                            {t('button.viewTrain.selected')}
                        </Button>
                        {/* <Button onClick={onClickAllChart}>
                            {t('button.viewTrain.all')}
                        </Button> */}
                        <Button onClick={onClickCurrentChart}>
                            {t('button.viewTrain.recent')}
                        </Button>
                        <Button onClick={onClickTopNChart}>
                            {t('button.viewTrain.toprated')}
                        </Button>
                        <DownloadBtnArea>
                            <Button onClick={onClickDownloadAllFeatureImportances}>
                                {t('button.download.allOfTrains')}
                            </Button>
                        </DownloadBtnArea>
                    </Col>
                </FilterArea>
            </Row>
            <Row>
                <ButtonArea>
                    <Button variant='info' onClick={() => onResetZoom(chart.current)}>
                        <i className='fas fa-sync' /> {t('button.zoom.init')}
                    </Button>
                    <Button variant='info' onClick={() => onZoomIn(chart.current)}>
                        <i className='fas fa-plus' /> {t('button.zoom.in')}
                    </Button>
                    <Button variant='info' onClick={() => onZoomOut(chart.current)}>
                        <i className='fas fa-minus' /> {t('button.zoom.out')}
                    </Button>
                </ButtonArea>
            </Row>
            <Row>
                <ChartArea height={500}>
                    <Chart type='bar' ref={chart} options={chartOption} data={chartData} redraw={true} plugins={plugins} />
                </ChartArea>
            </Row>
        </Col>
    )
}

export default TableClsFIContainer