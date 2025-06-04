import { useEffect, useState } from 'react'
import styled from 'styled-components'

import {Table} from 'react-bootstrap'

const TestDataResultLayout = styled.div`
    max-height: 600px;
    overflow: scroll;
`

const TestDataResultTR = styled.tr<{isInCorrect: boolean, isBlind: boolean}>`
    ${(props) => props.isInCorrect && 'background-color: white; color: black;'}
`

const cfg = ['tp', 'tn', 'fp', 'fn']
const perf = ['accuracy', 'precision', 'recall', 'f1_score']

const TestDataResult = ({ resultList, onClick }: any) => {
    const [header, setHeader] = useState<any>()

    useEffect(() => {
        if (!resultList || !resultList[0] || !resultList[0].data) return

        resultList.forEach((result:any) => {
            result.props = JSON.parse(result.data)
        })
        let _header = Object.keys(resultList[0].props)
        _header = _header.filter((h:string) => {
            if (['kaier_id', 'sr.true', 'sr.prediction', 'sr.true_label', 'sr.predicted_label'].includes(h)) return false
            return true
        })
        
        if (resultList[0].save_path?.includes('vision-cls-ml')) {
            _header = [..._header.filter((header: string) => !(cfg.includes(header) || perf.includes(header))).sort()]           
        }

        setHeader(_header)
    }, [resultList])

    return (
        <TestDataResultLayout>
            <div style={{overflow: 'auto', height: '500px'}}>
                <Table>
                <thead style={{backgroundColor: '#37404a', position: 'sticky', top:'0'}}>
                    <tr>
                        <th>No.</th>
                        <th>Data</th>                        
                        <th>Ground Truth</th>                        
                        <th>Prediction</th>                        
                        <th>Result</th>                        
                        {
                            header && header.map((h: any) => {
                                return <th style={{position: 'sticky', top:'0'}} key={`headeritem-${h}`}>{h}</th>
                            })
                        }
                    </tr>
                </thead>
                <tbody>
                    {resultList && resultList.map((result: any, index: number) => {
                        let isInCorrect = false
                        let isBlind = false
                        const sorted_predicted = result.props['sr.predicted_label'].split(',').sort().join()
                        const sorted_true = result.props['sr.true_label'].split(',').sort().join()
                        if (result.props['sr.true_label'] && result.props['sr.predicted_label'] && (sorted_predicted !== sorted_true)) isInCorrect = true
                        if (result.props['sr.true_label'] && result.props['sr.predicted_label'] && (sorted_true === '-')) isBlind = true
                        return (
                            <TestDataResultTR key={`datarow-${index}-${result.props.kaier_id}`} onClick={() => onClick(result.id)} isInCorrect={isInCorrect} isBlind={isBlind}>
                                <td>{index + 1}</td>
                                <td>{result.props.kaier_id}</td>
                                <td>{sorted_true}</td>
                                <td>{sorted_predicted}</td>
                                <td>{isBlind? '-':(isInCorrect ? 'X' : 'O')}</td>
                                {
                                    header&&header.map((h: any) => {
                                        return (<td key={`dataitem-${result.props.kaier_id}-${h}-${index}`}>{(result.props[h]*100||0).toFixed(2)+'%'}</td>)
                                    })
                                }
                            </TestDataResultTR>
                        )
                    })}
                </tbody>
            </Table>
        </div>
        </TestDataResultLayout>
    )
}

export default TestDataResult
