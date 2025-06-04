export type MenuItemType = {
    key: string
    label: string
    isTitle?: boolean
    icon?: string
    url?: string
    badge?: {
        variant: string
        text: string
    }
    parentKey?: string
    target?: string
    children?: MenuItemType[]
    isUse?: boolean
}

const MENU_ITEMS: MenuItemType[] = [
    {
        key: 'home',
        label: 'Home',
        isTitle: false,
        icon: 'mdi mdi-home',
        url: '/dashboard'
    },
    {
        key: 'vision',
        label: 'Vision data',
        icon: 'mdi mdi-television-guide',
        children: [
            {
                key: 'vision-cls',
                label: 'Classification',
                icon: 'mdi mdi-book-settings',
                url: '/vision/cls',
                parentKey: 'vision',
            }
        ]
    },
    {
        key: 'table',
        label: 'Table data',
        icon: 'mdi mdi-table',
        children: [
            {
                key: 'table-cls',
                label: 'Classification',
                icon: 'mdi mdi-book-settings',
                url: '/table/cls',
                parentKey: 'table',
            }
        ]
    },
    {
        key: 'config',
        label: 'Configuration',
        icon: 'mdi mdi-cog',
        url: '/configuration'
    }
]

export { MENU_ITEMS }


export type ProfileOption = {
    label: string;
    icon: string;
    redirectTo: string;
}

const profileMenus: ProfileOption[] = [
    // {
    //     label: 'My Account',
    //     icon: 'mdi mdi-account-circle',
    //     redirectTo: '#',
    // },
    {
        label: 'Settings',
        icon: 'mdi mdi-account-edit',
        redirectTo: '/config/setting',
    },
    {
        label: 'Logout',
        icon: 'mdi mdi-logout',
        redirectTo: '/auth/logout',
    },
]

export { profileMenus }