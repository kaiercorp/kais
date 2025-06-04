import { createContext, useState, ReactNode } from 'react'
import { FilterType, initialFilter } from 'common'

type FilterContextValueType = {
    filter: FilterType
    useFilter: boolean 
}

type FilterContextType = {
    filterContextValue: FilterContextValueType
    updateFilterContextValue: (value: Partial<FilterContextValueType>) => void
}

const initialFilterContextValue = {
    filter: initialFilter,
    useFilter: false
}

const FilterContext = createContext<FilterContextType>({filterContextValue: initialFilterContextValue, updateFilterContextValue: () => {}}) 

const FilterContextProvider = ({children} : {children: ReactNode}) => {
    const [filterContextValue, setFilterContextValue] = useState<FilterContextValueType>(initialFilterContextValue)
    
    const updateFilterContextValue = (value: Partial<FilterContextValueType>) => {
        setFilterContextValue((filterContextValue) => ({...filterContextValue, ...value}))
    }

    return (
        <FilterContext.Provider value={{filterContextValue, updateFilterContextValue}}>
            {children}
        </FilterContext.Provider>
    )
}

export {FilterContextProvider, FilterContext}