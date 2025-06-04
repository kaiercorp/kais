import { useEffect, useRef, useState, useCallback, useMemo, ChangeEvent } from "react"
import { useTranslation } from "react-i18next"
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
import { Line } from 'react-chartjs-2'
import zoomPlugin from 'chartjs-plugin-zoom'
import annotationPlugin from 'chartjs-plugin-annotation'
import { logger } from "helpers"
import { Button, Col, Row } from "react-bootstrap"
import styled from "styled-components"
import { LabelSelect2, LabelInput } from "components"
import { useSocket } from 'hooks'
import { TrialType, TrainType, ChartType } from 'common'
import gaussian from 'gaussian'
import Decimal from 'decimal.js'

type TSADChartDatasetType = {
    label: string;
    data: number[];
    backgroundColor: string;
    borderColor: string;
    pointRadius: number;
    tension: number;
    yAxisID: string;
}

const FilterArea = styled(Row)`
margin-top: 5px;
margin-left: 5px;

& button {
    margin-right: 5px;
}
`

const ChartArea = styled(Row)`
display: flex;
flex-wrap: wrap;
height: 500px;
& canvas {
    padding-right: 0;
}
margin-bottom: 20px;
`

const ButtonArea = styled.div`
margin-top: 5px;
& button {
    margin-right: 5px;
}
`

const ChartConfigurationContainer = styled.div`
display: flex;
flex-direction: row;
justify-content: space-between;
padding-right: 24px;

th, td {
    padding-top: 0.65rem;
    min-width: 100px;
}
`

const TrainChartArea = styled.div`
    flex-basis: 60%;
    padding-right: 0;
`

