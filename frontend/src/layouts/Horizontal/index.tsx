import React, { Suspense, useEffect, useContext } from 'react'
import { Outlet } from 'react-router-dom'
import { Container } from 'react-bootstrap'
import { useToggle } from 'hooks'
import * as layoutConstants from 'appConstants'
import { changeBodyAttribute } from 'helpers'
import { LayoutContext } from 'contexts'

// code splitting and lazy loading
// https://blog.logrocket.com/lazy-loading-components-in-react-16-6-6cea535c0b52
const Topbar = React.lazy(() => import('../Topbar/'))
const Navbar = React.lazy(() => import('./Navbar'))
const Footer = React.lazy(() => import('../Footer'))

const loading = () => <div className='text-center'></div>

const HorizontalLayout = () => {
  const [isMenuOpened, toggleMenu] = useToggle()
  const { layoutContextValue } = useContext(LayoutContext)

  /**
   * Open the menu when having mobile screen
   */
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

  /*
   * layout defaults
   */
  useEffect(() => {
    changeBodyAttribute('data-layout', layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL)
    changeBodyAttribute('data-leftbar-theme', layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_DEFAULT)
    changeBodyAttribute('data-leftbar-compact-mode', layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED)
  }, [])

  useEffect(() => {
    changeBodyAttribute('data-layout-color', layoutContextValue.layoutColor)
  }, [layoutContextValue.layoutColor])

  useEffect(() => {
    changeBodyAttribute('data-layout-mode', layoutContextValue.layoutWidth)
  }, [layoutContextValue.layoutWidth])

  return (
    <div className='wrapper'>
      <div className='content-page'>
        <div className='content'>
          <Suspense fallback={loading()}>
            <Topbar
              openLeftMenuCallBack={openMenu}
              navCssClasses='topnav-navbar topnav-navbar-dark'
              topbarDark={true}
            />
          </Suspense>

          <Suspense fallback={loading()}>
            <Navbar isMenuOpened={isMenuOpened} />
          </Suspense>

          <Container fluid>
            <Outlet />
          </Container>
        </div>

        <Suspense fallback={loading()}>
          <Footer />
        </Suspense>

      </div>
    </div>
  )
}

export default HorizontalLayout
