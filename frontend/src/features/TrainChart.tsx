import { useEffect, useRef, useState } from 'react'
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
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import zoomPlugin from 'chartjs-plugin-zoom'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { convertUtcTimeToLocalTimeMonthly, ApiFetchChartData, ApiFetchTrial } from 'helpers'

const ChartArea = styled(Row)`
  height: 400px;
`

const ButtonArea = styled.div`
  margin-bottom: 5px;
`

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler, zoomPlugin, {
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

const chartTitle = (tooltipItems: any) => {
  return 'Epoch ' + tooltipItems[0].label.split(',')[0]
}

const chartFooter = (tooltipItems: any) => {
  return tooltipItems[0].label.split(',')[1]
}

const getLossLabelForValue = (val: any) => {
  return lossDatasets.labels[val][0]
}

const getAccuracyLabelForValue = (val: any) => {
  return accuracyDatasets.labels[val][0]
}

const lossOptions = {
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
        font: {
          size: 14,
          weight: '400',
        },
      },
    },
    title: {
      display: true,
      text: 'Loss',
      color: '#000000',
      font: {
        size: 18,
        weight: '600',
      },
    },
    zoom: {
      zoom: {
        wheel: { enabled: true, },
        mode: "x" as const
      },
      pan: { enabled: true, mode: 'y' as const },
      limits: {
        y: {
          min: -20,
          max: 120,
        },
      },
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
  },
  scales: {
    x: {
      axis: 'x' as const,
      position: 'bottom' as const,
      title: {
        display: true,
        text: 'Epoch',
      },
      ticks: {
        callback: (val: any, index: number) => {
          return getLossLabelForValue(val)
        },
      },
    },
    y: {
      axis: 'y' as const,
      position: 'left' as const,
      title: {
        display: true,
        text: 'Loss',
      },
    },
  },
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
        font: {
          size: 14,
          weight: '400',
        },
      },
    },
    title: {
      display: true,
      text: 'Accuracy',
      color: '#000000',
      font: {
        size: 18,
        weight: '600',
      },
    },
    zoom: {
      zoom: {
        wheel: { enabled: true, },
        mode: "x" as const
      },
      pan: { enabled: true, mode: 'y' as const },
      limits: {
        y: {
          min: -20,
          max: 120,
        },
      },
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
  },
  scales: {
    x: {
      axis: 'x' as const,
      position: 'bottom' as const,
      title: {
        display: true,
        text: 'Epoch',
      },
      ticks: {
        callback: (val: any, index: number) => {
          return getAccuracyLabelForValue(val)
        },
      },
    },
    y: {
      axis: 'y' as const,
      position: 'left' as const,
      title: {
        display: true,
        text: 'Accuracy (%)',
      },
    },
  },
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

type TTrainChart = {
  show: boolean
  isRun: boolean
  trial: any
}

