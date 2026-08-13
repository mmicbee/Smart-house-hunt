import { useState } from 'react'
import PropertyCard from '../components/PropertyCard'

// Temporary mock data — will be replaced by API calls later
const mockProperties = [
  {
    id: 1,
    title: 'Modern 2-Bedroom Apartment',
    location: 'Kakamega Town, near MMUST',
    price: 15000,
    bedrooms: 2,
    bathrooms: 1,
    verified: true,
    amenities: ['wifi', 'parking'],
    image: 'https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?w=600'
  },
  {
    id: 2,
    title: 'Spacious 3-Bedroom Bungalow',
    location: 'Amaseko, Kakamega',
    price: 22000,
    bedrooms: 3,
    bathrooms: 2,
    verified: true,
    amenities: ['wifi', 'water'],
    image: 'https://images.unsplash.com/photo-1564013799919-ab600027ffc6?w=600'
  },
  {
    id: 3,
    title: 'Affordable Bedsitter',
    location: 'Shinyalu, Kakamega',
    price: 6500,
    bedrooms: 1,
    bathrooms: 1,
    verified: false,
    amenities: ['water'],
    image: 'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=600'
  },
  {
    id: 4,
    title: 'Luxury 4-Bedroom Maisonette',
    location: 'Milimani, Kakamega',
    price: 45000,
    bedrooms: 4,
    bathrooms: 3,
    verified: true,
    amenities: ['wifi', 'parking', 'security'],
    image: 'https://images.unsplash.com/photo-1600596542815-ffad4c1539a9?w=600'
  },
  {
    id: 5,
    title: 'Student Friendly 1-Bedroom',
    location: 'Near Kibabii University',
    price: 9000,
    bedrooms: 1,
    bathrooms: 1,
    verified: false,
    amenities: ['wifi'],
    image: 'https://images.unsplash.com/photo-1493809842364-78817add7ffb?w=600'
  },
  {
    id: 6,
    title: 'Family House with Compound',
    location: 'Lurambi, Kakamega',
    price: 18000,
    bedrooms: 3,
    bathrooms: 2,
    verified: true,
    amenities: ['parking', 'water'],
    image: 'https://images.unsplash.com/photo-1600585154340-be6161a56a0c?w=600'
  },
]

export default function Home() {
  const [search, setSearch] = useState('')
  const [maxPrice, setMaxPrice] = useState('')

  const filtered = mockProperties.filter(p => {
    const matchesSearch = p.title.toLowerCase().includes(search.toLowerCase()) ||
                          p.location.toLowerCase().includes(search.toLowerCase())
    const matchesPrice = maxPrice ? p.price <= Number(maxPrice) : true
    return matchesSearch && matchesPrice
  })

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Hero / Search */}
      <div className="bg-primary-600 rounded-2xl p-8 mb-10 text-white">
        <h1 className="text-3xl md:text-4xl font-bold mb-2">
          Find Your Next Home in Kenya
        </h1>
        <p className="text-primary-100 mb-6">
          Verified listings • AI recommendations • Fraud protection
        </p>
        
        <div className="flex flex-col sm:flex-row gap-3">
          <input
            type="text"
            placeholder="Search by location or title..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="flex-1 px-4 py-3 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-white"
          />
          <input
            type="number"
            placeholder="Max price (KES)"
            value={maxPrice}
            onChange={(e) => setMaxPrice(e.target.value)}
            className="w-full sm:w-40 px-4 py-3 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-white"
          />
          <button className="bg-white text-primary-600 font-semibold px-6 py-3 rounded-lg hover:bg-primary-50 transition">
            Search
          </button>
        </div>
      </div>

      {/* Results */}
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-semibold">
          {filtered.length} properties found
        </h2>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map(property => (
          <PropertyCard key={property.id} property={property} />
        ))}
      </div>

      {filtered.length === 0 && (
        <div className="text-center py-16 text-gray-500">
          No properties match your search. Try adjusting the filters.
        </div>
      )}
    </div>
  )
}
