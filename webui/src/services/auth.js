export function getAuth() {
  const token = localStorage.getItem("token") || "";
  const userId = Number(localStorage.getItem("userId") || "0");
  const username = localStorage.getItem("username") || "";
  return { token, userId, username, isAuthed: !!token && userId > 0 };
}
export function setAuth({ token, userId, username }) {
  localStorage.setItem("token", token);
  localStorage.setItem("userId", String(userId));
  localStorage.setItem("username", username || "");
}
export function clearAuth() {
  localStorage.removeItem("token");
  localStorage.removeItem("userId");
  localStorage.removeItem("username");
}
