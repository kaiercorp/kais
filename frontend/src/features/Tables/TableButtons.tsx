import { KeyboardEvent, useEffect, useState, useContext } from 'react'
import { Button, Form, InputGroup } from 'react-bootstrap'
import { Link, useLocation } from 'react-router-dom'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { CustomConfirm } from 'components'
import { TrialContext, FilterContext } from 'contexts'
import { ApiDeleteTrials, ApiDownloadFile, convertNowTimeForFilename } from 'helpers'
import { APICore } from 'helpers/api/apiCore'
import { useQueryClient } from '@tanstack/react-query'


const ButtonArea = styled.div`
  width: 100%;
  height: 35px;
  margin-bottom: 10px;
`

const ButtonAreaLeft = styled.div`
  float: left;
  & button {
    margin-right: 5px;
  }
`

const ButtonAreaRight = styled.div`
  float: right;
  padding-top: 10px;
  margin-right: 20px;
`
const SearchArea = styled.div`
  margin-left: 5px;
  display: inline-block;
  & span {
    padding: 5px;
    cursor: pointer;
  }
`

type TableButtonsProps = {
  onSearch?: any
  openCompareModal?: () => void
  openFilterModal?: () => void
}

const api = new APICore()

const TableButtons = ({ onSearch, openCompareModal, openFilterModal }: TableButtonsProps) => {
  const [t] = useTranslation('translation')
  const location = useLocation()
  
  const [projectpath, setProjectpath] = useState(location.pathname)

  const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)
  const { filterContextValue, updateFilterContextValue } = useContext(FilterContext)
  const { selectedRows } = trialContextValue
  const { useFilter} = filterContextValue 
  
  const deleteTrials = ApiDeleteTrials()

  const queryClient = useQueryClient()
  const user = api.getLoggedInUser()

  useEffect(() => {
    let newpath = ''
    let paths = projectpath.split('/')
    paths = paths.slice(0, paths.length - 1)
    newpath = paths.join('/')
    setProjectpath(newpath)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const handleDelete = () => {
    if (typeof selectedRows === 'undefined') return 

    if (
      selectedRows.length > 0 &&
      selectedRows.filter((row: any) => ['train', 'additional_train', 'test'].includes(row.state)).length < 1
    ) {
      CustomConfirm({
        onConfirm: () => {
          deleteTrials.mutate(
            selectedRows.map((row: any) => row.trial_id)
          )
          updateTrialContextValue({selectedRows: []})
        },
        onCancel: () => { },
        message: t('ui.confirm.delete', { count: selectedRows.length }),
      })
    }
  }

  const handleDownloadReport = (e:any) => {
    if (typeof selectedRows === 'undefined') return 
    e.stopPropagation()

    if (
      selectedRows.length > 0 &&
      selectedRows.filter((row: any) => ['train', 'additional_train', 'cancel', 'fail', 'test', 'finish_test', 'test_cancel', 'test_fail', 'idle'].includes(row.state)).length < 1
    ) {
      const trial = selectedRows[0]
      const ids = selectedRows.map((row: any) => row.trial_id).join(',')
      const filename = `${trial.engine_type}_${ids}_${convertNowTimeForFilename()}_report.zip`
      ApiDownloadFile(queryClient, `/trial/report/downloads`, filename, {trial_ids: ids})
    }
  }

  const toggleUseFilter = () => {
    updateFilterContextValue({useFilter: !useFilter})
  }

  const [searchValue, setSearchValue] = useState('')
  const handleChange = (e: any) => {
    setSearchValue(e.target.value)
  }
  const handleClick = (e: any) => {
    onSearch(searchValue)
  }
  const handleEnter = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      onSearch(searchValue)
    }
  }

  return (
    <ButtonArea>
      <ButtonAreaLeft>
        {/* {
          openCompareModal && (
            <Button
              variant='outline-primary'
              onClick={openCompareModal}
              disabled={
                selectedRows.length === 0 ||
                !selectedRows.every(
                  (sr: any) =>
                    sr.state === 'finish' ||
                    sr.state === 'finish_test' ||
                    (sr.state === 'finish-fail' && sr.best_model_download_path)
                )
              }
            >
              {t('button.compare')}
            </Button>
          )
        } */}
        {
          <Button
            variant='outline-danger'
            onClick={handleDelete}
            disabled={
              typeof selectedRows === 'undefined' ||
              selectedRows.length === 0 ||
              selectedRows.filter((row: any) => ['train', 'additional_train', 'test'].includes(row.state)).length > 0
            }
          >
            {t('button.delete.selected')}
          </Button>
        }
        {
          openFilterModal && (
            <Button variant='outline-info' onClick={openFilterModal}>
              <i className='mdi mdi-filter-outline' />
              <span>{t('button.filter')}</span>
            </Button>
          )
        }
        <Button variant='outline-info' onClick={toggleUseFilter}>
          {useFilter === true ? <i className='mdi mdi-filter-outline' /> : <i className='mdi mdi-filter-remove-outline' />}
          {useFilter === true ? <span>{t('button.filter.off')}</span> : <span>{t('button.filter.on')}</span>}
        </Button>
        {user && selectedRows && selectedRows.length > 0
        && (selectedRows.filter((row: any) => ['train', 'additional_train', 'cancel', 'fail', 'test', 'finish_test', 'test_cancel', 'test_fail', 'idle'].includes(row.state)).length < 1) 
        && <Button onClick={handleDownloadReport}>
          Download reports
        </Button>
        }
        <SearchArea>
          <InputGroup className=''>
            <Form.Control name='search-table' placeholder='search' value={searchValue} onKeyUp={handleEnter} onChange={handleChange} />
            <InputGroup.Text id='search' onClick={handleClick}>
              <i className='mdi mdi-magnify search-icon' />
            </InputGroup.Text>
          </InputGroup>
        </SearchArea>
      </ButtonAreaLeft>

      <ButtonAreaRight>
        <Link to={projectpath}>
          <span>{t('button.go.project')}</span>
        </Link>
      </ButtonAreaRight>

    </ButtonArea>
  )
}

export default TableButtons