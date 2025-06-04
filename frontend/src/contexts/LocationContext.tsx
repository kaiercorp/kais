import { createContext, useState, ReactNode } from 'react'

type LocationContextValueType = {
    location: string
    locationValue?: {
        trial: string 
    }
}

type LocationContextType = {
    locationContextValue: LocationContextValueType
    updateLocationContextValue: (value: LocationContextValueType) => void 
}

export const LocationContext = createContext<LocationContextType>({locationContextValue: {location: ''}, updateLocationContextValue: () => {}}) 

export const LocationContextProvider = ({ children }: {children: ReactNode}) => {
    const [locationContextValue, setLocationContextValue] = useState<LocationContextValueType>({location: ''});

    const updateLocationContextValue = (value: LocationContextValueType) => {
        setLocationContextValue(value)
    }

    return (
        <LocationContext.Provider value={{ locationContextValue, updateLocationContextValue}}>
            {children}
        </LocationContext.Provider>
    )
}