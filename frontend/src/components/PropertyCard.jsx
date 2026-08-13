import { Link } from 'react-router-dom'

export default function PropertyCard({ property }) {
  return (
    <Link to={`/properties/${property.id}`} className="block group">
      <div className="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-md transition duration-300">
        <div className="relative h-48 bg-gray-200">
          <img
            src={property.image || 'https://images.unsplash.com/photo-1560448204-e02f11c3be0e?w=600'}
            alt={property.title}
            className="w-full h-full object-cover group-hover:scale-105 transition duration-300"
          />
          {property.verified && (
            <span className="absolute top-3 left-3 bg-green-500 text-white text-xs px-2 py-1 rounded-full">
              ✓ Verified
            </span>
          )}
        </div>
        
        <div className="p-4">
          <div className="flex justify-between items-start">
            <h3 className="font-semibold text-lg text-gray-900 line-clamp-1">
              {property.title}
            </h3>
            <p className="font-bold text-primary-600">
              KES {property.price.toLocaleString()}
            </p>
          </div>
          
          <p className="text-gray-500 text-sm mt-1">{property.location}</p>
          
          <div className="flex items-center gap-4 mt-3 text-sm text-gray-600">
            <span>🛏 {property.bedrooms} bed</span>
            <span>🚿 {property.bathrooms} bath</span>
            {property.amenities?.includes('wifi') && <span>📶 WiFi</span>}
          </div>
        </div>
      </div>
    </Link>
  )
}
