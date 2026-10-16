async function createViewingRequest(propertyId, preferredDate, message = "") {
  return apiRequest("/viewings", {
    method: "POST",
    body: JSON.stringify({
      property_id: propertyId,
      preferred_date: preferredDate,
      message: message
    })
  });
}

async function getMyViewings() {
  return apiRequest("/viewings/mine");
}

async function getLandlordViewings() {
  return apiRequest("/viewings/landlord");
}

async function getViewing(id) {
  return apiRequest(`/viewings/${id}`);
}

async function updateViewingStatus(id, status) {
  return apiRequest(`/viewings/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status })
  });
}

async function acceptViewing(id) {
  return apiRequest(`/viewings/${id}/accept`, {
    method: "POST"
  });
}

async function rejectViewing(id) {
  return apiRequest(`/viewings/${id}/reject`, {
    method: "POST"
  });
}

async function cancelViewing(id) {
  return apiRequest(`/viewings/${id}/cancel`, {
    method: "POST"
  });
}
