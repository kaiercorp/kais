import { createContext, ReactNode, useState } from 'react'
import { LayoutTypes, LayoutColor, LayoutWidth, SideBarTheme, SideBarWidth } from 'appConstants'

type LayoutContextValueType = {
    layoutColor: LayoutColor
    layoutType: LayoutTypes
    layoutWidth: LayoutWidth
    leftSideBarTheme: SideBarTheme
    leftSideBarType: SideBarWidth
}

type LayoutContextType = {
    layoutContextValue: LayoutContextValueType
    updateLayoutContextValue: (value: LayoutContextValueType) => void
}

const initialLayoutContextValue = {
    layoutColor: LayoutColor.LAYOUT_COLOR_DARK,
    layoutType: LayoutTypes.LAYOUT_ONEPAGE,
    layoutWidth: LayoutWidth.LAYOUT_WIDTH_FLUID,
    leftSideBarTheme: SideBarTheme.LEFT_SIDEBAR_THEME_DARK,
    leftSideBarType: SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED,
}

export const LayoutContext = createContext<LayoutContextType>({layoutContextValue: initialLayoutContextValue, updateLayoutContextValue: () => {}})

export const LayoutContextProvider = ({children}: {children: ReactNode}) => {
    const [layoutContextValue, setLayoutContextValue] = useState<LayoutContextValueType>(initialLayoutContextValue)

    const updateLayoutContextValue = (value: LayoutContextValueType) => {
       setLayoutContextValue(value)
    }

    return (
        <LayoutContext.Provider value={{layoutContextValue, updateLayoutContextValue}}>
            {children}
        </LayoutContext.Provider>
    )
}