const GaussianChartArea = styled.div`
    flex-basis: 40%;
    padding-right: 0;
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
    datasets: []
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

const chartTitle = (tooltipItems: any) => {
    if (tooltipItems[0].label.includes('train')) return tooltipItems[0].label
    return tooltipItems[0].label.split(',')[0]
}

const initializeChartOption = (
    options: {
        title?: string,
        scales?: {
            [key: string]: {
                [key: string]: any
            }
        },
        displayLegend?: boolean,
    } = {}
) => ({
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
        intersect: false,
        mode: 'index' as const,
    },
    plugins: {
        legend: {
            display: typeof (options.displayLegend) === 'boolean' ? options.displayLegend : true,
            position: 'top' as const,
            labels: {
                boxHeight: 1,
                color: '#888888',
                font: { size: 12, weight: '400' },
            },
        },
        title: {
            display: true,
            text: options.title ? options.title : 'Train Chart',
            color: '#000000',
            font: { size: 14, weight: '600', },
        },
        zoom: {
            zoom: {
                wheel: { enabled: true, mode: 'x' as const },
                mode: 'x' as const
            },
            pan: {
                enabled: true,
                mode: 'x' as const,
            },
        },
        tooltip: {
            callbacks: {
                title: chartTitle,
                label: (context: any) => {
                    let label = context.parsed.y
                    return label
                },
            },
        },
        annotation: {
            annotations: Array<any>()
        }
    },
    scales: {
        ...options.scales
    },
})

const trainChartOption = initializeChartOption({
    scales: {
        x: {
            axis: 'x' as const,
            position: 'bottom' as const,
            title: {
                display: true,
                text: 'Time'
            },
            ticks: {
                callback: (val: any, index: number) => {
                    return val
                }
            }
        },
        y: {
            type: 'linear' as const,
            display: true,
            position: 'left' as const,
            title: {
                display: true,
                text: 'Input Signal',
                color: 'blue'
            },
            grid: {
                drawOnChartArea: false,
            },
        },
        y1: {
            type: 'linear' as const,
            display: true,
            position: 'right' as const,
            title: {
                display: true,
                text: 'Anomaly Score',
                color: 'red'
            },
        }
    }
})

const initChartData = (chart: any) => {
    if (chart && chart.current) {
        while (chart.current.data.labels.length > 0) {
            chart.current.data.labels.pop()
        }

        chart.current.data.datasets.forEach((dataset: any) => {
            while (dataset.data.length > 0) {
                dataset.data.pop()
            }
        })

        chart.current.update()
    }
}

const initAnnotation = (option: any) => {
    if (!option) return
    option.plugins.annotation.annotations = []
}

const addAnnotation = (option: any, label: string, value: number, axis: string) => {
    const annotationValue =
        axis === 'x' ?
            {
                xMin: value,
                xMax: value
            }
            :
            {
                yMin: value,
                yMax: value,
            }

    option.plugins.annotation.annotations.pop()

    option.plugins.annotation.annotations.push({
        type: 'line' as const,
        borderColor: 'black',
        borderWidth: 1,
        borderDash: [6, 6],
        label: {
            display: true,
            backgroundColor: 'rgba(200,200,200,0.8)',
            drawTime: 'afterDatasetsDraw',
            content: label,
            position: 'start' as const
        },
        yScaleID: axis,
        ...annotationValue
    })
}

const random_rgb = () => {
    var o = Math.round, r = Math.random, s = 255
    return [o(r() * s), o(r() * s), o(r() * s)]
}

const refineChartData = (chartRef: any, dataset: { [key: string]: number | number[] | string, data: number[] }[], options: { addAnnotation?: (option: any, label: string, value: number, axis: string) => void, chartOption?: any, threshold?: number, annotationAxis?: string } = {}) => {
    if (!chartRef.current) { return }


    // let newDatasets = Array<any>()

    // let first_data = JSON.parse(chartData[0].learning_chart_material.String)
    // chartRef.current.data.labels = Array.from({ length: first_data.pred.length }, (v, i) => i)

    // newDatasets.push({
    //     label: `True`,
    //     data: first_data.true,
    //     borderColor: `rgb(53, 162, 235)`,
    //     backgroundColor: 'rgba(53, 162, 235, 0.5)',
    //     pointRadius: 0,
    //     tension: 0.1,
    //     yAxisID: 'y',
    // })
    // addAnnotation(chartOption, `Threshold ${first_data.threshold}`, first_data.threshold)

    // for (let index = 0; index < chartData.length; index++) {
    //     let data = JSON.parse(chartData[index].learning_chart_material.String)
    //     let rgb = random_rgb()
    //     newDatasets.push({
    //         label: `Trial${chartData[index].train_id}`,
    //         data: data.pred,
    //         borderColor: `rgb(${rgb[0]}, ${rgb[1]}, ${rgb[2]})`,
    //         backgroundColor: `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, 0.5)`,
    //         pointRadius: 0,
    //         tension: 0.1,
    //         yAxisID: 'y1',
    //     })

    if (options.addAnnotation && options.chartOption && options.threshold && options.annotationAxis) {
        options.addAnnotation(options.chartOption, `Threshold ${options.threshold}`, options.threshold, options.annotationAxis)

    }

    chartRef.current.data.labels = Array.from({ length: dataset[0].data.length }, (v, i) => i)
    chartRef.current.data.datasets = dataset
    chartRef.current.update()
}

type TSADTrainContainerProps = {
    trial: TrialType
    trains: TrainType[] | undefined
}

const TSADTrainContainer = ({ trial, trains }: TSADTrainContainerProps) => {
    const [t] = useTranslation('translation')
    const trainChart: any = useRef<ChartJS | null>(null)
    const gaussianChart: any = useRef<ChartJS | null>(null)

    const [socketConnected, setSocketConnected] = useState(false)
    const [trainOptions, setTrainOptions] = useState<any>([])
    const [selectedTrains, setSelectedTrains] = useState<{ value: string, label: string }[]>([])
    const [input, setInput] = useState<{ threshold: string, error: string }>({ threshold: '', error: '' })

    const gaussianInfo = useMemo(() => {
        if (!trainChart.current || !trainChart.current.data) { return false }

        const datasets = trainChart.current.data.datasets
        if (datasets.length !== 2) { return false }

        const trainDataset =
            datasets
                .filter((dataset: TSADChartDatasetType) => {
                    return dataset.label.startsWith('Trial')
                })[0]
                .data

        const predictErrorRange =
            trainDataset
                .reduce((acc: number[], value: number) => {
                    if (!acc.includes(value)) {
                        acc.push(value)
                    }

                    return acc
                }, [])
                .sort((a: number, b: number) => a - b)

        const mean =
            trainDataset
                .reduce((acc: number, value: number) => new Decimal(acc).plus(new Decimal(value)), 0)
                .div(trainDataset.length)
                .toNumber()

        const variance =
            trainDataset
                .reduce((acc: number, value: number) => Decimal.pow(value - mean, 2), 0)
                .div(trainDataset.length)
                .toNumber()

        if (variance === 0) { return false }

        const distribution = gaussian(mean, variance)

        const length = predictErrorRange.length;
        const pdfDataset =
            Array(length).fill(null)
                .map((value: number, index: number) => {
                    return distribution.pdf(mean + (index - Math.floor(length / 2)))
                })
        const offsetX = mean - Math.floor(length / 2)

        const gaussianChartOption = initializeChartOption({
            title: 'Gaussian Distribution',
            scales: {
                x: {
                    axis: 'x' as const,
                    position: 'bottom' as const,
                    title: { display: true, text: 'P(X)', },
                    ticks: {
                        callback: (val: any, index: number) => {
                            return pdfDataset[index]
                        }
                    }
                },
                y: {
                    type: 'linear' as const,
                    display: true,
                    position: 'left' as const,
                    title: {
                        display: true,
                        text: 'X',
                        color: 'blue'
                    },
                    grid: {
                        drawOnChartArea: false,
                    },
                },
            },
            displayLegend: false,
        })

        return {
            mean,
            variance,
            distribution,
            trainDataset,
            pdfDataset,
            offsetX,
            predictErrorRange,
            gaussianChartOption
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trainChart.current?.data?.datasets])

    useEffect(() => {
        if (!trial || trial.trial_id === 0 || !trains || trains.length < 1) return

        initChartData(trainChart)
        initAnnotation(trainChartOption)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial, trains])

    useEffect(() => {
        if (!Array.isArray(trains)) return

        setTrainOptions(trains.map((train: any) => {
            return {
                value: train.uuid,
                label: `Trial${train.local_id}`
            }
        }))
    }, [trains])

    useEffect(() => {
        if (!gaussianInfo) { return }

        const dataset = [{
            label: `gaussian`,
            data: gaussianInfo.predictErrorRange.slice(0, Math.floor(gaussianInfo.predictErrorRange.length / 2)),
            borderColor: `rgb(53, 162, 235)`,
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        }, {
            label: `gaussian`,
            data: gaussianInfo.predictErrorRange.slice(Math.floor(gaussianInfo.predictErrorRange.length / 2)).reverse(),
            borderColor: `rgb(53, 162, 235)`,
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
            pointRadius: 0,
            tension: 0.1,
        }]

        refineChartData(gaussianChart, dataset)
    }, [gaussianInfo])

    useEffect(() => {
        if (/(^$)|(^[0-9]+\.$)/.test(input.threshold)) {
            if (trainChartOption && trainChart) {
                initAnnotation(trainChartOption)
                trainChart.current.update()
            }
            resetGaussianChartAnnotation()
            return
        }

        if (trainChartOption && trainChart) {
            addAnnotation(trainChartOption, `Threshold ${input.threshold}`, Number(input.threshold), 'y1')
            trainChart.current.update()
        }

        updateGaussianChartAnnotation()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [input.threshold])

    useEffect(() => {
        if (/(^$)|(^[0-9]+\.$)/.test(input.threshold)) {
            resetGaussianChartAnnotation()
        } else {
            updateGaussianChartAnnotation()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [gaussianInfo, gaussianChart.current])

    const resetGaussianChartAnnotation = () => {
        if (gaussianInfo && gaussianChart.current) {
            initAnnotation(gaussianInfo.gaussianChartOption)
            gaussianChart.current.update()
        }
    }

    const updateGaussianChartAnnotation = () => {
        if (gaussianInfo && gaussianChart.current) {
            addAnnotation(gaussianInfo.gaussianChartOption, `Threshold ${input.threshold}`, Number(input.threshold), 'y')
            gaussianChart.current.update()
        }
    }

    const onResetZoom = (chart: any) => {
        if (!chart) return
        chart.resetZoom()
    }

    const onZoomIn = (chart: any) => {
        if (!chart) return
        chart.zoom(1.01)
    }

    const onZoomOut = (chart: any) => {
        if (!chart) return
        chart.zoom(0.99)
    }

    const handleSocketMessage = useCallback((e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)

            if (Array.isArray(trains) && Array.isArray(msg)) {
                setSelectedTrains(
                    trains
                        .filter((train: TrainType) => {
                            return msg.some((data: ChartType) => {
                                return train.local_id === data.train_id
                            })
                        })
                        .map((train: TrainType) => ({
                            value: train.uuid,
                            label: `Trial${train.local_id}`
                        }))
                )
            }

            if (Array.isArray(msg) && msg.length >= 1) {
                const dataset = [{
                    label: `Input`,
                    data: JSON.parse(msg[0].learning_chart_material.String).true,
                    borderColor: `rgb(53, 162, 235)`,
                    backgroundColor: 'rgba(53, 162, 235, 0.5)',
                    pointRadius: 0,
                    tension: 0.1,
                    yAxisID: 'y',
                }].concat(
                    msg.map((data: any) => {
                        const rgb = random_rgb()
                        return {
                            label: `Trial${data.train_id}`,
                            data: JSON.parse(data.learning_chart_material.String).score,
                            borderColor: `rgb(${rgb[0]}, ${rgb[1]}, ${rgb[2]})`,
                            backgroundColor: `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, 0.5)`,
                            pointRadius: 0,
                            tension: 0.1,
                            yAxisID: 'y1',
                        }
                    })
                )

                refineChartData(trainChart, dataset, { addAnnotation, chartOption: trainChartOption, threshold: Number(input.threshold), annotationAxis: 'y1' })
            }
        } catch (e) {
            logger.error(e)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trains, trainChart])
    const ws = useSocket(`/trials/chart/tsad/${trial?.trial_id}`, 'Chart', handleSocketMessage, { setSocketConnected, shouldCleanup: true, shouldConnect: trial && trial.trial_id && trial.trial_id !== 0 && trains && trains.length })

    // useEffect(() => {
    //     if (!trial || trial.trial_id === 0 || !trains || trains.length < 1) return

    //     initChartData(chart)
    //     initAnnotation(chartOption)
    //     // eslint-disable-next-line react-hooks/exhaustive-deps
    // }, [trial, trains])

    // useEffect(() => {
    //     if (!trains) return
    //     const tos = trains.map((train: any) => {
    //         return {
    //             value: train.uuid,
    //             label: `Trial${train.local_id}`
    //         }
    //     })
    //     setTrainOptions(tos)
    // }, [trains])


    const onClickSelectedChart = () => {
        if (isWebSocketReady() && selectedTrains.length) {
            const uuids = selectedTrains.map((train: any) => {
                return train.value
            })
            sendChartDataRequest(JSON.stringify({ filter: 'SELECTED', uuids: uuids }))
            initChartData(trainChart)
            initAnnotation(trainChartOption)
        }
    }

    const onClickCurrentChart = () => {
        if (isWebSocketReady()) {
            sendChartDataRequest(JSON.stringify({ filter: 'CURRENT', topn: 10 }))
            initChartData(trainChart)
            initAnnotation(trainChartOption)
        }
    }

    const onClickTopNChart = () => {
        if (isWebSocketReady()) {
            sendChartDataRequest(JSON.stringify({ filter: 'TOPN', topn: 10 }))
            initChartData(trainChart)
            initAnnotation(trainChartOption)
        }
    }

    const isWebSocketReady = () => {
        return socketConnected && ws && ws.current
    }

    const sendChartDataRequest = (msg: string) => {
        ws.current!.send(msg)
    }

    const handleChangeThreshold = (e: ChangeEvent<HTMLInputElement>) => {
        const threshold = e.target.value

        if (/^$|^\d+\.$|^\d+(\.\d{1,2})?$/g.test(threshold)) {
            setInput({ threshold, error: calculateError(threshold) })
        }
    }

    const handleChangeError = (e: ChangeEvent<HTMLInputElement>) => {
        let error = (e.target.value)

        if (/^$|^\d+\.$|^\d+(\.\d{1,2})?$/g.test(error)) {
            if (Number.isNaN(Number(error)) || Number(error) > 100) {
                return
            }

            if (gaussianInfo) {
                const threshold = calculateThreshold(gaussianInfo.predictErrorRange, gaussianInfo.trainDataset, error)
                setInput({ threshold, error })
            }
        }
    }

    const calculateError = (threshold: string): string => {
        if (threshold === '' || Number.isNaN(Number(threshold))) { return '' }
        if (!trainChart.current || !trainChart.current.data) { return '' }

        const datasets = trainChart.current.data.datasets
        if (datasets.length !== 2) { return '' }

        const trainDataset =
            datasets
                .filter((dataset: TSADChartDatasetType) => {
                    return dataset.label.startsWith('Train')
                })[0]
                .data

        const countOfError =
            trainDataset
                .filter((value: number) => {
                    return value > Number(threshold)
                })
                .length

        return ((countOfError / trainDataset.length) * 100).toFixed(2)
    }

    const calculateThreshold = (range: number[], anomalyScore: number[], error: string): string => {
        if (error === '' || Number.isNaN(Number(error))) { return '' }

        let left = 0
        let right = range.length - 1
        while (left <= right) {
            const mid = Math.floor((left + right) / 2)

            const percent = (anomalyScore.reduce((acc, value) => {
                if (value >= range[mid]) {
                    acc++
                }

                return acc
            }, 0) / anomalyScore.length) * 100

            if (percent >= Number(error)) {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }

        return range[right].toFixed(2)
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
                        <Button onClick={onClickSelectedChart} disabled={selectedTrains.length < 1}>
                            {t('button.viewTrain.selected')}
                        </Button>
                        <Button onClick={onClickCurrentChart}>
                            {t('button.viewTrain.recent10')}
                        </Button>
                        <Button onClick={onClickTopNChart}>
                            {t('button.viewTrain.top10')}
                        </Button>
                    </Col>
                </FilterArea>
            </Row>
            <Row>
                <ChartConfigurationContainer>
                    <ButtonArea>
                        <Button variant='info' onClick={() => onResetZoom(trainChart.current)}>
                            <i className='fas fa-sync' /> {t('button.zoom.init')}
                        </Button>
                        <Button variant='info' onClick={() => onZoomIn(trainChart.current)}>
                            <i className='fas fa-plus' /> {t('button.zoom.in')}
                        </Button>
                        <Button variant='info' onClick={() => onZoomOut(trainChart.current)}>
                            <i className='fas fa-minus' /> {t('button.zoom.out')}
                        </Button>
                    </ButtonArea>
                    <Row sm='7'>
                        <Col>
                            <LabelInput title={t('ui.train.threshold.k')} name='threshold_k' value={input.threshold} onChange={handleChangeThreshold} errors={{ hasError: false }} />
                        </Col>
                        <Col style={{ height: '41px' }}>
                            <LabelInput title='Error (%)' name='error' value={input.error} onChange={handleChangeError} errors={{ hasError: false }} />
                        </Col>
                    </Row>
                </ChartConfigurationContainer>
                <ChartArea>
                    <TrainChartArea>
                        <Line ref={trainChart} options={trainChartOption} data={chartData} redraw={true} plugins={plugins} />
                    </TrainChartArea>
                    <GaussianChartArea>
                        {gaussianInfo && <Line ref={gaussianChart} options={gaussianInfo.gaussianChartOption} data={chartData} redraw={true} plugins={plugins} />}
                    </GaussianChartArea>
                </ChartArea>
            </Row>
        </Col>
    )
}

export default TSADTrainContainer