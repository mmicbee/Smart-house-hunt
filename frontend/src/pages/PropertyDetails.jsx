import { useParams, Link } from 'react-router-dom'

export default function PropertyDetails() {
  const { id } = useParams()

  // Mock data for now
  const property = {
    id,
    title: 'Modern 2-Bedroom Apartment',
    location: 'Kakamega Town, near MMUST',
    price: 15000,
    bedrooms: 2,
    bathrooms: 1,
    verified: true,
    description: 'A clean and modern 2-bedroom apartment located close to MMUST. Ideal for students or young professionals. Water and electricity available. Secure compound.',
    amenities: ['WiFi', 'Parking', 'Water', 'Electricity', 'Security'],
    landlord: 'John Wekesa',
    image: 'https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?w=800'
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <Link to="/" className="text-primary-600 hover:underline mb-4 inline-block">
        ← Back to listings
      </Link>

      <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
        <img
          src={property.image}
          alt={property.title}
          className="w-full h-80 object-cover"
        />

        <div className="p-6 md:p-8">
          <div className="flex flex-col md:flex-row md:justify-between md:items-start gap-4">
            <div>
              <h1 className="text-2xl md:text-3xl font-bold text-gray-900">
                {property.title}
              </h1>
              <p className="text-gray-500 mt-1">{property.location}</p>
              {property.verified && (
                <span className="inline-block mt-2 bg-green-100 text-green-700 text-sm px-3 py-1 rounded-full">
                  ✓ Verified Landlord
                </span>
              )}
            </div>
            <div className="text-right">
              <p className="text-3xl font-bold text-primary-600">
                KES {property.price.toLocaleString()}
              </p>
              <p className="text-gray-500 text-sm">per month</p>
            </div>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-8">
            <div className="bg-gray-50 rounded-lg p-4 text-center">
              <p className="text-2xl font-semibold">{property.bedrooms}</p>
              <p className="text-gray-500 text-sm">Bedrooms</p>
            </div>
            <div className="bg-gray-50 rounded-lg p-4 text-center">
              <p className="text-2xl font-semibold">{property.bathrooms}</p>
              <p className="text-gray-500 text-sm">Bathrooms</p>
            </div>
            <div className="bg-gray-50 rounded-lg p-4 text-center">
              <p className="text-2xl font-semibold">1</p>
              <p className="text-gray-500 text-sm">Parking</p>
            </div>
            <div className="bg-gray-50 rounded-lg p-4 text-center">
              <p className="text-2xl font-semibold">Yes</p>
              <p className="text-gray-500 text-sm">Water</p>
            </div>
          </div>

          <div className="mt-8">
            <h2 className="text-lg font-semibold mb-2">Description</h2>
            <p className="text-gray-700 leading-relaxed">{property.description}</p>
          </div>

          <div className="mt-8">
            <h2 className="text-lg font-semibold mb-3">Amenities</h2>
            <div className="flex flex-wrap gap-2">
              {property.amenities.map(a => (
                <span key={a} className="bg-primary-50 text-primary-700 px-3 py-1 rounded-full text-sm">
                  {a}
                </span>
              ))}
            </div>
          </div>

          <div className="mt-10 flex flex-col sm:flex-row gap-4">
            <button className="flex-1 bg-primary-600 text-white py-3 rounded-lg font-semibold hover:bg-primary-700 transition">
              Contact Landlord
            </button>
            <button className="flex-1 border border-primary-600 text-primary-600 py-3 rounded-lg font-semibold hover:bg-primary-50 transition">
              Save to Favorites
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
