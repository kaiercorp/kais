import { createContext, useState, ReactNode } from 'react'
import { ProjectType } from 'common'

type ProjectContextValueType = {
    selectedProject: ProjectType
}

type ProjectContextType = {
    projectContextValue: ProjectContextValueType 
    updateProjectContextValue: (value: ProjectContextValueType) => void
}

const ProjectContext = createContext<ProjectContextType>({projectContextValue: {selectedProject: {}}, updateProjectContextValue: () => {}})

const ProjectContextProvider = ({children}: {children: ReactNode}) => {
    const [projectContextValue, setProjectContextValue] = useState<ProjectContextValueType>({selectedProject: {}})

    const updateProjectContextValue = (value: ProjectContextValueType) => {
        setProjectContextValue(value)
    }

    return (
        <ProjectContext.Provider value={{ projectContextValue, updateProjectContextValue }}>
            {children}
        </ProjectContext.Provider>
    )
}

export { ProjectContextProvider, ProjectContext} 