import React from "react";
import { ProjectListComponent } from "../../components/project";
import { UserRole } from "../../hooks/types";

type ProjectPageProps = {
  setUserProjectRoleMap: React.Dispatch<React.SetStateAction<Record<number, UserRole>>>;
};

export const AllProjects: React.FC<ProjectPageProps> = ({ setUserProjectRoleMap }) => (
  <ProjectListComponent
    URI="http://localhost:4000/api/projects"
    setUserProjectRoleMap={setUserProjectRoleMap}
  />
);

export const AdminProjects: React.FC<ProjectPageProps> = ({ setUserProjectRoleMap }) => (
  <ProjectListComponent
    URI="http://localhost:4000/api/projects/admin"
    setUserProjectRoleMap={setUserProjectRoleMap}
  />
);

export const ManagerProjects: React.FC<ProjectPageProps> = ({ setUserProjectRoleMap }) => (
  <ProjectListComponent
    URI="http://localhost:4000/api/projects/manager"
    setUserProjectRoleMap={setUserProjectRoleMap}
  />
);

export const AllocatedProjects: React.FC<ProjectPageProps> = ({ setUserProjectRoleMap }) => (
  <ProjectListComponent
    URI="http://localhost:4000/api/projects/assigned"
    setUserProjectRoleMap={setUserProjectRoleMap}
  />
);
