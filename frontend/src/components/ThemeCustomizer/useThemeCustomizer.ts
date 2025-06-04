import { useState, useEffect, useCallback, useContext } from 'react'
import * as layoutConstants from 'appConstants'
import { LayoutContext } from 'contexts'

export default function useThemeCustomizer() {
    const { layoutContextValue, updateLayoutContextValue } = useContext(LayoutContext)

    const { layoutColor, layoutType, layoutWidth, leftSideBarType, leftSideBarTheme } = layoutContextValue

    const [disableLayoutWidth, setDisableLayoutWidth] = useState<boolean>(false)
    const [disableSidebarTheme, setDisableSidebarTheme] = useState<boolean>(false)
    const [disableSidebarType, setDisableSidebarType] = useState<boolean>(false)

    /**
     * change state based on props changes
     */
    const _loadStateFromProps = useCallback(() => {
        setDisableLayoutWidth(
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_DETACHED &&
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_FULL
        )

        setDisableSidebarTheme(
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL &&
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_DETACHED
        )
        setDisableSidebarType(layoutType !== layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL)
    }, [layoutType])

    useEffect(() => {
        _loadStateFromProps()
    }, [_loadStateFromProps])

    /**
     * On layout change
     */
    const changeLayoutType = (layout: string) => {
        let layoutType = layoutConstants.LayoutTypes.LAYOUT_VERTICAL
        switch (layout) {
            case 'topnav':
                layoutType = layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL
                break
            case 'detached':
                layoutType = layoutConstants.LayoutTypes.LAYOUT_DETACHED
                break
            case 'full':
                layoutType = layoutConstants.LayoutTypes.LAYOUT_FULL
                break
        }

        updateLayoutContextValue({...layoutContextValue, layoutType })
    }

    /**
     * Change the layout color
     */
    const changeLayoutColorScheme = (mode: string) => {
        let layoutColor = layoutConstants.LayoutColor.LAYOUT_COLOR_LIGHT
        switch (mode) {
            case 'dark':
                layoutColor = layoutConstants.LayoutColor.LAYOUT_COLOR_DARK
                break
        }

        updateLayoutContextValue({...layoutContextValue, layoutColor})
    }

    /**
     * Change the width mode
     */
    const changeWidthMode = (mode: string) => {
        let layoutWidth = layoutConstants.LayoutWidth.LAYOUT_WIDTH_FLUID
        let leftSideBarType = layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED
        switch (mode) {
            case 'boxed':
                layoutWidth = layoutConstants.LayoutWidth.LAYOUT_WIDTH_BOXED
                leftSideBarType = layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_CONDENSED
                break
        }

        updateLayoutContextValue({...layoutContextValue, layoutWidth, leftSideBarType }) 
    }

    /**
     * Changes the theme
     */
    const changeLeftSidebarTheme = (theme: string) => {
        let leftSideBarTheme = layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_DARK
        switch (theme) {
            case 'default':
                leftSideBarTheme = layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_DEFAULT
                break
            case 'light':
                leftSideBarTheme = layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_LIGHT
                break
        }
        
        updateLayoutContextValue({...layoutContextValue, leftSideBarTheme})
    }

    /**
     * Change the leftsiderbar type
     */
    const changeLeftSiderbarType = (type: string) => {
        let leftSideBarType = layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED
        switch (type) {
            case 'condensed':
                leftSideBarType = layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_CONDENSED
                break
            case 'scrollable':
                leftSideBarType = layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_SCROLLABLE
                break
        }
        
        updateLayoutContextValue({...layoutContextValue, leftSideBarType})
    }

    /**
     * Reset everything
     */
    const reset = () => {
        changeLayoutType(layoutConstants.LayoutTypes.LAYOUT_VERTICAL)
        changeLayoutColorScheme(layoutConstants.LayoutColor.LAYOUT_COLOR_LIGHT)
        changeWidthMode(layoutConstants.LayoutWidth.LAYOUT_WIDTH_FLUID)
        changeLeftSidebarTheme(layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_DEFAULT)
        changeLeftSiderbarType(layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED)
    }

    return {
        layoutColor,
        layoutType,
        layoutWidth,
        leftSideBarType,
        leftSideBarTheme,
        disableLayoutWidth,
        disableSidebarTheme,
        disableSidebarType,
        changeLayoutType,
        changeLayoutColorScheme,
        changeWidthMode,
        changeLeftSidebarTheme,
        changeLeftSiderbarType,
        reset,
    }
}
