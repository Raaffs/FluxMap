import { tasks } from "../hooks/types";


export const groupTasksWithMap = (
    tasksData: tasks[], 
    groupByField: keyof tasks, 
    filterField?: keyof tasks, 
    filterValue?: string
  ): { [key: string]: string | number }[] => {
    const map = new Map<string, number>();
  
  tasksData.forEach(task => {
      // Apply filter if needed
      if (filterField && task[filterField] !== filterValue) return;
  
      const key = task[groupByField] as string;
      if (map.has(key)) {
        map.set(key, map.get(key)! + 1);  // Increment the count
      } else {
        map.set(key, 1);  // Initialize the count
      }
    });
  
    // Convert Map to array format
    return Array.from(map, ([key, count]) => ({ [groupByField]: key, count }));
  };


  export const groupTasksByPeriod = (
  tasksData: tasks[],
  groupByField: keyof tasks,
  period: "week" | "day"
): { period: string; count: number }[] => {
  const map = new Map<string, number>();

  tasksData.forEach((task) => {
    const date = task[groupByField] as string;
    if (!date || date === "0") return; // Skip invalid dates

    const taskDate = new Date(date);

    // Calculate ISO Week Number
    const year = taskDate.getFullYear();
    const startDate = new Date(year, 0, 1);
    const diff = taskDate.getTime() - startDate.getTime();
    const days = Math.floor(diff / (1000 * 3600 * 24));
    const weekNumber = Math.ceil((days + 1) / 7); // Adding 1 to include the first week

    const periodLabel =
      period === "week"
        ? `${year}-W${String(weekNumber).padStart(2, '0')}`
        : taskDate.toISOString().split("T")[0]; // Use 'YYYY-MM-DD' for daily grouping

    if (map.has(periodLabel)) {
      map.set(periodLabel, map.get(periodLabel)! + 1);
    } else {
      map.set(periodLabel, 1);
    }
  });

  // Sorting periods numerically for correct chronological order
  const sortedEntries = Array.from(map.entries()).sort(([periodA], [periodB]) => {
    return periodA.localeCompare(periodB);
  });

  return sortedEntries.map(([period, count]) => ({ period, count }));
};
