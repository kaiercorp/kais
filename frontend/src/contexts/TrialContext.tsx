import { ReactNode, createContext, useState } from 'react'
import { TrialType } from 'common'

type TrialContextValueType = {
//    code, error, loading
    trials: TrialType[] | undefined
    trainMode: string | undefined
    selectedRows: TrialType[] | undefined
    requestData: TrialType | undefined
}

type TrialContextType = {
    trialContextValue: TrialContextValueType 
    updateTrialContextValue: (value: Partial<TrialContextValueType>) => void
}

const initialTrialContextValue = {
    trials: undefined,
    trainMode: undefined,
    selectedRows: undefined,
    requestData: undefined
}

const TrialContext = createContext<TrialContextType>({trialContextValue: initialTrialContextValue, updateTrialContextValue: () => {}})

const TrialContextProvider = ({children} : {children: ReactNode}) => {
    const [trialContextValue, setTrialContextValue] = useState<TrialContextValueType>(initialTrialContextValue)    
    
    const updateTrialContextValue = (value: Partial<TrialContextValueType>) => {
        setTrialContextValue({...trialContextValue, ...value})
    }
    
    return (
        <TrialContext.Provider value={{trialContextValue, updateTrialContextValue}}>
            {children}
        </TrialContext.Provider>
    )
}

export {TrialContextProvider, TrialContext}