"use client"
import { AssetContext } from '@/app/context/context'
import Image from 'next/image'
import React, { useContext } from 'react'

const About = () => {
    const asset = useContext(AssetContext);
    return <p className='text-[18px] leading-relaxed text-white'>{asset.description}</p>
}

export default About