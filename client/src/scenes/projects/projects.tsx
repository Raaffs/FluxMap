import React from "react";
import { ProjectListComponent } from "../../components/project";

export const AllProjects = () => {
  return <ProjectListComponent URI="http://localhost:4000/api/projects" />
};

export const AdminProjects = () => {
  return <ProjectListComponent URI="http://localhost:4000/api/projects/admin" />
};

export const ManagerProjects = () => {
  return <ProjectListComponent URI="http://localhost:4000/api/projects/manager" />
};
export const AllocatedProjects = () => {
  return <ProjectListComponent URI="http://localhost:4000/api/projects/assigned" />
};


