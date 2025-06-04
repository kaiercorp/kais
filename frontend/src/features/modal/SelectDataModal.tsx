import { useEffect, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query' 
import { Button, Col, Modal, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

import styled from 'styled-components'
import { IconButton } from 'components'
import { ApiFetchDirs, /*ApiFetchDatasets,*/ objDeepCopy } from 'helpers'
import { QUERY_KEY } from 'helpers'
import { DirectoryType } from 'common'

type ModalProps = {
  show: boolean
  selectData: (path: string, id: number) => void
  toggle: () => void
  isTest: boolean,
  dataType: string
  directoryId: number | undefined
  onDirectoryIdChange: (directoryId: number) => void 
}

type RenderProps = {
  dirs: DirectoryType[] 
  parent: string
  depth: number
  openSubDir: (parent_id: number) => void
  closeSubDir: (dirs: DirectoryType[], path: string) => void
  selectDir: (path: string, id: number) => void
  isTest: boolean
}

interface BlankProps {
  depth: number
}

const Blank = styled.div<BlankProps>`
  width: ${(props) => 26 * props.depth + 'px'};
  height: 10px;
`

const renderTree = ({ dirs, parent, depth, openSubDir, closeSubDir, selectDir, isTest }: RenderProps) => {
  if (dirs === null) return null
  const parent_path = parent === '' ? '' : parent + '/'
  
  return dirs.map((dir: any) => {
    return (
      <Row key={`dir-tree-${parent}/${dir.name}`} style={{padding: '10px 15px'}}>
        <Col xxl={12}>
          <Row>
            {depth > 0 ? <Blank depth={depth} /> : null}
            {(dir.is_trainable !== true && isTest !== true)
              ?(
                <Blank depth={1} />
              )
              :(dir.is_testable !== true && isTest === true)
              ?(
                <Blank depth={1} />
              )
              :(<IconButton
                color='success'
                icon='mdi-form-select'
                onClick={() => {
                  selectDir(`${parent_path}${dir.name}`, dir.id)
                }}
              />)
            }
            {dir.is_leaf ? (
              <Blank depth={1} />
            ) : dir.is_open ? (
              <IconButton
                color='danger'
                icon='mdi-minus-thick'
                onClick={() => closeSubDir(dirs, `${parent_path}${dir.name}`)}
              />
            ) : (
              <IconButton color='info' icon='mdi-plus-thick' onClick={() => openSubDir(dir.id)} />
            )}
            <Col className='col-auto'>
              <div style={{ width: '24px', height: '24px' }}>
                <span className='avatar-title bg-light text-secondary rounded'>
                  <i className={'mdi mdi-folder-zip'}></i>
                </span>
              </div>
            </Col>
            <Col className='ps-0'>{dir.name}</Col>
          </Row>
          {dir.dirs
            ? renderTree({
                dirs: dir.dirs,
                parent: parent_path + dir.name,
                depth: depth + 1,
                openSubDir: openSubDir,
                closeSubDir: closeSubDir,
                selectDir: selectDir,
                isTest
              })
            : null}
        </Col>
      </Row>
    )
  })
}

const SelectDataModal = ({ show, selectData, toggle, isTest, dataType, directoryId, onDirectoryIdChange }: ModalProps) => {
  const [t] = useTranslation('translation')
  const fetchDirs = ApiFetchDirs()
  // const { classes } = ApiFetchDatasets(dataType, directoryId)

  const queryClient = useQueryClient()
  const directoriesQueryData = queryClient.getQueryData<DirectoryType[]>([QUERY_KEY.fetchDirectories])
  const [ directories, setDirectories ] = useState<DirectoryType[] | undefined>(directoriesQueryData)
  
  useEffect(() => {
    setDirectories(directoriesQueryData)
  }, [directoriesQueryData])

  useEffect(() => {
    if (show) {
      fetchDirs.mutate({parentId: 0, dataType})
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [show, dataType])

  const handleOpenSubDir = (parentId: number) => {
    fetchDirs.mutate({parentId, dataType})
  }

  const handleCloseSubDir = (dirs: DirectoryType[], path: string) => {
    const newDirs: DirectoryType[] = objDeepCopy(dirs)
    const pathVariables = path?.split('/') || []

    let target: DirectoryType | undefined;
    pathVariables.forEach((path: string) => {
      if (!path) return 
      target = newDirs.find((dir: DirectoryType) => dir.name === path)
    })
    
    if (target) {
      delete target['dirs']
      target['is_open'] = false         
    }
   
    queryClient.setQueryData([QUERY_KEY.fetchDirectories], newDirs)
    setDirectories(newDirs)
  }

  const handleSelectData = (path: string, id: number) => {
    selectData(path, id)
    
    queryClient.invalidateQueries({queryKey: [`${QUERY_KEY.fetchDatasets}_${dataType}`], exact: true})
    
    onDirectoryIdChange(id)

    toggle()
  }

  return (
    <Modal show={show} keyboard={true} onHide={toggle}>
      <Modal.Header>{t('ui.dataset.title')}</Modal.Header>
      <Modal.Body>
        {directories ? (
          <Row className='mx-n1 g-0'>
            {renderTree({
              dirs: directories,
              parent: '',
              depth: 0,
              openSubDir: handleOpenSubDir,
              closeSubDir: handleCloseSubDir,
              selectDir: handleSelectData,
              isTest
            })}
          </Row>
        ) : (
          t('ui.dataset.nodata')
        )}
      </Modal.Body>
      <Modal.Footer>
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
          cancel
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default SelectDataModal
