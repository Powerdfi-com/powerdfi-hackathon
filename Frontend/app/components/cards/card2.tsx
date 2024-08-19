import Image from 'next/image';
import React from 'react'

const Card2 = ({ image, title }: { image: string; title: string }) => {
    return (
        <div className='flex flex-col w-full max-w-sm bg-blue-shade rounded-lg'>
            <div className='aspect-video w-full relative'>
                <Image src={image} alt="card" fill={true} objectFit='cover' className='rounded-t-lg' />
            </div>
            <div className='p-4 text-sm text-white'>{title}</div>
        </div>
    )
}

export default Card2