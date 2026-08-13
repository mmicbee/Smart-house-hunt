import { Link } from 'react-router-dom'

export default function Navbar() {
  return (
    <nav className="bg-white shadow-sm sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <Link to="/" className="text-xl font-bold text-primary-600">
            🏠 Smart House Hunt
          </Link>
          
          <div className="hidden md:flex items-center space-x-6">
            <Link to="/" className="text-gray-700 hover:text-primary-600">
              Browse Houses
            </Link>
            <Link to="/login" className="text-gray-700 hover:text-primary-600">
              Login
            </Link>
            <Link
              to="/register"
              className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700 transition"
            >
              Register
            </Link>
          </div>
        </div>
      </div>
    </nav>
  )
}
