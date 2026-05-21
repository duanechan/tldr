import { TokenKey } from "@/lib/constants";

export function useAuth() {
  const token = localStorage.getItem(TokenKey);
  const isAuthenticated = !!token && token !== "null" && token !== "undefined";
  return { isAuthenticated };
}
