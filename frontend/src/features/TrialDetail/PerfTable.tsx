import { PopoverLabel } from 'components'
import ConfigCol from './ConfigCol'

import styled from 'styled-components'
import { Col } from 'react-bootstrap'

import { useTranslation } from 'react-i18next'
import { useEffect, useState } from 'react'

type TResultOrder = {
    [key: string]: number
}
const resultOrder: TResultOrder = {
    wa: 1,
    auroc: 2,
    prauc: 3,
    f1: 4,
    precision: 5,
    recall: 6,
    uwa: 7,
    mae: 8,
    mse: 9,
    rmse: 10,
    image_accuracy: 11, 
    image_f1_score: 12,
    image_precision: 13,
    image_recall: 14,
    label_accuracy: 15,
    label_f1_score: 16,
    label_precision: 17,
    label_recall: 18,
  }

const TestDetailConfigCol = styled(Col)`
  font-size: 12px;
  & table {
    width: 100%;
    border: 1px solid grey;
    border-collapse: collapse;

    & th {
      font-weight: 600;
      text-align: center;
      padding: 5px;
      border: 1px solid grey;
      color: #ffffff;
      height: 47px;
    }

    & td {
      text-align: center;
      padding: 5px 10px;
      border: 1px solid grey;
      height: 47px;
    }
  }
`

const PerfTable = ({ perfStr, title, isTestDetail }: any) => {
    const [t] = useTranslation('translation')
    
    const [perf, setPerf] = useState<any>()
    const [result_keys, setResultKeys] = useState<any>()
    
    const PerfTableConfigCol = isTestDetail ? TestDetailConfigCol : ConfigCol

    useEffect(() => {
        if (perfStr) {
            let _perf = JSON.parse(perfStr)
            let _result_keys = Object.keys(_perf)
            _result_keys.sort(function(x, y) {
                const xidx = resultOrder[x] || 999
                const yidx = resultOrder[y] || 999
                return xidx - yidx
            })
            setResultKeys(_result_keys)
            setPerf(_perf)
        }
    }, [perfStr])

    if (!perf) return <></>
    return (
        <PerfTableConfigCol>
            <table>
                <tbody>
                    <tr>
                        <th colSpan={2}><b>{title?title:'Best 모델 성능'}</b></th>
                    </tr>
                    {
                        result_keys.map((key:any) => {
                            if(key.startsWith('b_') || key.startsWith('c_') || key.startsWith('Epoch') || key.startsWith('Model')) return null
                            return (
                                <tr key={`targets-${key}`}>
                                    <th>{t(`metric.${key}`)} <PopoverLabel name={`${key}`}>{t(`metric.${key}.desc`)}</PopoverLabel></th>
                                    <td>
                                        {['f1', 'precision', 'recall', 'uwa', 'wa', 'auroc', 'prauc', 'image_accuracy', 'image_precision', 'image_recall', 'image_f1_score', 'label_accuracy', 'label_precision', 'label_recall', 'label_f1_score'].includes(key)?Number(perf[key] * 100 || 0).toFixed(2)+'%':Number(perf[key]).toFixed(3)}
                                    </td>
                                </tr>
                            )
                        })
                    }
                </tbody>
            </table>
        </PerfTableConfigCol>
    )
}

export default PerfTable