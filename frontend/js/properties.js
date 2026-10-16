async function createProperty(propertyData) {
  return apiRequest("/properties", {
    method: "POST",
    body: JSON.stringify(propertyData)
  });
}

async function listMyProperties() {
  return apiRequest("/properties/mine");
}

async function deleteProperty(id) {
  return apiRequest(`/properties/${id}`, { method: "DELETE" });
}

async function uploadPropertyImage(id, file) {
  const formData = new FormData();
  formData.append("image", file);
  return apiRequest(`/properties/${id}/images`, {
    method: "POST",
    body: formData,
    headers: {}
  });
}

async function loadProperties(filters = {}) {
  const params = new URLSearchParams();
  if (filters.location && filters.location.trim()) {
    params.set("location", filters.location.trim());
  }
  if (filters.min_price !== undefined && filters.min_price !== "") {
    params.set("min_price", filters.min_price);
  }
  if (filters.max_price !== undefined && filters.max_price !== "") {
    params.set("max_price", filters.max_price);
  }
  if (filters.bedrooms !== undefined && filters.bedrooms !== "") {
    params.set("bedrooms", filters.bedrooms);
  }
  if (filters.bathrooms !== undefined && filters.bathrooms !== "") {
    params.set("bathrooms", filters.bathrooms);
  }
  if (filters.amenities && filters.amenities.trim()) {
    params.set("amenities", filters.amenities.trim());
  }

  const query = params.toString();
  return apiRequest(query ? `/properties?${query}` : "/properties");
}

async function getProperty(id) {
  return apiRequest(`/properties/${id}`);
}

async function listFavorites() {
  return apiRequest("/favorites");
}

async function getFavoriteIds() {
  return apiRequest("/favorites/ids");
}

async function checkFavorite(id) {
  return apiRequest(`/favorites/check/${id}`);
}

async function addFavorite(id) {
  return apiRequest(`/favorites/${id}`, { method: "POST" });
}

async function removeFavorite(id) {
  return apiRequest(`/favorites/${id}`, { method: "DELETE" });
}

