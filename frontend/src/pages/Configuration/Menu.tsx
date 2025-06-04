import { useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import { logger, objDeepCopy } from 'helpers'
import { Button, Card, Table } from 'react-bootstrap'
import { LocationContext } from 'contexts'

import { ApiFetchMenu, ApiUpdateMenu } from 'helpers/api'
import { MenuType } from 'common'
import { DefaultSwitch } from 'components'
import { ButtonArea, CardHeaderLeft } from '../../components/Containers'

const MenuLanding = () => {
    const [t] = useTranslation('translation')
    const [menus, setMenus] = useState<any[]>([])

    const { updateLocationContextValue } = useContext(LocationContext)

    const { menu } = ApiFetchMenu()
    const updateMenu = ApiUpdateMenu()

    useEffect(() => {
        logger.log(`Change Location to ${t('title.menu')}`)
        updateLocationContextValue({ location: 'menu' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        setMenus(menu)
    }, [menu])

    const handleChangeMenu = (key: string) => {
        setMenus(objDeepCopy(menus).map((menu: MenuType) => {
            if (menu.key === key) {
                return menu.key === key ? {
                    ...menu,
                    isUse: !menu.isUse,
                    children: menu.children ? menu.children.map((child: MenuType) => {
                        return {
                            ...child,
                            isUse: !menu.isUse,
                        }
                    }) : []
                } : menu
            }

            return {
                ...menu,
                children: menu.children ? menu.children.map((child: MenuType) => {
                    return child.key === key ? {
                        ...child,
                        isUse: !child.isUse,
                    } : child
                }) : []
            }
        }))
    }

    const handleSaveBtn = () => {
        updateMenu.mutate(menus)
    }

    const handleCancleBtn = () => {
        setMenus(menu)
    }

    return (
        <>
            <Card>
                <Card.Header>
                    <CardHeaderLeft>MENUS</CardHeaderLeft>
                </Card.Header>
                <Card.Body>
                    <Table className="mb-0" bordered>
                        <thead>
                            <tr>
                                <th colSpan={2}>#</th>
                                <th>Use</th>
                                <th>Title</th>
                                <th>Group</th>
                                <th>key</th>
                                <th>parent</th>
                                <th>icon</th>
                            </tr>
                        </thead>
                        {menus.map((menu: any, index: any) => {
                            return (
                                <tbody key={`menu-${menu.key}`}>
                                    <tr >
                                        <th scope="row">{index}</th>
                                        <td></td>
                                        <td>
                                            <DefaultSwitch
                                                disabled={menu.key === 'configuration'}
                                                label={""}
                                                checked={menu.isUse}
                                                onChange={() => handleChangeMenu(menu.key)}
                                            />
                                        </td>
                                        <td>{menu.label}</td>
                                        <td>{menu.group}</td>
                                        <td>{menu.key}</td>
                                        <td>{menu.parentKey}</td>
                                        <td>{menu.icon}</td>
                                    </tr>
                                    {menu.children && menu.children.map((mc: any, jndex: any) => {
                                        return (
                                            <tr key={`menu-sub-${mc.key}`}>
                                                <th></th>
                                                <td>{jndex}</td>
                                                <td>
                                                    <DefaultSwitch
                                                        disabled={mc.key === 'menus' && mc.parentKey === 'configuration'}
                                                        label={""}
                                                        checked={mc.isUse}
                                                        onChange={() => handleChangeMenu(mc.key)}
                                                    />
                                                </td>
                                                <td>{mc.label}</td>
                                                <td>{mc.group}</td>
                                                <td>{mc.key}</td>
                                                <td>{mc.parentKey}</td>
                                                <td>{mc.icon}</td>
                                            </tr>
                                        )
                                    })}
                                </tbody>
                            )
                        })}
                    </Table>

                </Card.Body>
            </Card>
            <ButtonArea>
                <Button variant={'danger'} onClick={handleCancleBtn}>{t('button.cancel')}</Button>
                <Button onClick={handleSaveBtn}>{t('button.save')}</Button>
            </ButtonArea>
        </>
    )
}

export default MenuLanding 