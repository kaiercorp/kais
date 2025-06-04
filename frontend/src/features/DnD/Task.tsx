import { KEYCODE } from "appConstants";
import { Draggable, DraggableProvided, DraggableStateSnapshot } from "react-beautiful-dnd"
import styled from "styled-components"

const getBackgroundColor = ({ isSelected, isGhosting }: any): string => {
    if (isGhosting) {
        return '#464f5b'
    }

    if (isSelected) {
        return '#464f5b'
    }

    return '#404954'
}

const getColor = ({ isSelected, isGhosting }: any): string => {
    if (isGhosting) {
        return '#ffffff'
    }
    if (isSelected) {
        return '#ffffff'
    }
    return '#e3eae'
}

const getBorder = ({ isSelected, isGhosting }: any): string => {
    if (isGhosting) {
        return '1px solid rgba(0, 0, 0, 0.125)'
    }
    if (isSelected) {
        return '1px solid #ffffff'
    }
    return '1px solid rgba(0, 0, 0, 0.125)'
}

interface IContainer {
    isDragging?: boolean
    isGhosting?: boolean
    isSelected?: boolean
}

const Container = styled.div<IContainer>`
    font-size: 12px;
    background-color: ${(props) => getBackgroundColor(props)};
    color: ${(props) => getColor(props)};
    padding: 2px;
    margin-bottom: 1px;
    border-radius: 3px;
    border: ${(props => getBorder(props))};
    ${(props) =>
        props.isDragging ? `box-shadow: 2px 2px 1px #464f5b;` : ''} ${(
            props,
        ) =>
            props.isGhosting
                ? 'opacity: 0.8;'
                : ''}
  
    /* needed for SelectionCount */
    position: relative;
  
    /* avoid default outline which looks lame with the position: absolute; */
    &:focus {
      outline: none;
      border-color: 'grey';
    }
  `;

const primaryButton = 0

const Task = ({
    task, index,
    isSelected, isGhosting,
    toggleSelection, toggleSelectionInGroup,
    multiSelectTo
}: any) => {

    const onKeyDown = (
        event: any,
        provided: DraggableProvided,
        snapshot: DraggableStateSnapshot,
    ) => {
        if (event.defaultPrevented) {
            return
        }

        if (snapshot.isDragging) {
            return
        }

        if (event.keyCode !== KEYCODE.enter) {
            return
        }

        // we are using the event for selection
        event.preventDefault()

        performAction(event)
    };

    // Using onClick as it will be correctly
    // preventing if there was a drag
    const onClick = (event: any) => {
        if (event.defaultPrevented) {
            return
        }

        if (event.button !== primaryButton) {
            return
        }

        // marking the event as used
        event.preventDefault()
        performAction(event)
    }

    const onTouchEnd = (event: any) => {
        if (event.defaultPrevented) {
            return
        }

        // marking the event as used
        // we would also need to add some extra logic to prevent the click
        // if this element was an anchor
        event.preventDefault()
        toggleSelectionInGroup(task.id)
    }

    // Determines if the platform specific toggle selection in group key was used
    const wasToggleInSelectionGroupKeyUsed = (event: any) => {
        const isUsingWindows = navigator.platform.indexOf('Win') >= 0
        return isUsingWindows ? event.ctrlKey : event.metaKey
    }

    // Determines if the multiSelect key was used
    const wasMultiSelectKeyUsed = (event: any) => event.shiftKey

    const performAction = (event: any) => {
        if (wasToggleInSelectionGroupKeyUsed(event)) {
            toggleSelectionInGroup(task.id)
            return
        }

        if (wasMultiSelectKeyUsed(event)) {
            multiSelectTo(task.id)
            return
        }

        toggleSelection(task.id)
    }

    return (
        <Draggable draggableId={`${task.id}`} index={index}>
            {(provided: DraggableProvided, snapshot: DraggableStateSnapshot) => {
                return (
                    <div>
                        <Container
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            {...provided.dragHandleProps}
                            isDragging={snapshot.isDragging}
                            isGhosting={isGhosting}
                            isSelected={isSelected}
                            onClick={onClick}
                            onTouchEnd={onTouchEnd}
                            onKeyDown={(event: any) => onKeyDown(event, provided, snapshot)}
                        >
                            {task.content}
                        </Container>
                    </div>
                )
            }}
        </Draggable>
    )
}


export default Task