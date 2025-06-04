import { useEffect, useState } from "react"
import { Card, Form } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import styled from "styled-components"

const StyledCardBody = styled(Card.Body)`
    padding: 5px;
`

const ScrollTable = styled.table`
    display: block;
    border-collapse: collapse;
    color: white;
    font-size: 11px;
    overflow: auto;

    & thead {
        border: 1px solid #000;
    }

    & tbody {
        display: block;
        height: 180px;
        width: 100%;
    }

    & th {
        border: 1px solid #000;
        text-align: center;
        padding: 5px;
    }

    & td {
        border: 1px solid #000;
        text-align: center;
        padding: 5px;
    }

    & th:nth-child(1),
    & td:nth-child(1) {
        min-width: 40px;
    }

    & th:nth-child(2),
    & td:nth-child(2) {
        min-width: 60px;
    }

    & th:nth-child(n+3),
    & td:nth-child(n+3) {
        min-width: 70px;
        word-break: break-all;
    }
    
`

const DataTable = ({ dataset, selectedData, selectData }: any) => {
    const [t] = useTranslation('translation')

    const [header, setHeader] = useState<any>([])
    const [rows, setRows] = useState<any>([])

    useEffect(() => {
        if (!dataset) return
        setHeader(dataset[0])
        setRows(dataset.slice(1))
    }, [dataset])

    return (
        <StyledCardBody>
            <ScrollTable>
                <thead>
                    <tr>
                        <th>
                            {t('ui.test.selectdata')}
                        </th>
                        {
                            header && header.map((d: any) => {
                                return <th key={`datasetheader-${d}`}>{d}</th>
                            })
                        }
                    </tr>
                </thead>
                <tbody>
                    {
                        rows && rows.map((row: any, index: any) => {
                            return <tr key={`datasetrow-${index}`}>
                                <td>
                                    <Form.Check
                                        type='radio'
                                        id={`dataset-row${index}`}
                                        onChange={() => selectData(index)}
                                        value={index}
                                        checked={selectedData === index}
                                    />
                                </td>
                                {
                                    row.map((r: any, jndex: any) => {
                                        return <td key={`datasetrow-data-${index}-${jndex}`}>{r}</td>
                                    })
                                }
                            </tr>
                        })
                    } 
                </tbody>
            </ScrollTable>
        </StyledCardBody>
    )
}

export default DataTable