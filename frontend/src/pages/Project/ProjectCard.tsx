import { Card, Col, Dropdown, Row } from "react-bootstrap"
import { useLocation, useNavigate } from "react-router-dom"
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { ProjectType } from 'common'
import { stateColors } from 'appConstants'
import { convertUtcTime } from "helpers"

const CardTitle = styled.div`
    cursor: pointer;
`

type ProjectCardPros = { 
    project: ProjectType, 
    handleEditProject: (selectedProject: ProjectType) => void ,
    handleDeleteProject: (selectedProject: ProjectType) => void ,
    canEdit?:boolean
}

const ProjectCard = ({ project, canEdit=true, handleEditProject, handleDeleteProject }: ProjectCardPros) => {
    const [t] = useTranslation('translation')
    const navigate = useNavigate()
    const location = useLocation()

    const selectProject = () => {
        let targetPath = `./${project.project_id}`
        if (location.pathname === '/dashboard' && project.category) {
            const projectPath = project.category.split('.')
            targetPath = `/${projectPath[0]}/${projectPath[1]}/${project.project_id}`
        }
        navigate(targetPath)
    }

    let border = 'none'
    if (project && ((project.num_tests && project.num_tests > 0) || (project.num_trains && project.num_trains > 0))) {
        border = '1px solid blue'
    }

    return (
        <Card className="d-block" style={{border: border, cursor:'pointer'}} onDoubleClick={selectProject}>
            <Card.Body>
                {canEdit&&<Dropdown className="card-widgets" align="end" style={{cursor:'pointer'}}>
                    <Dropdown.Toggle
                        variant="link"
                        as="a"
                        className="card-drop arrow-none cursor-pointer p-0 shadow-none"
                    >
                        <i className="dripicons-dots-3"></i>
                    </Dropdown.Toggle>

                    <Dropdown.Menu>
                        <Dropdown.Item onClick={() => handleEditProject(project)}>
                            <i className="mdi mdi-pencil me-1"></i>{t('ui.label.edit')}
                        </Dropdown.Item>
                        {
                            border === 'none' && (
                                <Dropdown.Item onClick={() => handleDeleteProject(project)}>
                                    <i className="mdi mdi-delete me-1"></i>{t('ui.label.delete')}
                                </Dropdown.Item>
                            )
                        }
                    </Dropdown.Menu>
                </Dropdown>}

                <h4 className="mt-0">
                    <CardTitle className="text-title">
                        {project.project_name}
                    </CardTitle>
                </h4>

                <div>
                    <p className="text-muted font-13 my-3">
                        {t('ui.label.created')} {convertUtcTime(project.created_at)}
                    </p>
                </div>

                {project.description 
                ? (
                    <p className="text-muted font-13 my-3">
                        {project.description}
                    </p>
                )
                : (
                    <p className="text-muted font-13 my-3">
                        -
                    </p>
                )}

                <Row>
                    <Col>
                    <span className="pe-2 text-nowrap mb-2 d-inline-block">
                        <i className="mdi mdi-format-list-bulleted-type" style={{color: stateColors['total']}}> </i>
                        <b>{project.num_trials}</b> Trials
                    </span>
                    </Col>

                    <Col>
                    <span className="pe-2 text-nowrap mb-2 d-inline-block">
                        <i className="mdi mdi-expansion-card-variant" style={{color: stateColors['train']}}> </i>
                        <b>{project.num_trains}</b> Trains
                    </span>
                    </Col>

                    <Col>
                    <span className="pe-2 text-nowrap mb-2 d-inline-block">
                        <i className="mdi mdi-expansion-card-variant" style={{color: stateColors['test']}}> </i>
                        <b>{project.num_tests}</b> Tests
                    </span>
                    </Col>
                </Row>
               
            </Card.Body>
        </Card>
    )
}

export default ProjectCard