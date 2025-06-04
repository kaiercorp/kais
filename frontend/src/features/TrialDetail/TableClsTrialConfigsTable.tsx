import { Col, Row } from 'react-bootstrap'

import { useEffect, useState } from 'react'
import PerfTable from './PerfTable'
import ConfigCol from './ConfigCol'
import TrialBase from './TrialBase'
import TableClsTrainInfo from './TableClsTrainInfo'


type TTrialConfigsTable = {
    trial: any,
    config: any,
    showPerfTable?: boolean
}

const TableClsTrialConfigsTable = ({ trial, config, showPerfTable }: TTrialConfigsTable) => {
    const [perf, setPerf] = useState<any>(null)
    useEffect(() => {
        setPerf(null)
        if (!trial) return
        if (trial.perf) {
            setPerf(trial.perf)
        } else if (trial.parent_trial && trial.parent_trial.perf) {
            setPerf(trial.parent_trial.perf)
        }
    }, [trial])

    if (!trial) return <Col />
    if (!config || !config.train_config) return <Col />

    return (
        <Row>
            <ConfigCol sm='4'>
                <TrialBase trial={trial} tconfig={config} />
            </ConfigCol>
            <ConfigCol sm='4'>
                <TableClsTrainInfo config={config} />
            </ConfigCol>
            { showPerfTable &&
                <ConfigCol sm='4'>
                    <PerfTable perfStr={perf} />
                </ConfigCol>
            }
        </Row>
    )
}

export default TableClsTrialConfigsTable
