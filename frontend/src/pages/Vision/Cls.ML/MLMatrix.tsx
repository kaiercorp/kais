import { useEffect, useState } from 'react'
import styled from 'styled-components'

const ResultTableContainer = styled.div`
  background-color: #ffffff;
  padding: 10px;
  display: flex;
  color: #5f5f5f;
`

const LeftHeaderArea = styled.div`
  width: 160px;
`

interface ILeftHeaderGroup {
    border?: string
    header?: boolean
}
const LeftHeaderGroup = styled.div<ILeftHeaderGroup>`
  margin: 0;
  border: ${(props) => props.border || '1px solid grey'};
  background: transparent;
  min-height: 48px;
  display: flex;
`

const LeftHeaderSubGroup = styled.div`
  width: 160px;
  margin: 0;
`

interface ILeftHeaderSub {
    position?: string
}
const LeftHeaderSub = styled.div<ILeftHeaderSub>`
  width: 100%;
  border: 1px solid grey;
  border-right: none;
  border-top: ${(props) => (props.position === 'first' && 'none') || '1px solid grey'};
  border-bottom: none;
  margin: 0;
  padding: 0 3px;
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  line-height: 34px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  word-break: break-all;
  text-align: center;
  height: 47px;
`

const RightArea = styled.div`
  width: calc(100% - 160px);
`

interface IRightHeader {
    length?: string | number
}
const RightHeader = styled.div<IRightHeader>`
  width: ${props => Number(props.length) * 87}px;
  max-width: 100%;
  border: 1px solid grey;
  height: 0px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
`
const ContentArea = styled.div`
  width: 100%;
  overflow-x: scroll;
  display: flex;
`

interface IResultCell {
    header?: boolean
    hover?: boolean
    highlight?: boolean
}
const ResultCell = styled.div<IResultCell>`
  font-size: 11px;
  border: 1px solid grey;
  border-top: 1px solid transparent;
  width: 87px;
  height: 47px;
  padding: 0 5px;
  background-color: ${(props) =>
        props.hover
            ? props.highlight
                ? '#CCFF66'
                : '#CCFF66'
            : 'transparent'};
  text-overflow: ellipsis;
  white-space: normal;
  word-break: break-all;
  text-align: ${(props) => (props.header ? 'center' : 'right')};
  font-weight: ${(props) => (props.header ? '600' : '400')};
  display: flex;
  align-items: center;
  justify-content: ${props => props.header ? 'center' : 'right'};
`

const ResultColumn = styled.div`
  border-top: 1px solid transparent;
`

const cfg = ['tp', 'tn', 'fp', 'fn']
const perf = ['accuracy', 'f1_score', 'precision', 'recall']
const default_col_keys = [...cfg, 'sums', ...perf]

const MLMatrix = ({ ml_matrixStr }: any) => {
    const [cellInfo, setCellInfo] = useState({ name: '', idx: -1 })

    const [ml_matrix, setMl_matrix] = useState<any>()
    useEffect(() => {
        if (!ml_matrixStr) return
        if (typeof(ml_matrixStr) === 'string') {
            const _ml_matrix = JSON.parse(ml_matrixStr)
            setMl_matrix(_ml_matrix)
        } else {
            setMl_matrix(ml_matrixStr)
        }
    }, [ml_matrixStr])

    const setCellName = (header: string, idx: number) => {
        setCellInfo({ name: header, idx: idx })
    }

    const resetCellName = () => {
        setCellInfo({ name: '', idx: -1 })
    }

    if (!ml_matrix) return <></>

    const row_keys = Object.keys(ml_matrix).sort((a: string, b: string) => {
        if (a === 'sum') return 1
        if (b === 'sum') return -1
        
        return a < b ? -1 : 1
    })
    // const col_keys = Object.keys(ml_matrix[row_keys[0]])
    const col_keys = default_col_keys

    let leftHeader = row_keys.map(function (key:string, index) {
        return (
            <LeftHeaderSub key={`${key}-leftHeader`} position={index === 0 ? 'first' : ''}>
                {key}
            </LeftHeaderSub>
        )
    })

    let columns = col_keys.map(function(header:string, index) {
        const col = (
            <ResultColumn key={header}>
                <ResultCell
                    key={`${header}_col_header_row`}
                    title={header}
                    hover={cellInfo.name === header}
                    header={true}
                >
                    {perf.includes(header) ? `Label ${header} (%)` : header}
                </ResultCell>
                {row_keys.map((key, jndex) => {
                    if (key === 'sum') {console.log(ml_matrix[key][header])}
                    return (
                        <ResultCell
                            key={`${header}_col_${key}_row`}
                            highlight={index === jndex}
                            onMouseEnter={() => setCellName(header, jndex)}
                            onMouseLeave={() => resetCellName()}
                            hover={cellInfo.idx === jndex || cellInfo.name === header}
                        >
                            {perf.includes(header) ? Number.isNaN(Number(ml_matrix[key][header])) ? ml_matrix[key][header] : (ml_matrix[key][header]*100).toFixed(2) : ml_matrix[key][header]}
                        </ResultCell>
                    )
                })}
            </ResultColumn>
        )

        return col
    })

    return (
        <ResultTableContainer>
            <LeftHeaderArea>
                <LeftHeaderGroup border='none' header={true}></LeftHeaderGroup>
                <LeftHeaderGroup>
                    <LeftHeaderSubGroup>{leftHeader}</LeftHeaderSubGroup>
                </LeftHeaderGroup>
            </LeftHeaderArea>
            <RightArea>
                <RightHeader length={columns.length}></RightHeader>
                <ContentArea>{columns}</ContentArea>
            </RightArea>
        </ResultTableContainer>
    )
}

export default MLMatrix
