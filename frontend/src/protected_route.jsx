import React from 'react';
import { Navigate } from 'react-router-dom'

export default function ProtectedRoute({children}) {
    const token = localStorage.getItem("accessToken")

    if(!token) {
        return <Navigate to="/LoginUser" replace ></Navigate>
    }

    return children;
}