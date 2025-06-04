import { MenuItemType } from 'appConstants'
import { useQueryClient } from '@tanstack/react-query'
import { ConfigType } from 'common' 
import { QUERY_KEY } from 'helpers/api'

const GetMenuItems = () => {
    const queryClient = useQueryClient()
    const configurationsQueryData = queryClient.getQueryData<{configs: ConfigType[], version: string}>([QUERY_KEY.configurations])
    
    let menuConfig = configurationsQueryData?.configs?.filter((config: any) => {
        if (config.config_key === 'MENU') return true
        return false
    })[0]

    if (!menuConfig || !menuConfig.config_val) return []

    return filterItems(JSON.parse(menuConfig.config_val)) 
}

const findAllParent = (menuItems: MenuItemType[], menuItem: MenuItemType): string[] => {
    let parents: string[] = []
    const parent = findMenuItem(menuItems, menuItem['parentKey'])

    if (parent) {
        parents.push(parent['key'])
        if (parent['parentKey']) parents = [...parents, ...findAllParent(menuItems, parent)]
    }

    return parents
}

const findMenuItem = (
    menuItems: MenuItemType[] | undefined,
    menuItemKey: MenuItemType['key'] | undefined
): MenuItemType | null => {
    if (menuItems && menuItemKey) {
        for (var i = 0; i < menuItems.length; i++) {
            if (menuItems[i].key === menuItemKey) {
                return menuItems[i]
            }
            var found = findMenuItem(menuItems[i].children, menuItemKey)
            if (found) return found
        }
    }
    return null
}

const filterItems = (menus:any[]) => {
    return menus.filter((menu: any) => {
        if (menu.children) {
            menu.children = filterItems(menu.children)
        }
        return menu.isUse
    })
}

export { GetMenuItems, findAllParent, findMenuItem }
