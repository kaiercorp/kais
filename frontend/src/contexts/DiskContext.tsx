import { createContext, useState, ReactNode } from 'react'
import { DiskType } from 'common'

type DiskContextValueType = {
    disks: DiskType[]
}

type DiskContextType = {
    diskContextValue: DiskContextValueType
    updateDiskContextValue: (value: DiskContextValueType) => void 
}

const initialDiskContextValue = {
    disks: []
}

const DiskContext = createContext<DiskContextType>({diskContextValue: initialDiskContextValue, updateDiskContextValue: () => {}})

const DiskContextProvider = ({children}: {children: ReactNode}) => {
    const [diskContextValue, setDiskContextValue] = useState<DiskContextValueType>(initialDiskContextValue)    
    
    const updateDiskContextValue = (value: DiskContextValueType) => {
        setDiskContextValue(value)
    }
    
    return (
        <DiskContext.Provider value = {{diskContextValue, updateDiskContextValue}}>
            {children}
        </DiskContext.Provider>    
    )
}

export {DiskContext, DiskContextProvider}