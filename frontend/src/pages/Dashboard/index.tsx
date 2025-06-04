import { useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import { logger } from 'helpers'
import { ApiFetchMenu } from 'helpers/api'
import { LocationContext, ProjectContext } from 'contexts'
import { Card, Col, Form, Row } from 'react-bootstrap'
import { GPUArea } from 'features'
import ProjectCard from 'pages/Project/ProjectCard'
import { ApiFetchProjects } from 'helpers/api'
import { ProjectType } from 'common'


const DashboardLanding = () => {
    const [t] = useTranslation('translation')

    const { updateLocationContextValue } = useContext(LocationContext)
    const { updateProjectContextValue } = useContext(ProjectContext)

    const { projects } = ApiFetchProjects('dashboard')
    const { menu } = ApiFetchMenu()

    useEffect(() => {
        logger.log(`Change Location to ${t('title.dashboard')}`)
        updateLocationContextValue({ location: 'dashboard' })

        const selectedProject = projects?.filter((project: ProjectType) => Number(project.project_id) === 0)
        updateProjectContextValue({ selectedProject: selectedProject ? selectedProject[0] : {} })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [t])

    const [menus, setMenus] = useState<any[]>([])
    useEffect(() => {
        if (!menu || menu === undefined) return

        let newMenus: any = []
        menu.forEach((m: any) => {
            if (m.key === 'vision' || m.key === 'table' || m.key === 'ts') {
                m.children.forEach((c: any) => {
                    c.parentLabel = m.label
                    c.category = c.key.replace('-', '.')

                    if (c.isUse) newMenus.push(c)
                })
            }
        })

        setMenus(newMenus)

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [menu])

    return (
        <Form>
            <Row>
                <Col sm={6} style={{ marginLeft: '5px' }}>
                    <GPUArea />
                </Col>
            </Row>
            <Row>
                {
                    menus && menus.map((menu: any) => {
                        return <Card key={menu.key}>
                            <Card.Header>{menu.parentLabel} {menu.label}</Card.Header>
                            <Card.Body>
                                <Row>
                                    {projects && projects.map((project: any) => {
                                        if (project.category !== menu.category) return null
                                        return (
                                            <Col md={6} xxl={3} key={'proj-' + project.project_id}>
                                                <ProjectCard project={project} canEdit={false} handleEditProject={() => { }} handleDeleteProject={() => { }} />
                                            </Col>
                                        )
                                    })}
                                </Row>
                            </Card.Body>
                        </Card>
                    })
                }
            </Row>
        </Form>
    )
}

export default DashboardLanding