"use client"
import { useRouter } from 'next/navigation'
import React, { useEffect } from 'react'

const Home = () => {
    const router = useRouter()
    useEffect(() => {
        router.replace("/i")
    }, [])
    return (
        <div>Home</div>
    )
}

export default Home