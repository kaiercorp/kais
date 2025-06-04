import React, { Suspense, useEffect, useContext } from 'react'
import { Container } from 'react-bootstrap'
import { Outlet } from 'react-router-dom'

import { useToggle } from 'hooks'
import { changeBodyAttribute } from 'helpers'
import * as layoutConstants from 'appConstants'
import { LayoutContext } from 'contexts'

const Topbar = React.lazy(() => import('./Topbar'))
const LeftSidebar = React.lazy(() => import('./LeftSidebar'))

const loading = () => <div className=''></div>

const OnePageLayout = () => {
  const [isMenuOpened, toggleMenu] = useToggle()
  const { layoutContextValue } = useContext(LayoutContext)

  const { layoutColor, leftSideBarTheme, leftSideBarType, layoutWidth } = layoutContextValue
  useEffect(() => {
    changeBodyAttribute('data-layout', layoutConstants.LayoutTypes.LAYOUT_ONEPAGE)
  }, [])

  useEffect(() => {
    changeBodyAttribute('data-layout-color', layoutColor)
  }, [layoutColor])

  useEffect(() => {
    changeBodyAttribute('data-layout-mode', layoutWidth)
  }, [layoutWidth])

  useEffect(() => {
    changeBodyAttribute('data-leftbar-theme', leftSideBarTheme)
  }, [leftSideBarTheme])

  useEffect(() => {
    changeBodyAttribute('data-leftbar-compact-mode', leftSideBarType)
  }, [leftSideBarType])

  const openMenu = () => {
    toggleMenu()

    if (document.body) {
      if (isMenuOpened) {
        document.body.classList.remove('sidebar-enable')
      } else {
        document.body.classList.add('sidebar-enable')
      }
    }
  }

  const isCondensed = leftSideBarType === layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_CONDENSED
  const isLight = leftSideBarTheme === layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_LIGHT

  return (
    <>
      <div className='wrapper'>
        <Suspense fallback={loading()}>
          <LeftSidebar isCondensed={isCondensed} isLight={isLight} />
        </Suspense>
        <div className='content-page'>
          <div className='content'>
            <Suspense fallback={loading()}>
                <Topbar openLeftMenuCallBack={openMenu} hideLogo={true} />
            </Suspense>
            <Container fluid>
              <Outlet />
            </Container>
          </div>
        </div>
      </div>
    </>
  )
}

export default OnePageLayout
