import { useEffect, useRef, useState } from 'react'
import SimpleBar from 'simplebar-react'
import { Link } from 'react-router-dom'
import styled from 'styled-components'

import AppMenu from 'layouts/Menu'

import logo from 'assets/images/kaier_logo_big.png'
import logoDark from 'assets/images/kaier_logo_big_white.png'
import { useQueryClient } from '@tanstack/react-query'
import { ConfigType } from 'common'
import { QUERY_KEY, ApiFetchMenu, objDeepCopy } from 'helpers'

const Logo = styled.div`
padding: 0px 20px 0px 20px;
width: 260px;
display: table-cell;
vertical-align: middle;

color: #ffffff;
font-weight: 800;
font-size: 30px;

& .subtitle {
  margin-left: 5px;
  font-weight: 300;
  font-size: 14px;
}
`

const LogoArea = styled.div`
text-align: center;
bottom: 40px;
left: 90px;
position: fixed;

& img {
  height: 20px;
}
`

const filterItems = (menus: any[]) => {
  return menus.filter((menu: any) => {
    if (menu.children) {
      menu.children = filterItems(menu.children)
    }
    return menu.isUse
  })
}

interface LeftSidebarProps {
  hideLogo?: boolean
  isLight: boolean
  isCondensed: boolean
}

const LeftSidebar = ({ isCondensed, isLight, hideLogo }: LeftSidebarProps) => {
  const menuNodeRef = useRef<HTMLDivElement>(null)
  const [menus, setMenus] = useState<any[]>([])

  const queryClient = useQueryClient()
  const configurationsQueryData = queryClient.getQueryData<{configs: ConfigType[], version: string}>([QUERY_KEY.configurations])
  const { isLoading: isMenuLoading, error: menuError, menu} = ApiFetchMenu()

  const handleOtherClick = (e: MouseEvent) => {
    if (menuNodeRef && menuNodeRef.current && menuNodeRef.current.contains(e.target as Node)) return

    if (document.body) {
      document.body.classList.remove('sidebar-enable')
    }
  }

  useEffect(() => {
    document.addEventListener('mousedown', handleOtherClick, false)

    return () => {
      document.removeEventListener('mousedown', handleOtherClick, false)
    }

  }, [])

  useEffect(() => {
    if (isMenuLoading || menuError) return

    if (menu.length === 0) return

    setMenus(filterItems(objDeepCopy(menu)))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [menu])

  return (
    <div className='leftside-menu' ref={menuNodeRef}>
      {!hideLogo && (
        <>
          <Link to='/' className='logo text-center logo-light'>
            <span className='logo-lg'>
              <Logo>
                KAI.S
                <span className='subtitle'>{configurationsQueryData?.version}</span>
              </Logo>
            </span>
            <span className='logo-sm'>
              <Logo>K</Logo>
            </span>
          </Link>

          <Link to='/' className='logo text-center logo-dark'>
            <span className='logo-lg'>
              <Logo>
                KAI.S
                <span className='subtitle'>{configurationsQueryData?.version}</span>
              </Logo>
            </span>
            <span className='logo-sm'>
              <Logo>K</Logo>
            </span>
          </Link>
        </>
      )}

      {!isCondensed && (
        <SimpleBar style={{ maxHeight: '100%' }} timeout={500} scrollbarMaxSize={320}>
          <AppMenu menuItems={menus} />
        </SimpleBar>
      )}
      {isCondensed && <AppMenu menuItems={menus} />}

      {!isCondensed && (
        <LogoArea>
          <span className='logo-lg'>
            <img src={isLight ? logo : logoDark} alt='logo' />
          </span>
        </LogoArea>
      )}
    </div>
  )
}

export default LeftSidebar
