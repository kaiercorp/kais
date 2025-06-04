import { useContext } from 'react'
import { Link } from 'react-router-dom'
import classNames from 'classnames'
import * as layoutConstants from 'appConstants'
import { useToggle, useViewport } from 'hooks'
import { notifications, searchOptions } from './data'
import LanguageDropdown from './LanguageDropdown'
import NotificationDropdown from './NotificationDropdown'
import SearchDropdown from './SearchDropdown'
import TopbarSearch from './TopbarSearch'
import AppsDropdown from './AppsDropdown'
import ProfileDropdown from './ProfileDropdown'
import { profileMenus } from 'appConstants'
import { LayoutContext } from 'contexts'

type TopbarProps = {
  hideLogo?: boolean
  navCssClasses?: string
  openLeftMenuCallBack?: () => void
  topbarDark?: boolean
}

const Topbar = ({ hideLogo, navCssClasses, openLeftMenuCallBack, topbarDark }: TopbarProps) => {
  const { width } = useViewport()
  const [isMenuOpened, toggleMenu] = useToggle()
  const { layoutContextValue, updateLayoutContextValue } = useContext(LayoutContext)
  const { layoutType, leftSideBarType } = layoutContextValue

  const containerCssClasses = !hideLogo ? 'container-fluid' : ''

  /**
   * Toggle the leftmenu when having mobile screen
   */
  const handleLeftMenuCallBack = () => {
    toggleMenu()
    if (openLeftMenuCallBack) openLeftMenuCallBack()

    switch (layoutType) {
      case layoutConstants.LayoutTypes.LAYOUT_VERTICAL:
        if (width >= 768) {
          if (leftSideBarType === 'fixed' || leftSideBarType === 'scrollable')
            updateLayoutContextValue({...layoutContextValue, leftSideBarType: layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_CONDENSED})
          if (leftSideBarType === 'condensed')
            updateLayoutContextValue({...layoutContextValue, leftSideBarType: layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED})
        }
        break

      case layoutConstants.LayoutTypes.LAYOUT_FULL:
        if (document.body) {
          document.body.classList.toggle('hide-menu')
        }
        break
      default:
        break
    }
  }

  return (
    <div className={classNames('navbar-custom', navCssClasses)}>
      <div className={containerCssClasses}>

        <ul className='list-unstyled topbar-menu float-end mb-0'>
          <li className='notification-list topbar-dropdown d-xl-none'>
            <SearchDropdown />
          </li>
          <li className='dropdown notification-list topbar-dropdown'>
            <LanguageDropdown />
          </li>
          <li className='dropdown notification-list'>
            <NotificationDropdown notifications={notifications} />
          </li>
          <li className='dropdown notification-list d-none d-sm-inline-block'>
            <AppsDropdown />
          </li>
          <li className="dropdown notification-list">
            <ProfileDropdown
              menuItems={profileMenus}
              username={'Dominic Keller'}
              userTitle={'Founder'}
            />
          </li>

        </ul>

        {/* toggle for vertical layout */}
        {(layoutType === layoutConstants.LayoutTypes.LAYOUT_VERTICAL ||
          layoutType === layoutConstants.LayoutTypes.LAYOUT_FULL) && (
            <button className='button-menu-mobile open-left' onClick={handleLeftMenuCallBack}>
              <i className='mdi mdi-menu' />
            </button>
          )}

        {/* toggle for horizontal layout */}
        {layoutType === layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL && (
          <Link to='#' className={classNames('navbar-toggle', { open: isMenuOpened })} onClick={handleLeftMenuCallBack}>
            <div className='lines'>
              <span></span>
              <span></span>
              <span></span>
            </div>
          </Link>
        )}

        {/* toggle for detached layout */}
        {layoutType === layoutConstants.LayoutTypes.LAYOUT_DETACHED && (
          <Link to='#' className='button-menu-mobile disable-btn' onClick={handleLeftMenuCallBack}>
            <div className='lines'>
              <span></span>
              <span></span>
              <span></span>
            </div>
          </Link>
        )}
        <TopbarSearch options={searchOptions} />
      </div>
    </div>
  )
}

export default Topbar
