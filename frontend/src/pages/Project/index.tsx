import { useEffect, useState, useContext } from 'react'
import { Button, Col, Row } from "react-bootstrap"
import { useLocation } from 'react-router-dom'

import ProjectCard from "./ProjectCard"
import { logger } from 'helpers'
import { ProjectModal } from 'features'
import { ProjectType } from 'common'
import { CustomConfirm } from 'components'
import { useTranslation } from 'react-i18next'
import { LocationContext, ProjectContext } from 'contexts'
import { ApiFetchProjects, ApiDeleteProject, ApiCreateProject, ApiUpdateProject } from 'helpers'

const ProjectLanding = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()
    const pathVariables = location.pathname.split('/')

    const { locationContextValue, updateLocationContextValue } = useContext(LocationContext)
    const { updateProjectContextValue } = useContext(ProjectContext)

    const category = locationContextValue.location.split('.project')[0];
    const { projects } = ApiFetchProjects(category)
    const createProject = ApiCreateProject(category)
    const updateProject = ApiUpdateProject(category)  
    const deleteProject = ApiDeleteProject(category)

    useEffect(() => {
        updateProjectContextValue({selectedProject: {project_id: 0}})
        updateLocationContextValue({location: `${pathVariables[1]}.${pathVariables[2]}.project`})
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [location])

    useEffect(() => {
        if (locationContextValue.location !== `${pathVariables[1]}.${pathVariables[2]}.project`) return
        logger.log(`Change Location to ${t('title.' + locationContextValue.location)}`)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [locationContextValue])

    const closeProjectModal = () => {
        setModalProject({
          ...modalProject,
          show: false,
        })
    }

    const addPrj = (data:any) => {
        createProject.mutate(data)
    }

    const editPrj = (data:any) => {
        updateProject.mutate(data)
    }

    const [modalProject, setModalProject] = useState({
        show: false,
        isEdit: false,
        data: {},
        category: '',
        projects: {},
        toggle: closeProjectModal,
        handleSubmit: addPrj,
    })

    const toggleAddProject = () => {
        setModalProject({
            ...modalProject,
            show: true,
            isEdit: false,
            data: {},
            category: category,
            projects: projects,
            handleSubmit: addPrj,
        })
    }

    const toggleEditProject = (selectedProject: ProjectType) => {
        setModalProject({
            ...modalProject,
            show: true,
            isEdit: true,
            data: selectedProject || {},
            category: category,
            projects: projects,
            handleSubmit: editPrj,
        })
    }

    const onDeleteProject = (selectedProject: ProjectType) => {
        if (!projects || projects.length < 1) return
        if (!selectedProject) return
    
        const prjName = projects.find(
            (project: { project_id: number | undefined }) => selectedProject?.project_id === project.project_id
        ).project_name
    
        CustomConfirm({
            onConfirm: () => {
                deleteProject.mutate(Number(selectedProject.project_id))
            },
            onCancel: () => {},
            message: t('ui.confirm.delete.project', { name: prjName }),
        })
    }

    return (
        <>
            <Row className="mb-1">
                <Col sm={4}>
                    <Button variant="danger" className="rounded-pill mb-1" onClick={toggleAddProject}>
                        <i className="mdi mdi-plus"></i> {t('ui.project.add')}
                    </Button>
                </Col>
                {/* <Col sm={8}>
                    <div className="text-sm-end">
                        <div className="btn-group mb-3">
                            <Button variant="primary">All</Button>
                        </div>
                        <ButtonGroup className="mb-3 ms-1">
                            <Button variant="light">Idle</Button>
                            <Button variant="light">Train</Button>
                            <Button variant="light">Test</Button>
                        </ButtonGroup>
                    </div>
                </Col> */}
            </Row>
            <Row>
                {projects && projects.map((project: any) => {
                    return (
                        <Col md={6} xxl={3} key={'proj-' + project.project_id}>
                            <ProjectCard project={project} handleEditProject={toggleEditProject} handleDeleteProject={onDeleteProject} />
                        </Col>
                    )
                })}
            </Row>
            <ProjectModal {...modalProject} />
        </>
    )
}

export default ProjectLanding