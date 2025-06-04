import { useRef, useState, useCallback, useEffect, forwardRef, useContext } from 'react'
import { DndProvider, useDrag, useDrop } from 'react-dnd'
import { HTML5Backend } from 'react-dnd-html5-backend'
import update from 'immutability-helper'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { TrialContext } from 'contexts'
import { ApiFetchTrialCompare, QUERY_KEY } from 'helpers'
import { useQueryClient } from '@tanstack/react-query'

interface IDndColumn {
    opacity?: number
}

interface IDndCell {
    type?: string
    color: string
}

const DndColumn = styled.div<IDndColumn>`
    border: 1px solid grey;
    padding: 0;
    margin-btottom: 0.5rem;
    backgournd-color: white;
    cursor: move;
    opacity: ${(props) => props.opacity};
    max-width: 170px;
  `

const cellColor: any = {
    title: 'white',
    max: 'red',
    inherit: 'inherit',
}
const DndCell = styled.div<IDndCell>`
    width: 100%;
    border: 1px solid grey;
    margin: 0;
    text-align: center;
    font-weight: ${(props) => (props.type === 'title' ? '600' : '400')};
    line-height: 33px;
    padding: 0 5px;
    color: ${(props) => cellColor[props.color] || 'inherit'};
    text-overflow: ellipsis;
    white-space: nowrap;
    word-break: break-all;
    overflow: hidden;
    ${(props) => (props.type === 'title' || props.type === 'bottom' ? 'border-bottom: none;' : '')}
`

type DndColType = {
    id: any
    index: number
    model: any
    moveColumn: any
}

const DndCol = ({ id, index, model, moveColumn }: DndColType) => {
    const ref = useRef<HTMLDivElement>(null)
    const [{ handlerId }, drop] = useDrop({
        accept: 'column',
        collect(monitor) {
            return {
                handlerId: monitor.getHandlerId(),
            }
        },
        hover(item: any, monitor) {
            if (!ref.current) {
                return
            }
            const dragIndex = item.index
            const hoverIndex = index
            // Don't replace items with themselves
            if (dragIndex === hoverIndex) {
                return
            }
            // Determine rectangle on screen
            const hoverBoundingRect = ref.current?.getBoundingClientRect()
            // Get vertical middle
            const hoverMiddleX = (hoverBoundingRect.right - hoverBoundingRect.left) / 2
            // Determine mouse position
            const clientOffset = monitor.getClientOffset()
            // Get pixels to the top
            const hoverClientX = (clientOffset?.x || 0) - hoverBoundingRect.left
            // Only perform the move when the mouse has crossed half of the items height
            // When dragging downwards, only move when the cursor is below 50%
            // When dragging upwards, only move when the cursor is above 50%
            // Dragging downwards
            // 상하 dnd -> 좌우 dnd로 변경함
            if (dragIndex < hoverIndex && hoverClientX < hoverMiddleX) {
                return
            }
            // Dragging upwards
            if (dragIndex > hoverIndex && hoverClientX > hoverMiddleX) {
                return
            }
            // Time to actually perform the action
            moveColumn(dragIndex, hoverIndex)
            // Note: we're mutating the monitor item here!
            // Generally it's better to avoid mutations,
            // but it's good here for the sake of performance
            // to avoid expensive index searches.
            item.index = hoverIndex
        },
    })

    const [{ isDragging }, drag] = useDrag({
        type: 'column',
        item: () => {
            return { id, index }
        },
        collect: (monitor) => ({
            isDragging: monitor.isDragging(),
        }),
    })

    const opacity = isDragging ? 0 : 1
    drag(drop(ref))

    return (
        <DndColumn opacity={opacity} ref={ref} data-handler-id={handlerId}>
            <DndCell color='title' type='title'>
                {model.name} ({model.local_id})
            </DndCell>
            <DndCell color='inherit'>{model.model}</DndCell>
            <DndCell color='inherit'>{model.test_db}</DndCell>
            <DndCell color='inherit'>{model.num_gpus}</DndCell>
            <DndCell color={model.maxList.includes('avg_time') ? 'max' : 'inherit'}>
                {(model.avg_time ? Number(model.avg_time) * 1000 : 0).toFixed(0)}ms
            </DndCell>
            <DndCell color={model.maxList.includes('accuracy') ? 'max' : 'inherit'}>
                {(model.accuracy ? Number(model.accuracy) * 100 : 0).toFixed(2)}%
            </DndCell>
            <DndCell color={model.maxList.includes('recall') ? 'max' : 'inherit'}>
                {(model.recall ? Number(model.recall) * 100 : 0).toFixed(2)}%
            </DndCell>
            <DndCell color={model.maxList.includes('precision') ? 'max' : 'inherit'}>
                {(model.precision ? Number(model.precision) * 100 : 0).toFixed(2)}%
            </DndCell>
            <DndCell color={model.maxList.includes('f1') ? 'max' : 'inherit'} type='bottom'>
                {(model.f1 ? Number(model.f1) * 100 : 0).toFixed(2)}%
            </DndCell>
        </DndColumn>
    )
}

const DndContainer = styled.div`
    display: flex;
  `

const DndHeaderArea = styled.div`
    width: 100px;
  `

interface IDndHeaderGroup {
    border?: string
}

const DndHeaderGroup = styled.div<IDndHeaderGroup>`
    margin: 0;
    border: ${(props) => props.border || '1px solid grey'};
    border-right: none;
    background: transparent;
    min-height: 35px;
    display: flex;
  `

const DndHeaderSubArea = styled.div`
    width: 100px;
    margin: 0;
  `

interface IDndHeaderSub {
    position?: string
}

