import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import axios from 'axios';

function AuthCallback({ onLogin }) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const token = searchParams.get('token');
    
    console.log('AuthCallback - Token received:', token ? 'Yes' : 'No');

    if (token) {
      // Store token
      localStorage.setItem('token', token);
      console.log('Token stored in localStorage');

      // Set axios default header
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      console.log('Fetching user profile...');

      // Fetch user profile
      axios.get('/api/v1/auth/profile')
        .then(response => {
          console.log('Profile fetched:', response.data);

          // Call parent login handler
          if (onLogin) {
            onLogin(response.data, token);
          }
          
          // Redirect to dashboard
          console.log('Redirecting to dashboard...');
          navigate('/dashboard');
        })
        .catch(error => {
          console.error('Failed to fetch profile:', error);
          console.error('Error details:', error.response?.data);
          navigate('/login');
        });
    } else {
      console.error('No token found in URL');
      navigate('/login');
    }
  }, [searchParams, navigate, onLogin]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-900 to-gray-800">
      <div className="text-center">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
        <p className="text-white text-lg">Completing login...</p>
      </div>
    </div>
  );
}

export default AuthCallback;

