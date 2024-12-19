import React from "react";
import { ProjectComponent } from "../../components/project";

export const AllProjects = () => {
  return <ProjectComponent URI="http://localhost:4000/api/projects" />
};

export const AdminProjects = () => {
  return <ProjectComponent URI="http://localhost:4000/api/projects/admin" />
};

export const ManagerProjects = () => {
  return <ProjectComponent URI="http://localhost:4000/api/projects/manager" />
};
export const AllocatedProjects = () => {
  return <ProjectComponent URI="http://localhost:4000/api/projects/assigned" />
};