const DndHeaderSub = styled.div<IDndHeaderSub>`
    width: 100%;
    border: 1px solid grey;
    border-right: none;
    border-top: ${(props) => (props.position === 'first' && 'none') || '1px solid grey'};
    border-bottom: ${(props) => (props.position === 'last' && 'none') || '1px solid grey'};
    margin: 0;
    color: white;
    text-align: center;
    font-weight: 400;
    line-height: 33px;
  `

const DndContentArea = styled.div`
    display: flex;
    overflow-x: scroll;
    max-width: 570px;
  `

function resToData(source: any) {
    const dataLength = source.length
    if (dataLength < 1) return []

    let max = {
        accuracy: [0, 0],
        f1: [0, 0],
        precision: [0, 0],
        recall: [0, 0],
        avg_time: [0, 99999],
    }

    for (var ii = 0; ii < dataLength; ii++) {
        if (!source[ii].test_result || !source[ii].test_result.String) continue

        source[ii]['test'] = JSON.parse(source[ii].test_result.String)
        var tests1 = source[ii].test
        var models1 = Object.keys(tests1)
        for (var jj = 0; jj < models1.length; jj++) {
            var test1 = tests1[models1[jj]].summary
            if (test1['accuracy'] > max.accuracy[1]) {
                max.accuracy = [ii + jj, test1['accuracy']]
            }
            if (test1['macro avg']['f1-score'] > max.f1[1]) {
                max.f1 = [ii + jj, test1['macro avg']['f1-score']]
            }
            if (test1['macro avg']['precision'] > max.precision[1]) {
                max.precision = [ii + jj, test1['macro avg']['precision']]
            }
            if (test1['macro avg']['recall'] > max.recall[1]) {
                max.recall = [ii + jj, test1['macro avg']['recall']]
            }
            if (test1['avg_time'] < max.avg_time[1]) {
                max.avg_time = [ii + jj, test1['avg_time']]
            }
        }
    }

    let data = []
    for (var i = 0; i < dataLength; i++) {
        var tests = source[i].test
        var models = Object.keys(tests)
        for (var j = 0; j < models.length; j++) {
            let maxList = []
            if (max.accuracy[0] === i + j) maxList.push('accuracy')
            if (max.f1[0] === i + j) maxList.push('f1')
            if (max.precision[0] === i + j) maxList.push('precision')
            if (max.recall[0] === i + j) maxList.push('recall')
            if (max.avg_time[0] === i + j) maxList.push('avg_time')

            var test = tests[models[j]]
            let col = {
                key: source[i]['trial_id'] + '-' + i + '-' + j,
                id: source[i]['trial_id'],
                local_id: source[i]['trial_id'],
                name: source[i]['trial_name'],
                model: models[j],
                accuracy: test.summary['accuracy'],
                f1: test.summary['macro avg']['f1-score'],
                precision: test.summary['macro avg']['precision'],
                recall: test.summary['macro avg']['recall'],
                avg_time: test.summary['avg_time'],
                num_gpus: source[i]['gpus'],
                test_db: source[i]['test_db'],
                maxList: maxList,
            }
            data.push(col)
        }
    }

    return data
}

const CompareTrial = forwardRef(({ toggle }: any, ref) => {
    const [t] = useTranslation('translation')

    const [cols, setCols] = useState<any>(null)
    
    const { trialContextValue } = useContext(TrialContext)

    const queryClient = useQueryClient()
    const { compareTrials } = queryClient.getQueryData<any>([QUERY_KEY.fetchTrialCompare])
    const fetchTrialCompare = ApiFetchTrialCompare()

    useEffect(() => {
        if (typeof trialContextValue.selectedRows === 'undefined') return 

        fetchTrialCompare.mutate(trialContextValue.selectedRows.map((r: any) => r.trial_id))
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        if (!compareTrials) return
        setCols(resToData(compareTrials))
    }, [compareTrials])

    const moveColumn = useCallback((dragIndex: number, hoverIndex: number) => {
        setCols((prevCols: any) =>
            update(prevCols, {
                $splice: [
                    [dragIndex, 1],
                    [hoverIndex, 0, prevCols[dragIndex]],
                ],
            })
        )
    }, [])

    const renderCol = useCallback((col: any, index: number) => {
        return <DndCol key={col.key} index={index} id={col.id} model={col} moveColumn={moveColumn} />
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <DndProvider backend={HTML5Backend}>
            <DndContainer>
                <DndHeaderArea>
                    <DndHeaderGroup border='none'></DndHeaderGroup>
                    <DndHeaderGroup border='none'></DndHeaderGroup>
                    <DndHeaderGroup>
                        <DndHeaderSubArea>
                            <DndHeaderSub position='first'>{t('ui.formatter.testdb')}</DndHeaderSub>
                            <DndHeaderSub>{t('ui.formatter.usedgpu')}</DndHeaderSub>
                            <DndHeaderSub>{t('ui.formatter.inftime')}</DndHeaderSub>
                            <DndHeaderSub>{t('ui.formatter.accuracy')}</DndHeaderSub>
                            <DndHeaderSub>{t('ui.label.recall')}</DndHeaderSub>
                            <DndHeaderSub>{t('ui.label.precision')}</DndHeaderSub>
                            <DndHeaderSub position='last'>F1 Score</DndHeaderSub>
                        </DndHeaderSubArea>
                    </DndHeaderGroup>
                </DndHeaderArea>
                <DndContentArea>{cols && cols.map((col: any, i: number) => renderCol(col, i))}</DndContentArea>
            </DndContainer>
        </DndProvider>
    )
})

export default CompareTrial