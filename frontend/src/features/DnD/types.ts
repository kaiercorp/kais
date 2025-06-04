export type TaskType = {
    id: string,
    content: string,
}

export type Column = {
    id: string,
    title: string,
    taskIds: string[],
}

export type ColumnMap = {
    [columnId: string]: Column,
}

export type TaskMap = {
    [taskId: string]: TaskType,
}

export type Entities = {
    columnOrder: string[],
    columns: ColumnMap,
    tasks: TaskMap,
}

export type DraggableRubric = {
    draggableId: string,
    type: string,
    source: DraggableLocation,
}

export type MovementMode = 'FLUID' | 'SNAP'

export type DragStart = {
    draggableId: string,
    type: string,
    source: DraggableLocation,
    mode: MovementMode,
}

export type Combine = {
    draggableId: string,
    droppableId: string,
}

export type DragUpdate = {
    draggableId: string,
    type: string,
    source: DraggableLocation,
    mode: MovementMode,
    // may not have any destination (drag to nowhere)
    destination: DraggableLocation,
    // populated when a draggable is dragging over another in combine mode
    combine: Combine,
};

export type DropReason = 'DROP' | 'CANCEL'

// published when a drag finishes
export type DropResult = {
    draggableId: string,
    type: string,
    source: DraggableLocation,
    mode: MovementMode,
    destination: DraggableLocation | null | undefined,
    combine: Combine | null | undefined,
    reason: DropReason,
}

export type DraggableLocation = {
    droppableId: string,
    index: number,
}

export type Args = {
    entities: Entities,
    selectedTaskIds: string[],
    source: DraggableLocation,
    destination: DraggableLocation,
}

export type ReorderResult = {
    entities: Entities,
    // a drop operations can change the order of the selected task array
    selectedTaskIds: string[],
}
