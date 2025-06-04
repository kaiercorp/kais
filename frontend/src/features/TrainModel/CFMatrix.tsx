import { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

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
  border-right: 1px solid grey;
  background: transparent;
  min-height: ${(props) => (props.header ? '60px' : '35px')};
  display: flex;
`

const LeftHeaderVertical = styled.div`
  width: 40px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  transform: rotate(270deg);
`

const LeftHeaderSubGroup = styled.div`
  width: 120px;
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
`

const RightArea = styled.div`
  width: calc(100% - 160px);
`

interface IRightHeader {
    length?: string | number
}
const RightHeader = styled.div<IRightHeader>`
  width: ${props => Number(props.length) * 80}px;
  max-width: 100%;
  border: 1px solid grey;
  height: 35px;
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
  width: 80px;
  height: ${(props) => (props.header ? '60px' : '35px')};
  padding: 0 5px;
  background-color: ${(props) =>
        props.hover
            ? props.highlight
                ? 'skyblue'
                : '#CCFF66'
            : props.highlight
                ? 'skyblue'
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

const CFMatrix = ({ cf_matrixStr }: any) => {
    const [t] = useTranslation('translation')
    const [cellInfo, setCellInfo] = useState({ name: '', idx: -1 })

    const [cf_matrix, setCf_matrix] = useState<any>()
    useEffect(() => {
        if (!cf_matrixStr) return
        if (typeof(cf_matrixStr) === 'string') {
            const _cf_matrix = JSON.parse(cf_matrixStr)
            setCf_matrix(_cf_matrix)
        } else {
            setCf_matrix(cf_matrixStr)
        }
    }, [cf_matrixStr])

    const setCellName = (header: string, idx: number) => {
        setCellInfo({ name: header, idx: idx })
    }

    const resetCellName = () => {
        setCellInfo({ name: '', idx: -1 })
    }

    if (!cf_matrix) return <></>

    const row_keys = Object.keys(cf_matrix).filter((key:string) => {
        if (['accuracy', 'macro avg', 'weighted avg'].includes(key)) return false
        if (['cm.total', 'cm.precision'].includes(key)) return false
        return true
    })
    const col_keys = Object.keys(cf_matrix[row_keys[0]]).filter((key:string) => {
        if (['accuracy', 'macro avg', 'weighted avg'].includes(key)) return false
        if (['cm.total', 'cm.recall'].includes(key)) return false
        return true
    })

    // left headers
    let leftHeader = row_keys.map(function (key:string, index) {
        return (
            <LeftHeaderSub key={`${key}-leftHeader`} position={index === 0 ? 'first' : ''}>
                {key}
            </LeftHeaderSub>
        )
    })
    leftHeader.push(<LeftHeaderSub key='sum-leftHeader'>{t('ui.label.total')}</LeftHeaderSub>)
    leftHeader.push(<LeftHeaderSub key='precision-leftHeader'>{t('ui.label.precision')}</LeftHeaderSub>)

    let columns = col_keys.map(function(header:string, index) {
        const col = (
            <ResultColumn key={header}>
                <ResultCell
                    key={`${header}_col_header_row`}
                    title={header}
                    hover={cellInfo.name === header}
                    header={true}
                >
                    {header}
                </ResultCell>
                {row_keys.map((key, jndex) => {
                    return (
                        <ResultCell
                            key={`${header}_col_${key}_row`}
                            highlight={index === jndex}
                            onMouseEnter={() => setCellName(header, jndex)}
                            onMouseLeave={() => resetCellName()}
                            hover={cellInfo.idx === jndex || cellInfo.name === header}
                        >
                            {cf_matrix[key][header]}
                        </ResultCell>
                    )
                })}
                <ResultCell key={`${header}_col_sum_row`} hover={header === cellInfo.name}>
                    {cf_matrix['cm.total'][header]}
                </ResultCell>
                <ResultCell key={`${header}_col_precision_row`} hover={header === cellInfo.name}>
                    {Number(cf_matrix['cm.precision'][header].split("%")[0]).toFixed(2)}%
                </ResultCell>
            </ResultColumn>
        )

        return col
    })

    const sumCol = (
        <ResultColumn key='sum-col'>
            <ResultCell header={true}>{t('ui.label.total')}</ResultCell>
            {row_keys.map((key, jndex) => {
                return (
                    <ResultCell key={`${key}_sum_col_row`} hover={cellInfo.idx === jndex}>
                        {cf_matrix[key]['cm.total']}
                    </ResultCell>
                )
            })}
            <ResultCell>{cf_matrix['cm.total']['cm.total']}</ResultCell>
            <ResultCell>-</ResultCell>
        </ResultColumn>
    )
    columns.push(sumCol)

    const resultCol = (
        <ResultColumn key='precision-col'>
            <ResultCell header={true}>{t('ui.label.recall')}</ResultCell>
            {row_keys.map((key, jndex) => {
                return (
                    <ResultCell key={`${key}_row_precision`} hover={cellInfo.idx === jndex}>
                        {Number(cf_matrix[key]['cm.recall'].split("%")[0]).toFixed(2)}%
                    </ResultCell>
                )
            })}
            <ResultCell>-</ResultCell>
            <ResultCell>-</ResultCell>
        </ResultColumn>
    )
    columns.push(resultCol)

    return (
        <ResultTableContainer>
            <LeftHeaderArea>
                <LeftHeaderGroup border='none'></LeftHeaderGroup>
                <LeftHeaderGroup border='none' header={true}></LeftHeaderGroup>
                <LeftHeaderGroup>
                    <LeftHeaderVertical>{t('ui.label.ground_truth')}</LeftHeaderVertical>
                    <LeftHeaderSubGroup>{leftHeader}</LeftHeaderSubGroup>
                </LeftHeaderGroup>
            </LeftHeaderArea>
            <RightArea>
                <RightHeader length={columns.length}>{t('ui.label.prediction')}</RightHeader>
                <ContentArea>{columns}</ContentArea>
            </RightArea>
        </ResultTableContainer>
    )
}

export default CFMatrix
