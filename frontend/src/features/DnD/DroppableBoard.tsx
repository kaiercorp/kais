import {
    Droppable,
    DroppableProvided,
    DroppableStateSnapshot
} from 'react-beautiful-dnd'
import { Col, Form, Row } from 'react-bootstrap'
import styled from 'styled-components'

import Task from './Task'

interface ITaskList {
    isDraggingOver?: boolean
}

const TaskList = styled.div<ITaskList>`
padding: 2px;
height: 200px;
transition: background-color 0.2s ease;
${(props) => props.isDraggingOver ? `background-color: #464f5b` : '#464f5b'};
`;

const StyledRow = styled(Row)`
height: 200px;
`

const Container  = styled(Col)`
border: 1px solid #4a525d;
padding: 0;
margin: 0 12px;
overflow: auto;
`

const DroppableBoard = ({ 
    label, tasks, droppableId, draggingTaskId, selectedTaskIds,
    toggleSelection, toggleSelectionInGroup, multiSelectTo
}: any) => {
    
    const getSelectedMap = (selectedTaskIds: string[]) =>
        selectedTaskIds.reduce(
            (previous: any, current: any) => {
                previous[current] = true
                return previous
            }, {}
        )

    return (
        <Form.Group className={'mb-1'}>
            <StyledRow>
                <Form.Label column='sm' sm={4}>
                    {label}
                </Form.Label>
                <Container>
                    <Droppable droppableId={droppableId}>
                        {(provided: DroppableProvided, snapshot: DroppableStateSnapshot) => (
                            <TaskList
                                ref={provided.innerRef}
                                isDraggingOver={snapshot.isDraggingOver}
                                {...provided.droppableProps}
                            >
                                {tasks&&tasks.map((task: any, index: any) => {
                                    let isSelected: boolean = false
                                    if (selectedTaskIds) {
                                        isSelected = Boolean(getSelectedMap(selectedTaskIds)[task.id])
                                    }

                                    const isGhosting: boolean = isSelected && Boolean(draggingTaskId) && draggingTaskId !== task.id
                                    return (
                                        <Task
                                            key={task.id}
                                            task={task}
                                            index={index}
                                            isSelected={isSelected}
                                            isGhosting={isGhosting}
                                            toggleSelection={toggleSelection}
                                            toggleSelectionInGroup={toggleSelectionInGroup}
                                            multiSelectTo={multiSelectTo}
                                        />
                                    )
                                })}
                                {provided.placeholder}
                            </TaskList>
                        )}
                    </Droppable>
                </Container>
            </StyledRow>
        </Form.Group>
    )
}

export default DroppableBoard