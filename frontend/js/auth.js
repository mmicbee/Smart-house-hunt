function saveAuth(data) {
  if (data.token) {
    localStorage.setItem("token", data.token);
  }

  if (data.user) {
    localStorage.setItem("user", JSON.stringify(data.user));
  }
}

async function registerUser(userData) {
  return apiRequest("/auth/register", {
    method: "POST",
    body: JSON.stringify(userData)
  });
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  window.location.href = "login.html";
}

function getCurrentUser() {
  const user = localStorage.getItem("user");
  return user ? JSON.parse(user) : null;
}
