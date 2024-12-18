import {useRetrieveProjectsFrom} from "../hooks/projects";


export const ProjectComponent = ({URI}:{URI:string}) => {
const { projects, loading, error } = useRetrieveProjectsFrom(URI);

if (loading) {
  return <div>Loading projects...</div>;
}

if (error) {
  return <div>Error: {error}</div>;
}

    return (
      <div>
        <h1>Project List</h1>
        <ul>
          {projects.map((project) => (
            <li key={project.projectID}>
              <strong>{project.projectName}</strong>
              <p>{project.projectDescription || "No description available"}</p>
              <p>
                Start: {project.projectStartDate} | Due: {project.projectDueDate}
              </p>
            </li>
          ))}
        </ul>
      </div>
    );
};
  
  