export function useAuth() {
  const token = localStorage.getItem("ACCESS_TOKEN");
  const isAuthenticated = !!token && token !== "null" && token !== "undefined";
  return { isAuthenticated };
}
