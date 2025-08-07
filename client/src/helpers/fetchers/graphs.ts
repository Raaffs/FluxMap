import { Graphs } from "../../hooks/types";

export async function fetchCompletedGraph(): Promise<Graphs> {
  const res = await fetch("http://localhost:4000/api/graph/completed", {
    method: "GET",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
  });

  return res.json();
}
