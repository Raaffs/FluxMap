import { UserRole } from "../hooks/types";

export const normalizeAccessMap = (
  rawAccess: Record<UserRole, number[]|null>
): Record<number, UserRole> => {
  const accessMap: Record<number, UserRole> = {};
  const roles: UserRole[] = ["admin", "manager", "user"];

  for (const role of roles) {
    const projectIds = rawAccess[role] || [];
    for (const id of projectIds) {
      // Only assign if it’s not already set (higher priority roles win)
      if (!accessMap[id]) {
        accessMap[id] = role;
      }
    }
  }
  return accessMap;
};
