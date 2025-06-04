import { useEffect, useState } from 'react'
import { Col, Row } from 'react-bootstrap'

import PerfTable from './PerfTable'
import TrialBase from './TrialBase'
import ConfigCol from './ConfigCol'

type TTrialConfigsTable = {
  trial: any
  testState?: string
}

const VisionClsTrialConfigsTable = ({ trial, testState }: TTrialConfigsTable) => {
  const [perf, setPerf] = useState<any>(null)
  useEffect(() => {
    if (!trial) return
    if (trial.perf) {
      setPerf(trial.perf)
    }
  }, [trial])

  if (!trial || !trial.params) return <Col />

  const config = trial.params
  if (!config || !config.train_config) return <Col />

  return (
    <Row>
      <ConfigCol sm='4'>
        <TrialBase trial={trial} tconfig={config} testState={testState} />
      </ConfigCol>
      {/* <ConfigCol sm='4'>
        <VisionClsTrainInfo config={config} />
      </ConfigCol> */}
      {
        (perf && !testState) && (
          <ConfigCol sm='4'>
            <PerfTable perfStr={perf} />
          </ConfigCol>
        )
      }
    </Row>
  )
}

export default VisionClsTrialConfigsTable
