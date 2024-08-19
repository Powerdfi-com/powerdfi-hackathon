"use client"
import { AssetContext } from '@/app/context/context'
import Image from 'next/image'
import React, { useContext } from 'react'

const Property = () => {
    const asset = useContext(AssetContext);
    return <ul className='grid grid-cols-2 gap-8 mt-6'>
        {
            asset.properties.map((property) => <li key={property[0]}>
                <div className="flex justify-between max-w-sm">
                    <div className="text-white text-[24px]">{property[0]}</div>
                    <div className='text-white text-[24px]'>{property[1]}</div>
                </div>
            </li>)
        }
    </ul>

}

export default Property