import {
    DragDropContext,
} from 'react-beautiful-dnd'
import { useTranslation } from 'react-i18next'

import { objDeepCopy } from 'helpers'
import DroppableBoard from './DroppableBoard'
import { useEffect, useState } from 'react'
import { multiSelect, mutliDragAwareReorder } from './utils'
import { DraggableLocation, DragStart, DropResult, Entities, ReorderResult, TaskMap, TaskType } from './types'


const TableCLSColumnsDnd = ({ input_columns, except_columns, onMoveItem }: any) => {
    const [t] = useTranslation('translation')

    const [draggingTaskId, setDraggingTaskId] = useState<any>(null)
    const [selectedTaskIds, setSelectedTaskIds] = useState<any>(null)
    const [entities, setEntities] = useState<any>(null)

    const [inputTasks, setInputTasks] = useState<TaskType[]>()
    const [exceptTasks, setExceptTasks] = useState<TaskType[]>()

    useEffect(() => {
        let tasks = Array<TaskType>()
        if (input_columns) {
            input_columns.forEach((col: string) => {
                tasks.push({ id: col, content: col })
            })
        }
        setInputTasks(tasks)
    }, [input_columns])

    useEffect(() => {
        let tasks = Array<TaskType>()
        if (except_columns) {
            except_columns.forEach((col: string) => {
                tasks.push({ id: col, content: col })
            })
        }
        setExceptTasks(tasks)
    }, [except_columns])

    useEffect(() => {
        let tasks: TaskType[] = Array<TaskType>()
        if (inputTasks && !exceptTasks) {
            tasks = objDeepCopy(inputTasks)
        } else if (!inputTasks && exceptTasks) {
            tasks = objDeepCopy(exceptTasks)
        } else if (inputTasks && exceptTasks) {
            tasks = inputTasks.concat(exceptTasks)
        } 

        const taskMap: TaskMap = tasks.reduce(
            (previous: TaskMap, current: TaskType): TaskMap => {
                previous[current.id] = current;
                return previous;
            },
            {},
        )

        const initialEntities: Entities = {
            columnOrder: ["droppable-inputlist", "droppable-exceptlist"],
            columns: {
                "droppable-inputlist": {
                    id: 'droppable-inputlist',
                    title: '',
                    taskIds: inputTasks?inputTasks.map((task: TaskType): string => task.id):[],
                },
                "droppable-exceptlist": {
                    id: 'droppable-exceptlist',
                    title: '',
                    taskIds: exceptTasks?exceptTasks.map((task: TaskType): string => task.id):[],
                },
            },
            tasks: taskMap,
        }
        setEntities(initialEntities)
    }, [inputTasks, exceptTasks])

    const toggleSelection = (taskId: string) => {
        const wasSelected: boolean = selectedTaskIds && selectedTaskIds.includes(taskId)

        const newTaskIds: string[] = (() => {
            // Task was not previously selected
            // now will be the only selected item
            if (!wasSelected) {
                return [taskId]
            }

            // Task was part of a selected group
            // will now become the only selected item
            if (selectedTaskIds.length > 1) {
                return [taskId]
            }

            // task was previously selected but not in a group
            // we will now clear the selection
            return []
        })()
        
        setSelectedTaskIds(newTaskIds)
    }

    const toggleSelectionInGroup = (taskId: string) => {
        const index: number = selectedTaskIds.indexOf(taskId)

        // if not selected - add it to the selected items
        if (index === -1) {
            const shallow: string[] = objDeepCopy(selectedTaskIds)
            setSelectedTaskIds([...shallow, taskId])
            return
        }

        // it was previously selected and now needs to be removed from the group
        const shallow: string[] = objDeepCopy(selectedTaskIds)
        shallow.splice(index, 1)
        setSelectedTaskIds(shallow)
    }

    // This behaviour matches the MacOSX finder selection
    const multiSelectTo = (newTaskId: string) => {
        const updated: any = multiSelect(
            entities,
            selectedTaskIds,
            newTaskId,
        )

        if (updated == null) {
            return
        }

        setSelectedTaskIds(updated)
    }

    const onDragStart = (start: DragStart) => {
        // if (!selectedTaskIds) return
        const id: string = start.draggableId;
        // const selected: any = selectedTaskIds.find(
        //     (taskId: string): boolean => taskId === id,
        // )

        // if dragging an item that is not selected - unselect all items
        if (!selectedTaskIds) {
            setSelectedTaskIds([id])
        }

        setDraggingTaskId(id)
    }

    const onDragEnd = (result: DropResult) => {
        const destination: DraggableLocation | null | undefined = result.destination;
        const source: DraggableLocation = result.source;

        // nothing to do
        if (!destination || result.reason === 'CANCEL') {
            setDraggingTaskId(null)
            return
        }
        if (destination.droppableId === source.droppableId) {
            setDraggingTaskId(null)
            return
        }

        const processed: ReorderResult = mutliDragAwareReorder({
            entities: entities,
            selectedTaskIds: objDeepCopy(selectedTaskIds),
            source,
            destination,
        })

        setEntities(processed.entities)
        setSelectedTaskIds(processed.selectedTaskIds)
        setDraggingTaskId(null)

        if (destination.droppableId === 'droppable-exceptlist') {
            let newExcept = objDeepCopy(except_columns)
            newExcept = newExcept.concat(selectedTaskIds)
            onMoveItem(newExcept)
            setSelectedTaskIds(null)
        } else if (destination.droppableId === 'droppable-inputlist') {
            let newExcept = objDeepCopy(except_columns)
            newExcept = newExcept.filter((col: string) => {
                return !selectedTaskIds.includes(col)
            })
            onMoveItem(newExcept)
            setSelectedTaskIds(null)
        }
    }

    return (
        <DragDropContext onDragStart={onDragStart} onDragEnd={onDragEnd}>
            <DroppableBoard
                label={t('ui.train.includeCol')}
                droppableId="droppable-inputlist"
                tasks={inputTasks}
                selectedTaskIds={selectedTaskIds}
                draggingTaskId={draggingTaskId}
                toggleSelection={toggleSelection}
                toggleSelectionInGroup={toggleSelectionInGroup}
                multiSelectTo={multiSelectTo}

            />
            <DroppableBoard
                label={t('ui.train.exceptCol')}
                droppableId="droppable-exceptlist"
                tasks={exceptTasks}
                selectedTaskIds={selectedTaskIds}
                draggingTaskId={draggingTaskId}
                toggleSelection={toggleSelection}
                toggleSelectionInGroup={toggleSelectionInGroup}
                multiSelectTo={multiSelectTo}

            />
        </DragDropContext>
    )
}

export default TableCLSColumnsDnd