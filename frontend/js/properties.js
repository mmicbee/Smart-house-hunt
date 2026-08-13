async function loadProperties() {
  try {
    const result = await apiRequest("/properties");
    console.log("Properties:", result);
  } catch (error) {
    console.error("Could not load properties:", error);
  }
}