const TrainChart = ({ show, isRun, trial }: TTrainChart) => {
  const [t] = useTranslation('translation')

  const [lastEpoch, setLastEpoch] = useState(0)
  const [checkStop, setCheckStop] = useState(false)
  const [epoch, setEpoch] = useState<number | undefined>()
  const [trialId, setTrialId] = useState<number | undefined>()

  const { chartData } = ApiFetchChartData(trial.trial_id, epoch)
  ApiFetchTrial(trialId)

  const handleCheckStop = () => {
    setCheckStop(!checkStop)
  }

  useEffect(() => {
    if (isRun) {
      setTrialId(trial.trial_id)

      setTimer(timer + 1)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [checkStop, isRun])

  const [timer, setTimer] = useState(1)

  useEffect(() => {
    if (!show) {
      return
    }

    if (isRun) {
      let timer = setTimeout(() => {
        changeQueryKeyOfFetchChartData(lastEpoch)
      }, 2000)

      return () => {
        clearTimeout(timer)
      }
    } else {
      changeQueryKeyOfFetchChartData(lastEpoch)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [timer])

  // 상세보기로 진입한 경우, 차트 데이터 한 번 호출
  useEffect(() => {
    if (show && trial.trial_id) {
      changeQueryKeyOfFetchChartData(0)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [show, trial.trial_id])

  const chartRefLoss: any = useRef(null)
  const handleResetZoomLoss = () => {
    if (!chartRefLoss || !chartRefLoss.current) return
    chartRefLoss.current.resetZoom()
  }

  const chartRefAccuracy: any = useRef<ChartJS | null>(null)
  const handleResetZoomAccuracy = () => {
    if (!chartRefAccuracy || !chartRefAccuracy.current) return
    chartRefAccuracy.current.resetZoom()
  }

  const changeQueryKeyOfFetchChartData = (epoch: number) => {
    setEpoch(epoch)
  }

  useEffect(() => {
    if (chartData !== null && chartData.length > 0) {
      if (lastEpoch === 0) {
        if (chartRefLoss && chartRefLoss.current) {
          chartRefLoss.current.data.datasets[0].data = [0]
          chartRefLoss.current.data.datasets[1].data = [0]

          chartRefLoss.current.update()
        }

        if (chartRefAccuracy && chartRefAccuracy.current) {
          chartRefAccuracy.current.data.datasets[0].data = [0]
          chartRefAccuracy.current.data.datasets[1].data = [0]

          chartRefAccuracy.current.update()
        }
      }

      let lastIdx = 0
      for (let index = 0; index < chartData.length; index++) {
        if (
          chartRefLoss.current.data.labels[chartRefLoss.current.data.labels.length - 1][0] <
          parseInt(chartData[index].cur_epoch)
        ) {
          const label = [
            parseInt(chartData[index].cur_epoch),
            convertUtcTimeToLocalTimeMonthly(chartData[index].utc_ts * 1000),
          ]
          chartRefLoss.current.data.labels.push(label)
          chartRefAccuracy.current.data.labels.push(label)
        }

        if (chartRefLoss && chartRefLoss.current) {
          chartRefLoss.current.data.datasets[0].data.push(chartData[index].loss.toFixed(2))
          chartRefLoss.current.data.datasets[1].data.push(chartData[index].valid_loss.toFixed(2))
        }

        if (chartRefAccuracy && chartRefAccuracy.current) {
          chartRefAccuracy.current.data.datasets[0].data.push(Number((chartData[index].acc * 100).toFixed(2)))
          chartRefAccuracy.current.data.datasets[1].data.push(Number((chartData[index].valid_acc * 100).toFixed(2)))
        }

        lastIdx = parseInt(chartData[index].cur_epoch)
      }

      while (lastIdx < chartRefLoss.current.data.labels[chartRefLoss.current.data.labels.length - 1][0]) {
        chartRefLoss.current.data.labels.pop()
        chartRefAccuracy.current.data.labels.pop()
      }

      setLastEpoch(lastIdx)

      if (chartRefLoss && chartRefLoss.current) {
        chartRefLoss.current.update()
      }

      if (chartRefAccuracy && chartRefAccuracy.current) {
        chartRefAccuracy.current.update()
      }
    } else {
      handleCheckStop()
    }

    setTimer(timer + 1)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chartData])

  return (
    <Row>
      <Col>
        <ButtonArea>
          <Button variant='info' onClick={handleResetZoomLoss}>
            <i className='fas fa-sync' /> {t('button.zoom.init')}
          </Button>
        </ButtonArea>
        <ChartArea>
          <Line ref={chartRefLoss} data={lossDatasets} options={lossOptions} redraw={true} height={400} />
        </ChartArea>
      </Col>
      <Col>
        <ButtonArea>
          <Button variant='info' onClick={handleResetZoomAccuracy}>
            <i className='fas fa-sync' /> {t('button.zoom.init')}
          </Button>
        </ButtonArea>
        <ChartArea>
          <Line ref={chartRefAccuracy} data={accuracyDatasets} options={accuracyOptions} redraw={true} height={400} />
        </ChartArea>
      </Col>
    </Row>
  )
}

export default TrainChart
