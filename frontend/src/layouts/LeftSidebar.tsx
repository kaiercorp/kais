import React, { useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import SimpleBar from 'simplebar-react'
import classNames from 'classnames'
import { GetMenuItems } from 'helpers'
import AppMenu from './Menu/'

type SideBarContentProps = {
  hideUserProfile: boolean
}

const SideBarContent = ({ hideUserProfile }: SideBarContentProps) => {
  return (
    <>
      {!hideUserProfile && (
        <div className='leftbar-user'>
          <Link to='/'>
            <span className='leftbar-user-name'>Dominic Keller</span>
          </Link>
        </div>
      )}
      <AppMenu menuItems={GetMenuItems()} />

      <div
        className={classNames('help-box', 'text-center', {
          'text-white': hideUserProfile,
        })}
      >
        <Link to='/' className='float-end close-btn text-white'>
          <i className='mdi mdi-close' />
        </Link>

        <h5 className='mt-3'>Unlimited Access</h5>
        <p className='mb-1'>Upgrade to plan to get access to unlimited reports</p>
        <button className={classNames('btn', 'btn-sm', hideUserProfile ? 'btn-outline-light' : 'btn-outline-primary')}>
          Upgrade
        </button>
      </div>
      <div className='clearfix' />
    </>
  )
}

type LeftSidebarProps = {
  hideLogo?: boolean
  hideUserProfile: boolean
  isLight: boolean
  isCondensed: boolean
}

const LeftSidebar = ({ isCondensed, isLight, hideLogo, hideUserProfile }: LeftSidebarProps) => {
  const menuNodeRef = useRef<HTMLDivElement>(null)

  /**
   * Handle the click anywhere in doc
   */
  const handleOtherClick = (e: MouseEvent) => {
    if (menuNodeRef && menuNodeRef.current && menuNodeRef.current.contains(e.target as Node)) return
    // else hide the menubar
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

  return (
    <div className='leftside-menu' ref={menuNodeRef}>

      {!isCondensed && (
        <SimpleBar style={{ maxHeight: '100%' }} timeout={500} scrollbarMaxSize={320}>
          <SideBarContent hideUserProfile={hideUserProfile} />
        </SimpleBar>
      )}
      {isCondensed && <SideBarContent hideUserProfile={hideUserProfile} />}
    </div>
  )
}

export default LeftSidebar
