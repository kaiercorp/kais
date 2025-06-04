import classNames from 'classnames'
import { useTranslation } from 'react-i18next'

import * as layoutConstants from 'appConstants'
import { useViewport, useSocket } from 'hooks'
import ProfileDropdown from 'layouts/Topbar/ProfileDropdown'
import { profileMenus } from 'appConstants'
import { Link } from 'react-router-dom'
import { useContext } from 'react'
import { logger } from 'helpers'
import { DiskInfo } from 'features'
import { LocationContext, ProjectContext, LayoutContext, DiskContext, GPUContext } from 'contexts'
import { APICore } from 'helpers/api/apiCore'
import { Card } from 'react-bootstrap'

interface TopbarProps {
  hideLogo?: boolean
  navCssClasses?: string
  openLeftMenuCallBack?: () => void
  topbarDark?: boolean
}

const api = new APICore()

const Topbar = ({ hideLogo, navCssClasses, openLeftMenuCallBack }: TopbarProps) => {
  const [t] = useTranslation('translation')
  const containerCssClasses = hideLogo ? '' : 'container-fluid'

  const { locationContextValue } = useContext(LocationContext)
  const { projectContextValue } = useContext(ProjectContext)
  const { diskContextValue, updateDiskContextValue } = useContext(DiskContext)
  const { updateGpuContextValue } = useContext(GPUContext)
  const { layoutContextValue, updateLayoutContextValue } = useContext(LayoutContext)
  const { leftSideBarType, layoutType } = layoutContextValue

  const user = api.getLoggedInUser()

  const { width } = useViewport()

  const handleLeftMenuCallBack = () => {
    if (openLeftMenuCallBack) openLeftMenuCallBack()

    if (width >= 768) {
      if (leftSideBarType === 'fixed' || leftSideBarType === 'scrollable') {
        updateLayoutContextValue({ ...layoutContextValue, leftSideBarType: layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_CONDENSED })
      } else if (leftSideBarType === 'condensed') {
        updateLayoutContextValue({ ...layoutContextValue, leftSideBarType: layoutConstants.SideBarWidth.LEFT_SIDEBAR_TYPE_FIXED })
      }
    }
  }

  const handleSocketMessage = (e: MessageEvent<any>) => {
    try {
      let system = JSON.parse(e.data)

      const disks = system.disk
      updateDiskContextValue({ disks: Array.isArray(disks) ? disks : [] })

      const gpus = system.gpus
      updateGpuContextValue({ gpus: Array.isArray(gpus) ? gpus : [], total: 0, working: 0, idle: 0, train: 0, test: 0 })
    } catch (e) {
      logger.error(e)
    }
  }
  useSocket('/sys', 'System', handleSocketMessage, { shouldCleanup: true, shouldConnect: true })

  return (
    <div className={classNames('navbar-custom', navCssClasses)}>
      <div className={containerCssClasses}>
        <span className={classNames('badge', 'p-1', 'bg-success')} style={{ marginRight: '10px', fontSize: '16px' }}>
          {projectContextValue.selectedProject?.project_name}
        </span>
        <span style={{ lineHeight: '60px', fontSize: '16px', fontWeight: '800' }}>
          <>{t(`title.${locationContextValue.location}`, locationContextValue.locationValue)}</>
        </span>

        {(layoutType === layoutConstants.LayoutTypes.LAYOUT_VERTICAL ||
          layoutType === layoutConstants.LayoutTypes.LAYOUT_FULL ||
          layoutType === layoutConstants.LayoutTypes.LAYOUT_ONEPAGE) && (
            <button className='button-menu-mobile open-left' onClick={handleLeftMenuCallBack}>
              <i className='mdi mdi-menu' />
            </button>
          )}

        <ul className='list-unstyled topbar-menu float-end mb-0'>
          {
            user ?
              <ProfileDropdown
                menuItems={profileMenus}
                username={user.username}
                userTitle={user.department}
              />
              : <div style={{ marginTop: '15px' }}><Link to='/auth/login'>Login</Link></div>
          }
        </ul>

        {
          diskContextValue.disks && diskContextValue.disks.map((disk: any) => {
            if (disk.path === "ROOT_PATH") {
              return null
            }
            return (<ul className='list-unstyled topbar-menu float-end mb-0' key={disk.path}>
              <DiskInfo diskInfo={disk} />
            </ul>)
          })
        }

        {
          (!diskContextValue.disks || diskContextValue.disks.length < 1) &&
          <ul className='list-unstyled topbar-menu float-end mb-0'>
            <Card><Card.Body>NO Data Root</Card.Body></Card>
          </ul>
        }
      </div>
    </div>
  )
}

export default Topbar
