// Admin API calls

async function getAdminReports() {
  return apiRequest("/admin/reports");
}

async function getAdminProperties(filters = {}) {
  const params = new URLSearchParams();
  if (filters.status && filters.status !== "all") {
    params.set("status", filters.status);
  }
  if (filters.property_type && filters.property_type !== "all") {
    params.set("property_type", filters.property_type);
  }
  const query = params.toString() ? `?${params.toString()}` : "";
  return apiRequest(`/admin/properties${query}`);
}

async function getAdminProperty(id) {
  return apiRequest(`/admin/properties/${id}`);
}

async function updatePropertyStatus(id, status) {
  return apiRequest(`/admin/properties/${id}/status`, {
    method: "PUT",
    body: JSON.stringify({ status })
  });
}

async function approveProperty(id) {
  return apiRequest(`/admin/properties/${id}/approve`, {
    method: "PUT"
  });
}

async function rejectProperty(id) {
  return apiRequest(`/admin/properties/${id}/reject`, {
    method: "PUT"
  });
}

async function deleteAdminProperty(id) {
  return apiRequest(`/admin/properties/${id}`, {
    method: "DELETE"
  });
}

async function getAdminUsers(filters = {}) {
  const params = new URLSearchParams();
  if (filters.status && filters.status !== "all") {
    params.set("status", filters.status);
  }
  if (filters.role && filters.role !== "all") {
    params.set("role", filters.role);
  }
  if (filters.search) {
    params.set("search", filters.search);
  }
  const query = params.toString() ? `?${params.toString()}` : "";
  return apiRequest(`/admin/users${query}`);
}

async function updateUserStatus(id, status) {
  return apiRequest(`/admin/users/${id}/status`, {
    method: "PUT",
    body: JSON.stringify({ status })
  });
}
