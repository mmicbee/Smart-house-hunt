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

async function loadProperties() {
  try {
    const result = await apiRequest("/properties");
    console.log("Properties:", result);
  } catch (error) {
    console.error("Could not load properties:", error);
  }
}
