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
  if (filters.property_type && filters.property_type.trim()) {
    params.set("property_type", filters.property_type.trim());
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

function getPropertyTypeLabel(propertyType) {
  const t = (propertyType || "").toLowerCase();
  switch (t) {
    case "bedsitter":
      return "Bedsitter";
    case "studio":
      return "Studio Apartment";
    case "house":
      return "Standalone House";
    case "apartment":
    default:
      return "Apartment";
  }
}

function getBedroomDisplay(property) {
  const type = (property.property_type || "").toLowerCase();
  const beds = parseInt(property.bedrooms, 10);
  if (type === "bedsitter") {
    return "Bedsitter (Single room)";
  }
  if (type === "studio") {
    return "Studio (Open plan)";
  }
  if (beds === 0) {
    return "0 Bedrooms";
  }
  if (beds === 1) {
    return "1 Bedroom";
  }
  return `${beds} Bedrooms`;
}

function getPropertyImageUrl(property, basePath = "") {
  if (property && property.images && property.images.length > 0) {
    const imgPath = property.images[0].file_path;
    if (imgPath.startsWith("http")) return imgPath;
    return `${API_BASE_URL.replace("/api", "")}/${imgPath}`;
  }
  const type = (property && property.property_type ? property.property_type : "").toLowerCase();
  if (type === "bedsitter") return `${basePath}assets/images/bedsitter.svg`;
  if (type === "studio") return `${basePath}assets/images/studio.svg`;
  if (type === "house") return `${basePath}assets/images/house.svg`;
  return `${basePath}assets/images/apartment.svg`;
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

