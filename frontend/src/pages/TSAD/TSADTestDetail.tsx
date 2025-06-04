import { useEffect, useRef, useState, useContext } from "react"
import { Button, Card, Col, Row } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import { useLocation, Link } from "react-router-dom"
import styled from 'styled-components'

import { TableClsTrialConfigsTable, PerfTable } from "features"
import { logger, objDeepCopy, ApiFetchTrial, ApiDownloadFile } from "helpers"
import { ModelsArea, TestModelCard, TrainModelContainer } from "features/TrainModel"
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
import { LocationContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'
import { useSocket } from 'hooks'
import { engine } from "appConstants/trial"

const CardHeaderLeft = styled.div`
float: left;
color: #ffffff;
font-weight: 600;
`

const CardHeaderRight = styled.div`
float: right;
`

const ChartArea = styled(Row)`
height: 500px;
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

const chartOption = {
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
                font: { size: 12, weight: '400' },
            },
        },
        title: {
            display: true,
            text: 'Train Chart',
            color: '#000000',
            font: { size: 14, weight: '600', },
        },
        zoom: {
            zoom: {
                wheel: { enabled: true, },
                mode: "x" as const
            },
            pan: { enabled: true, mode: 'y' as const },
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
        x: {
            axis: 'x' as const,
            position: 'bottom' as const,
            title: { display: true, text: 'Time', },
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
            title: { display: true, text: 'TRUE', color: 'blue'},
            grid: {
                drawOnChartArea: false,
            }
        },
        y1: {
            type: 'linear' as const,
            display: true,
            position: 'right' as const,
            title: {
                display: true,
                text: 'Train Value',
                color: 'red'
            }
        }
    },
}

const initAnnotation = (option: any) => {
    if (!option) return
    option.plugins.annotation.annotations = []
}

const addAnnotation = (option: any, label: string, value: any) => {
    option.plugins.annotation.annotations.push({
        type: 'line' as const,
        yMin: value,
        yMax: value,
        borderColor: 'black',
        borderWidth: 1,
        borderDash: [6,6],
        label: {
            display: true,
            backgroundColor: 'rgba(200,200,200,0.8)',
            drawTime: 'afterDatasetsDraw',
            content: label,
            position: 'start' as const
        },
    })
}

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

const refineChartData = (chartData: any, chartRef: any) => {
    if (!chartData || chartData.length < 1) return

    let newDatasets = Array<any>()

    let data = JSON.parse(chartData[0]['multi_samples_test_material']['String'])
    chartRef.current.data.labels = Array.from({ length: data.pred.length }, (v, i) => i)
    newDatasets.push({
        label: `True`,
        data: data.true,
        borderColor: `rgb(53, 162, 235)`,
        backgroundColor: 'rgba(53, 162, 235, 0.5)',
        pointRadius: 0,
        tension: 0.1,
        yAxisID: 'y'
    })
    addAnnotation(chartOption, `Threshold ${data.threshold}`, data.threshold)

    newDatasets.push({
        label: `Pred`,
        data: data.pred,
        borderColor: `rgb(255, 165, 0)`,
        backgroundColor: 'rgba(255, 165, 0, 0.5)',
        pointRadius: 0,
        tension: 0.1,
        yAxisID: 'y1'
    })

    newDatasets.push({
        label: `Score`,
        data: data.score,
        borderColor: `rgb(255, 99, 132)`,
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
        pointRadius: 0,
        tension: 0.1,
        yAxisID: 'y'
    })

    chartRef.current.data.datasets = newDatasets
    chartRef.current.update()
}

const TSADTestDetail = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()
    
    const [trialId, setTrialId] = useState<number | undefined>()
    const [prevLocation, setPrevLocation] = useState(location.pathname)
    const [selectedModel, setSelectedModel] = useState<any>()
    const [socketConnected, setSocketConnected] = useState(false)
    
    const { updateLocationContextValue } = useContext(LocationContext) 
    
    const { trial } = ApiFetchTrial(trialId) 
    const queryClient = useQueryClient()

    const chart: any = useRef<ChartJS | null>(null)

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
    
    useEffect(() => {
        const pathVariables = location.pathname.split('/')
        setPrevLocation(pathVariables.slice(0, pathVariables.length - 2).join('/'))
        setTrialId(Number(pathVariables.slice(-2, -1)[0]))
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const [models, setModels] = useState<any[]>()
    useEffect(() => {
        if (!trial || trial.project_id === 0) {
            setModels([])
            return
        }

        logger.log(`Change Location to ${t('title.ts.ad.test', { trial: trial.trial_name })}`)
        updateLocationContextValue({location: 'ts.ad.test', locationValue: { trial: `[${trial.trial_id}] ${trial.trial_name}` }})

        if (trial.test && trial.test.models) {
            let _models = objDeepCopy(trial.test.models)

            setModels(_models)
            setSelectedModel(_models[0])
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial])

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            if (e.data === 'Invalid request') {
                throw e
            }

            let msg = JSON.parse(e.data)
            if (msg.length > 0) {
                refineChartData(msg, chart)
            }
        } catch (e) {
            logger.error(e)
        }
    }

    const ws = useSocket(`/trials/test/result/${trial?.test?.id}`, 'Result Data', handleSocketMessage, {setSocketConnected, shouldCleanup: true, shouldConnect: !!selectedModel}) 

    useEffect(() => {
        if (!selectedModel) return 

        initChartData(chart)
        initAnnotation(chartOption)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedModel])

    useEffect(() => {
        if (!socketConnected || !ws || !ws.current) {

        } else {
            ws.current.send(JSON.stringify({ sample_id: trial.trial_id, engine_type: engine.ts_ad }))
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [socketConnected])

    const onSelectModel = (model: any) => {
        setSelectedModel(model)
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ model_name: model.name }))
    }

    const handleDownloadTestFiles = () => {
        ApiDownloadFile(queryClient, `/trial/test/download/${selectedModel.id}`, `${trial.trial_name}_${selectedModel.model}.zip`)
    }

    return (
        <Col>
            <Row>
                <Col>
                    <Card>
                        <Card.Header>
                            <CardHeaderLeft>{t('ui.train.title.info')}</CardHeaderLeft>
                            <CardHeaderRight><Link to={prevLocation}>{t('button.go.list')}</Link></CardHeaderRight>
                        </Card.Header>
                        <Card.Body>
                            {
                                trial && trial.parent_trial && <TableClsTrialConfigsTable trial={trial} config={trial.parent_trial.params} />
                            }
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            <Row>
                <Col>
                    <Card>
                        <TrainModelContainer isVertical={false}>
                            <ModelsArea isVertical={false}>
                                <div style={{ display: 'flex' }}>
                                    {models && models.map((model: any) => {
                                        return (
                                            <TestModelCard
                                                key={`model-list-${model.id}`}
                                                model={model}
                                                isSelected={selectedModel && (model.id === selectedModel.id)}
                                                onClick={() => onSelectModel(model)}
                                            />
                                        )
                                    })}
                                </div>
                            </ModelsArea>
                        </TrainModelContainer>
                    </Card>
                    <Card>
                        <Card.Body>
                            {
                                selectedModel && (
                                    <Row style={{ marginBottom: '10px' }}>
                                        <Col sm={2}>
                                            <TrainModelContainer isVertical={false}>
                                                <PerfTable perfStr={selectedModel.perf.String} title={t('ui.test.title.perf', { perf: selectedModel.name })} isTestDetail={true} />
                                            </TrainModelContainer>
                                            <Row style={{ marginTop: '10px' }}>
                                                <Col sm={2}></Col>
                                                <Col><Button style={{ width: '140px', height: '30px' }} onClick={handleDownloadTestFiles}>{t('button.download.resultfile')}</Button></Col>
                                                <Col sm={2}></Col>
                                            </Row>
                                        </Col>
                                    </Row>
                                )
                            }
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
                                <ChartArea>
                                    <Line ref={chart} options={chartOption} data={chartData} redraw={true} plugins={plugins} />
                                </ChartArea>
                            </Row>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Col>
    )
}

export default TSADTestDetail