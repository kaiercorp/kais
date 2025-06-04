import { useState, createContext, ReactNode } from 'react'

type LoginContextValueType = {
    isLoggedIn: boolean
}

type LoginContextType = {
    loginContextValue: LoginContextValueType
    updateLoginContextValue: (value: LoginContextValueType) => void 
} 

const LoginContext = createContext<LoginContextType>({loginContextValue: {isLoggedIn: false}, updateLoginContextValue: () => {}})

const LoginContextProvider = ({children}: {children: ReactNode}) => {
    const [loginContextValue, setLoginContextValue] = useState<LoginContextValueType>({isLoggedIn: false})    
    
    const updateLoginContextValue = (value: LoginContextValueType) => {
        setLoginContextValue(value)
    }

    return (
        <LoginContext.Provider value={{loginContextValue, updateLoginContextValue}}>
            {children}
        </LoginContext.Provider>
    )
}

export { LoginContextProvider, LoginContext }