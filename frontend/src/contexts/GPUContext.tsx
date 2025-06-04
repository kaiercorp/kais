import { createContext, useState, ReactNode } from 'react'
import { GPUType } from 'common'
import { logger } from 'helpers'

type GPUContextValueType = {
    total: number
    working: number
    idle: number
    train: number
    test: number
    gpus: GPUType[]
}

type GPUContextType = {
    gpuContextValue: GPUContextValueType
    updateGpuContextValue: (value: GPUContextValueType) => void
}

const initialGpuContextValue = {
    total: 0,
    working: 0,
    idle: 0,
    train: 0,
    test: 0,
    gpus: []    
}

const GPUContext = createContext<GPUContextType>({gpuContextValue: initialGpuContextValue, updateGpuContextValue: () => {}})

const GPUContextProvider = ({children}: {children: ReactNode}) => {
    const [gpuContextValue, setGpuContextValue] = useState<GPUContextValueType>(initialGpuContextValue)

    const updateGpuContextValue = (value: GPUContextValueType) => {
        try {
            let gpus = value.gpus
            setGpuContextValue({
                total: Number(gpus.length),
                idle: Number(gpus.filter((gpu: any) => gpu.state === 'idle').length),
                working: Number(gpus.filter((gpu: any) => gpu.state === 'train').length) +
                         Number(gpus.filter((gpu: any) => gpu.state === 'test').length),
                train: Number(gpus.filter((gpu: any) => gpu.state === 'train').length),
                test: Number(gpus.filter((gpu: any) => gpu.state === 'test').length),
                gpus: Array.isArray(gpus) ? 
                        gpus.sort((a: GPUType, b: GPUType) => { return a.id < b.id ? -1 : a.id > b.id ? 1 : 0 }) : []
            })
        } catch (e) {
            logger.error(e)
        }
    }
    
    return (
        <GPUContext.Provider value={{gpuContextValue ,updateGpuContextValue}}>
            {children}
        </GPUContext.Provider>
    )
}

export { GPUContext, GPUContextProvider } 